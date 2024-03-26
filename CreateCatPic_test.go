package main

import (
	"bytes"
	"database/sql"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func createMultipartRequest(uri string, paramName, path string) (*http.Request, error) {
	fileContent := bytes.NewBufferString("fake cat pic content")
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, path)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, fileContent)
	if err != nil {
		return nil, err
	}
	writer.Close()

	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

func setupDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Unable to open sqlite database: %v", err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS cat_pics (id TEXT PRIMARY KEY, data BLOB NOT NULL);")
	if err != nil {
		t.Fatalf("Unable to create table: %v", err)
	}
	return db
}

func TestCreateCatPic(t *testing.T) {
	db := setupDB(t)
	defer db.Close()
	handler := CreateCatPic(db)
	t.Run("Valid file upload", func(t *testing.T) {
		req, err := createMultipartRequest("/catpics", "catpic", "cat.jpg")
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}
	})

	t.Run("File too large", func(t *testing.T) {
    var b bytes.Buffer
    w := multipart.NewWriter(&b)

    fw, err := w.CreateFormFile("catpic", "large_cat.jpg")
    if err != nil {
        t.Fatalf("Failed to create form file: %v", err)
    }
    largeData := make([]byte, maxUploadSize+1) // just over 10 MB
    if _, err := fw.Write(largeData); err != nil {
        t.Fatalf("Failed to write large data to form file: %v", err)
    }
    w.Close()

    req, err := http.NewRequest("POST", "/catpics", &b)
    if err != nil {
        t.Fatalf("Failed to create request: %v", err)
    }
    req.Header.Set("Content-Type", w.FormDataContentType())

    rr := httptest.NewRecorder()
    handler := CreateCatPic(db)
    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusRequestEntityTooLarge {
        t.Errorf("Handler returned wrong status code for large file: got %v want %v", status, http.StatusRequestEntityTooLarge)
    }
})
    t.Run("Invalid file", func(t *testing.T) {
        body := new(bytes.Buffer)
        writer := multipart.NewWriter(body)

        writer.Close()

        req, err := http.NewRequest("POST", "/catpics", body)
        if err != nil {
            t.Fatalf("Failed to create request: %v", err)
        }
        req.Header.Set("Content-Type", writer.FormDataContentType())

        rr := httptest.NewRecorder()
        handler.ServeHTTP(rr, req)

        if status := rr.Code; status != http.StatusBadRequest {
            t.Errorf("Handler returned wrong status code for invalid file: got %v want %v", status, http.StatusBadRequest)
        }
    })
}
