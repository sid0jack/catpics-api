package main

import (
	"database/sql"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

// updateCatPic godoc
// @Summary Update a cat picture
// @Description Update an existing cat picture with new image data
// @Tags catpics
// @Accept  mpfd
// @Produce  json
// @Param   id      path     string                 true  "Cat Picture ID"
// @Param   catpic  formData file                  true  "New Cat Picture"
// @Success 200     {string} string                "ok"
// @Failure 400     {object} map[string]string     "Bad Request"
// @Failure 404     {object} map[string]string     "Not Found"
// @Failure 500     {object} map[string]string     "Internal Server Error"
// @Router /catpics/{id} [put]
func updateCatPic(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		if err := r.ParseMultipartForm(10 << 20); err != nil {
			jsonError(w, "File too large or invalid", http.StatusBadRequest)
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
			jsonError(w, "Invalid file", http.StatusBadRequest)
			return
		}

		stmt, err := db.Prepare("UPDATE cat_pics SET data = ? WHERE id = ?")
		if err != nil {
			jsonError(w, "Error preparing SQL statement", http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		result, err := stmt.Exec(fileBytes, id)
		if err != nil {
			jsonError(w, "Error updating the cat picture", http.StatusInternalServerError)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			jsonError(w, "Error checking update result", http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			jsonError(w, "Cat picture not found", http.StatusNotFound)
			return
		}

		jsonResponse(w, "Cat picture updated successfully", http.StatusOK)
	}
}
