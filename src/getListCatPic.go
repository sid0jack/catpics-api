package main

import (
	"database/sql"
	"net/http"
)

// CatPicInfo represents the cat picture metadata.
type CatPicInfo struct {
	ID string `json:"id"`
	// If you have additional metadata such as the original filename, include it here.
}

// listCatPics godoc
// @Summary List all cat pictures
// @Description Get a list of all cat pictures' metadata
// @Tags catpics
// @Accept  json
// @Produce  json
// @Success 200 {array} CatPicInfo
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /catpics [get]
func listCatPics(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id FROM cat_pics")
		if err != nil {
			jsonError(w, "Server error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var pics []CatPicInfo
		for rows.Next() {
			var pic CatPicInfo
			if err := rows.Scan(&pic.ID); err != nil {
				jsonError(w, "Server error", http.StatusInternalServerError)
				return
			}
			pics = append(pics, pic)
		}

		// Check for errors from iterating over rows
		if err = rows.Err(); err != nil {
			jsonError(w, "Server error", http.StatusInternalServerError)
			return
		}

		jsonResponse(w, pics, http.StatusOK)
	}
}
