package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	//"strings"
	"time"

	_ "github.com/lib/pq"
)

const (
	DBUSER = "XXXUSERNAMEXXX"
	DBPASS = "XXXPASSWORDXXX"
	DBNAME = "soma_development"
	PORT   = "3333"
)

func msToTime(ms string) (time.Time, error) {
	msInt, err := strconv.ParseInt(ms, 8, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(0, msInt*int64(time.Millisecond)), nil
}

func main() {

	dbinfo := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", DBUSER, DBNAME, DBPASS)
	db, err := sql.Open("postgres", dbinfo)

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()

	if err != nil {
		log.Fatal("[ERROR] Could not establish a connection with the database")
		log.Fatal(err)
	}

	// The sql.DB should not have a lifetime beyond the scope of the function.
	defer db.Close()

	rows, err := db.Query("SELECT latitude, longitude, timestamp FROM location")

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

		fmt.Printf("%3v | %3v | %3v\n", latitude, longitude, timestamp)
	}

	http.Handle("/", ApiHandler{db, indexHandler})
	log.Printf("Listening on :%s", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}
