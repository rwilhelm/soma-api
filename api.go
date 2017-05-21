package main

import (
	"log"
	"net/http"
)

type Location struct {
	Latitude  float64 `json:latitude`
	Longitude float64 `json:longitude`
	Timestamp float64 `json:timestamp`
}

type AllLocationsResponse struct {
	Location []string
}

func (a *api) allLocations(w http.ResponseWriter, r *http.Request) {
	rows, err := a.db.Query("SELECT latitude, longitude, timestamp FROM location")

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {

		var (
			latitude  float64
			longitude float64
			timestamp string
		)

		err = rows.Scan(&latitude, &longitude, &timestamp)
		if err != nil {
			log.Fatal(err)
		}

		w.Write([]byte(timestamp))
	}
}
