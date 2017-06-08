package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
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

	http.HandleFunc("/upload", api.uploadTrip)

	// Get all locations
	http.HandleFunc("/location", BasicAuth(api.getAllLocations))

	// Get all or a user's locations
	// curl https://soma.uni-koblenz.de:5000/location/:id
	http.HandleFunc("/location/", BasicAuth(api.getUserLocations))

	// Generate a new LimeSurvey token for ...?q=<id>
	http.HandleFunc("/token/new", BasicAuth(api.generateToken))

	log.Printf("Listening on :%s", PORT)
	http.ListenAndServeTLS(":"+PORT, CERT, KEY, nil)
	//http.ListenAndServe(":"+PORT, nil)
}

func logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s requested %s", r.RemoteAddr, r.URL)
		h.ServeHTTP(w, r)
	})
}
