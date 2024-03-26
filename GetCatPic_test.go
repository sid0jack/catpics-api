package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func TestMain(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	testID := "test-get-id"
	_, err := db.Exec("INSERT INTO cat_pics (id, data) VALUES (?, ?)", testID, []byte("test cat pic data"))
	if err != nil {
		t.Fatalf("Failed to insert test record: %v", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/catpics/{id}", GetCatPicByID(db)).Methods("GET")

	tt := []struct {
		name       string
		catPicID   string
		wantStatus int
	}{
		{name: "Get Existing CatPic", catPicID: testID, wantStatus: http.StatusOK},
		{name: "Get Non-Existing CatPic", catPicID: "non-existing-id", wantStatus: http.StatusNotFound},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/catpics/"+tc.catPicID, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			req = mux.SetURLVars(req, map[string]string{"id": tc.catPicID})

			r.ServeHTTP(rr, req)

			if rr.Code != tc.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, tc.wantStatus)
			}

		})
	}
}
