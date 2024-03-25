// getCatPic.go
package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// getCatPicByID godoc
// @Summary Get a cat picture by ID
// @Description Get a cat picture by its unique ID
// @Tags catpics
// @Accept  json
// @Produce  jpeg
// @Param   id   path  string  true  "Cat Picture ID"
// @Success 200  {object}  CatPicResponse
// @Failure 404  {object}  map[string]string
// @Router /catpics/{id} [get]
func getCatPicByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var data []byte
		err := db.QueryRow("SELECT data FROM cat_pics WHERE id = ?", id).Scan(&data)
		switch {
		case err == sql.ErrNoRows:
			http.NotFound(w, r)
		case err != nil:
			log.Printf("Error querying database: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		default:
			w.Header().Set("Content-Type", "image/jpeg")
			if _, err := w.Write(data); err != nil {
				log.Printf("Error writing image to response: %v", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}
	}
}
