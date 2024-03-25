package main

import (
	"database/sql"
	"log"
	"net/http"

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
			log.Fatal(err)
		}
		defer db.Close()

		statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS cat_pics (id TEXT PRIMARY KEY, data BLOB NOT NULL);")
		statement.Exec()

		router := mux.NewRouter()
		router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
		router.HandleFunc("/catpics", createCatPic(db)).Methods("POST")
		router.HandleFunc("/catpics/{id}", getCatPicByID(db)).Methods("GET")
		router.HandleFunc("/catpics/{id}", deleteCatPic(db)).Methods("DELETE")
		router.HandleFunc("/catpics", listCatPics(db)).Methods("GET")
		router.HandleFunc("/catpics/{id}", updateCatPic(db)).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8080", router))
}
