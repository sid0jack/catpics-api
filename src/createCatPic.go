package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func jsonResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

func jsonError(w http.ResponseWriter, message string, statusCode int) {
	jsonResponse(w, map[string]string{"error": message}, statusCode)
}

// createCatPic godoc
// @Summary Create a cat picture
// @Description Add a new cat picture to the collection
// @Tags catpics
// @Accept  mpfd
// @Produce  json
// @Param   catpic   formData  file  true  "Cat Picture"
// @Success 201  {object}  CatPic
// @Failure 400  {object}  map[string]string
// @Router /catpics [post]
func createCatPic(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if err := r.ParseMultipartForm(10 << 20); err != nil {
            jsonError(w, "File too large", http.StatusBadRequest)
            return
        }

        file, _, err := r.FormFile("catpic")
        if err != nil {
            jsonError(w, "Invalid file", http.StatusBadRequest)
            return
        }
        defer file.Close()

        fileBytes, err := ioutil.ReadAll(file)
        if err != nil {
            jsonError(w, "Error reading file", http.StatusInternalServerError)
            return
        }

        id := uuid.NewString()

        stmt, err := db.Prepare("INSERT INTO cat_pics (id, data) VALUES (?, ?)")
        if err != nil {
            jsonError(w, "Error preparing database operation", http.StatusInternalServerError)
            return
        }
        defer stmt.Close()

        _, err = stmt.Exec(id, fileBytes)
        if err != nil {
            jsonError(w, "Error executing database operation", http.StatusInternalServerError)
            return
        }

        jsonResponse(w, map[string]string{"id": id}, http.StatusCreated)
    }
}
