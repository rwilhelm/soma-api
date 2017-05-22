package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"unsafe"

	_ "github.com/lib/pq"
)

type Location struct {
	Accuracy  float64 `json:"accuracy,omitempty"`
	Altitude  float64 `json:"altitude,omitempty"`
	Bearing   float64 `json:"bearing,omitempty"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Speed     float64 `json:"speed,omitempty"`
	Timestamp float64 `json:"timestamp"`
}

type Trip struct {
	TripUUID   string     `json:"uuid"`
	ClientUUID string     `json:"device_id"`
	Locations  []Location `json:"locationData"`
}

func (a *api) getAllLocations(w http.ResponseWriter, r *http.Request) { // {{{

	var (
		data []Location
		rows *sql.Rows
	)

	limit := r.URL.Query().Get("l")

	if limit != "" {
		r, err := a.db.Query("SELECT latitude, longitude, timestamp FROM location LIMIT $1", limit)
		handleError(err)
		defer r.Close()
		rows = r
	} else {
		r, err := a.db.Query("SELECT latitude, longitude, timestamp FROM location")
		handleError(err)
		defer r.Close()
		rows = r
	}

	for rows.Next() {
		var (
			latitude  float64
			longitude float64
			timestamp float64
		)

		err := rows.Scan(&latitude, &longitude, &timestamp)

		handleError(err)

		//log.Println(fmt.Printf("%14.4f\n", timestamp))
		//t, _ := fmt.Printf("%14.4f", timestamp)
		//log.Println(t)

		data = append(data, Location{Latitude: latitude, Longitude: longitude, Timestamp: timestamp})
	}

	//d := json.NewDecoder(strings.NewReader(data))
	//d.UseNumber()
	//err := d.Decode(&data)
	//handleError(err)

	// FIXME timestamp must be a number

	response, err := json.Marshal(data)
	handleError(err)

	//log.Printf("[OUTGOING] Sending %d locations (%d)", len(data), unsafe.Sizeof(response))

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Write(response)
} // }}}

func (a *api) getUserLocations(w http.ResponseWriter, r *http.Request) { // {{{
	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	uuid := strings.Split(r.URL.Path, "/")[2]

	rows, err := a.db.Query(`SELECT latitude, longitude, timestamp FROM location JOIN trip ON location.trip_id = trip.id JOIN client ON trip.client_id = client.id WHERE client.device_id = $1`, uuid)

	handleError(err)

	defer rows.Close()

	var data []Location

	for rows.Next() {
		var (
			accuracy  float64
			altitude  float64
			bearing   float64
			latitude  float64
			longitude float64
			speed     float64
			timestamp float64
		)

		err := rows.Scan(&accuracy, &altitude, &bearing, &latitude, &longitude, &speed, &timestamp)
		handleError(err)

		data = append(data, Location{accuracy, altitude, bearing, latitude, longitude, speed, timestamp})
	}

	if len(data) == 0 {
		// https://tools.ietf.org/html/rfc7231#section-6.3.5
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response, err := json.Marshal(data)
	handleError(err)
	w.Write(response)
} // }}}

func (a *api) uploadTrip(w http.ResponseWriter, r *http.Request) { // {{{

	var (
		clientID int
		tripID   int
		d        Trip
	)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&d)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}

	defer r.Body.Close()

	if len(d.ClientUUID) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("[BAD REQUEST] ClientUUID undefined")
		return
	}

	if len(d.TripUUID) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("[BAD REQUEST] TripUUID undefined")
		return
	}

	if len(d.Locations) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("[BAD REQUEST] Locations missing")
		return
	}

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

	//log.Printf("%s -> %d -> %d", d.TripUUID, clientID, tripID)

	var locationsInserted int64
	for _, l := range d.Locations {
		stmt, err := a.db.Prepare("INSERT INTO location (id, trip_id, accuracy, altitude, bearing, latitude, longitude, timestamp, speed) VALUES (DEFAULT, $1, $2, $3, $4, $5, $6, $7, $8)")
		handleError(err)

		res, err := stmt.Exec(tripID, l.Accuracy, l.Altitude, l.Bearing, l.Latitude, l.Longitude, l.Timestamp, l.Speed)
		handleError(err)

		//lastId, err := res.LastInsertId()
		//handleError(err)

		rowCnt, err := res.RowsAffected()
		handleError(err)

		locationsInserted += rowCnt
	}

	log.Printf("[INCOMING] c:%d/%s t:%d/%s l:%d", clientID, d.ClientUUID, tripID, d.TripUUID, locationsInserted)
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
} // }}}

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
		log.Println(err)
	}
}

func fatalError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
