package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	_ "github.com/lib/pq"
)

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp float64 `json:"timestamp"`
}

type Trip struct {
	TripUUID   string     `json:"uuid"`
	ClientUUID string     `json:"device_id"`
	Locations  []Location `json:"location"`
}

func (a *api) getAllLocations(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	rows, err := a.db.Query("SELECT latitude, longitude, timestamp FROM location LIMIT 42")

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var data []Location

	for rows.Next() {
		var (
			latitude  float64
			longitude float64
			timestamp float64
		)

		err := rows.Scan(&latitude, &longitude, &timestamp)
		if err != nil {
			log.Fatal(err)
		}

		data = append(data, Location{latitude, longitude, timestamp})
	}

	response, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(response)
}

func (a *api) getUserLocations(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	uuid := strings.Split(r.URL.Path, "/")[2]

	rows, err := a.db.Query(`SELECT latitude, longitude, timestamp FROM location JOIN trip ON location.trip_id = trip.id JOIN client ON trip.client_id = client.id WHERE client.device_id = $1`, uuid)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var data []Location

	for rows.Next() {
		var (
			latitude  float64
			longitude float64
			timestamp float64
		)

		err := rows.Scan(&latitude, &longitude, &timestamp)
		if err != nil {
			log.Fatal(err)
		}

		data = append(data, Location{latitude, longitude, timestamp})
	}

	if len(data) == 0 {
		// https://tools.ietf.org/html/rfc7231#section-6.3.5
		w.WriteHeader(http.StatusNotFound)
	} else {
		response, err := json.Marshal(data)
		w.Write(response)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (a *api) uploadTrip(w http.ResponseWriter, r *http.Request) {

	var (
		clientID int
		tripID   int
		d        Trip
	)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&d)
	handleError(err)
	defer r.Body.Close()

	err = a.db.QueryRow(`
	SELECT id FROM client WHERE device_id = $1`, d.ClientUUID).Scan(&clientID)
	if err != nil {
		err = a.db.QueryRow(`
		INSERT INTO client (device_id) VALUES ($1) RETURNING id`, d.ClientUUID).Scan(&clientID)
		handleError(err)
	}

	err = a.db.QueryRow(`
	INSERT INTO trip (uuid, client_id) VALUES ($1, $2) RETURNING id`, d.TripUUID, clientID).Scan(&tripID)
	handleError(err)

	log.Printf("%s -> %d -> %d", d.TripUUID, clientID, tripID)

	// rows, err := a.db.Query(`
	// INSERT INTO location (id, trip_id, accuracy, altitude, bearing, latitude, longitude, timestamp, speed)
	// VALUES (DEFAULT, $1, $2, $3, $4, $5, $6, $7, $8`)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// defer rows.Close()

	// for rows.Next() {
	// 	var (
	// 		latitude  float64
	// 		longitude float64
	// 		timestamp float64
	// 	)

	// 	err := rows.Scan(&latitude, &longitude, &timestamp)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	data = append(data, Location{latitude, longitude, timestamp})
	// }

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	//if len(data) == 0 {
	//	// https://tools.ietf.org/html/rfc7231#section-6.3.5
	//	w.WriteHeader(http.StatusNotFound)
	//} else {
	//	response, err := json.Marshal(data)
	//	w.Write(response)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}
}

func eatRows(rows *sql.Rows) (d Trip) {
	for rows.Next() {
		err := rows.Scan(&d)
		handleError(err)
	}
	return d
}

func countRows(rows *sql.Rows) (count int) {
	for rows.Next() {
		err := rows.Scan(&count)
		handleError(err)
	}
	return count
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
