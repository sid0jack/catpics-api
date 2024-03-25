package main

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
)

// deleteCatPic godoc
// @Summary Delete a cat picture
// @Description Delete a cat picture by its unique identifier
// @Tags catpics
// @Accept  json
// @Produce  json
// @Param   id  path  string  true  "Cat Picture ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 404 {object} map[string]string "Not Found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /catpics/{id} [delete]
func deleteCatPic(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		result, err := db.Exec("DELETE FROM cat_pics WHERE id = ?", id)
		if err != nil {
			jsonError(w, "Server error", http.StatusInternalServerError)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			jsonError(w, "Error checking deletion result", http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			jsonError(w, "Cat picture not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
