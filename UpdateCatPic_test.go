package main

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func TestUpdateCatPic(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()

    // Insert a test record to update
    testID := uuid.NewString()
    _, err := db.Exec("INSERT INTO cat_pics (id, data) VALUES (?, ?)", testID, []byte("original data"))
    if err != nil {
        t.Fatalf("Failed to insert test record: %v", err)
    }

    // Setup the HTTP request with a multipart form
    var b bytes.Buffer
    w := multipart.NewWriter(&b)
    fw, err := w.CreateFormFile("catpic", "new_pic.jpg")
    if err != nil {
        t.Fatal(err)
    }
    _, err = fw.Write([]byte("new cat pic data"))
    if err != nil {
        t.Fatal(err)
    }
    w.Close()

    req, err := http.NewRequest("PUT", "/catpics/"+testID, &b)
    if err != nil {
        t.Fatal(err)
    }
    req.Header.Set("Content-Type", w.FormDataContentType())

    // Use mux to set variables for the route
    req = mux.SetURLVars(req, map[string]string{"id": testID})

    rr := httptest.NewRecorder()
    router := mux.NewRouter()
    router.HandleFunc("/catpics/{id}", UpdateCatPic(db)).Methods("PUT")
    router.ServeHTTP(rr, req)

    // Check the status code
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    // Verify the record was updated in the database
    var newData []byte
    err = db.QueryRow("SELECT data FROM cat_pics WHERE id = ?", testID).Scan(&newData)
    if err != nil {
        t.Fatalf("Failed to fetch updated record: %v", err)
    }
    if string(newData) != "new cat pic data" {
        t.Errorf("record was not updated with new data")
    }

    // Additional tests for scenarios like updating a non-existing cat pic could follow a similar pattern
}
