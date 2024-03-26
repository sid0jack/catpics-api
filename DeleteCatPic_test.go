package main

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal("Failed to open sqlite database:", err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS cat_pics (id TEXT PRIMARY KEY, data BLOB NOT NULL);")
	if err != nil {
		t.Fatal("Failed to create table:", err)
	}

	return db
}

func TestDeleteCatPic(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	testID := "test-cat-pic-id"
	_, err := db.Exec("INSERT INTO cat_pics (id, data) VALUES (?, ?)", testID, []byte("test data"))
	if err != nil {
		t.Fatalf("Failed to insert test record: %v", err)
	}

	req, err := http.NewRequest("DELETE", "/catpics/"+testID, nil)
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}

	req = mux.SetURLVars(req, map[string]string{"id": testID})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteCatPic(db))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM cat_pics WHERE id = ?", testID).Scan(&count)
	if err != nil {
		t.Fatal("Failed to query database:", err)
	}
	if count != 0 {
		t.Errorf("Record was not deleted from the database")
	}
}
