package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func TestListCatPics(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/catpics", ListCatPics(db)).Methods("GET")

	tt := []struct {
		name         string
		setupData    func(db *sql.DB)
		expectedCode int
		expectedSize int
	}{
		{
			name: "Empty Database",
			setupData: func(db *sql.DB) {
				db.Exec("DELETE FROM cat_pics")
			},
			expectedCode: http.StatusOK,
			expectedSize: 0,
		},
		{
			name: "Database With Records",
			setupData: func(db *sql.DB) {
				db.Exec("INSERT INTO cat_pics (id, data) VALUES (?, ?)", "id1", []byte("test data 1"))
				db.Exec("INSERT INTO cat_pics (id, data) VALUES (?, ?)", "id2", []byte("test data 2"))
			},
			expectedCode: http.StatusOK,
			expectedSize: 2,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupData != nil {
				tc.setupData(db)
			}

			req, err := http.NewRequest("GET", "/catpics", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedCode)
			}

			var pics []CatPicResponse
			err = json.NewDecoder(rr.Body).Decode(&pics)
			if err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			if len(pics) != tc.expectedSize {
				t.Errorf("handler returned wrong number of items: got %v want %v", len(pics), tc.expectedSize)
			}
		})
	}
}
