package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/sid0jack/catpics-api/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

type CatPic struct {
	ID  string `json:"id"`
	Data []byte `json:"-"`
}

// @title Cat Pics API
// @version 1.0
// @description This is a simple set of API's to store and retrieve cat pictures.
// @host localhost:8080
// @BasePath /
func main() {
		db, err := sql.Open("sqlite3", "./catpics.sqlite3")
    if err != nil {
        log.Fatalf("Error opening database: %v", err)
    }
    defer db.Close()

    statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS cat_pics (id TEXT PRIMARY KEY, data BLOB NOT NULL);")
    if err != nil {
        log.Fatalf("Error preparing statement: %v", err)
    }
    defer statement.Close()

    _, err = statement.Exec()
    if err != nil {
        log.Fatalf("Error executing statement: %v", err)
    }

		router := mux.NewRouter()
		router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
		router.HandleFunc("/catpics", createCatPic(db)).Methods("POST")
		router.HandleFunc("/catpics/{id}", getCatPicByID(db)).Methods("GET")
		router.HandleFunc("/catpics/{id}", deleteCatPic(db)).Methods("DELETE")
		router.HandleFunc("/catpics", listCatPics(db)).Methods("GET")
		router.HandleFunc("/catpics/{id}", updateCatPic(db)).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8080", router))
}

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

type CatPicResponse struct {
	ID string `json:"id"`
}

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
