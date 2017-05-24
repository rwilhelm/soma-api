package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

const (
	KEY    = "/etc/letsencrypt/live/soma.uni-koblenz.de/privkey.pem"
	CERT   = "/etc/letsencrypt/live/soma.uni-koblenz.de/cert.pem"
	DBUSER = "XXXUSERNAMEXXX"
	DBPASS = "XXXPASSWORDXXX"
	DBNAME = "soma_development"
	PORT   = "5000"
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

	// Get all locations
	http.HandleFunc("/location", BasicAuth(api.getAllLocations))

	// Get locations of user
	// curl https://soma.uni-koblenz.de:5000/location/cfcad754069f612a
	http.HandleFunc("/location/", api.getUserLocations)

	// Get all trips
	//http.HandleFunc("/trip", api.getAllTrips)

	// Get trips of user
	//http.HandleFunc("/trip/", api.getUserTrips)

	http.HandleFunc("/upload", api.uploadTrip)

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
