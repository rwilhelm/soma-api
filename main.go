package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

const (
	DBUSER = "asdf"
	DBPASS = "qwer"
	DBNAME = "soma_development"
	PORT   = "3333"
)

type api struct {
	db *sql.DB
}

func main() {
	dbinfo := fmt.Sprintf(
		"user=%s dbname=%s password=%s sslmode=disable", DBUSER, DBNAME, DBPASS)

	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("[ERROR] Could not establish a connection with the database")
		log.Fatal(err)
	}

	defer db.Close()

	api := api{db: db}

	http.HandleFunc("/loc", api.allLocations)

	log.Printf("Listening on :%s", PORT)
	http.ListenAndServe(":"+PORT, nil)
}

func logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s requested %s", r.RemoteAddr, r.URL)
		h.ServeHTTP(w, r)
	})
}
