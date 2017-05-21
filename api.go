package main

import (
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

	//var uuidString = "13818EA5-27E9-429E-8C82-2E68DB2EA98D"
	//var uuidRegexp = "[A-F0-9]{8}-[A-F0-9]{4}-4[A-F0-9]{3}-[89aAbB][a-f0-9]{3}-[a-f0-9]{12}"

	uuid := strings.Split(r.URL.Path, "/")[2]

	//match, _ := regexp.MatchString(uuidString, uuidRegexp)
	//fmt.Println(match)

	//fmt.Println(reflect.TypeOf(uuid))

	////r, _ := regexp.Compile(uuidRegexp)

	log.Println(uuid)
	//fmt.Println(reflect.TypeOf(uuid))
	//fmt.Println(reflect.TypeOf(uuid[0]))

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
