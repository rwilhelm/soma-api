package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"time"
)

const (
	DB_USER = "XXXUSERNAMEXXX"
	DB_PASS = "XXXPASSWORDXXX"
	DB_NAME = "soma_development"
	PORT    = "3333"
)

var badTimestamp = regexp.MustCompile(`^[0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9].[0-9][0-9][0-9][0-9][0-9][0-9]$`)
var bts = "1495230019.084993"

func msToTime(ms string) (time.Time, error) {
	msInt, err := strconv.ParseInt(ms, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(0, msInt*int64(time.Millisecond)), nil
}

func main() {

	dbinfo := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", DB_USER, DB_NAME, DB_PASS)
	db, err := sql.Open("postgres", dbinfo)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	rows, err := db.Query("SELECT latitude, longitude, timestamp FROM location")

	if err != nil {
		panic(err)
	}

	fmt.Println(reflect.TypeOf(dbinfo)) // => string
	fmt.Println(reflect.TypeOf(db))     // => *sql.DB
	fmt.Println(reflect.TypeOf(rows))   // => *sql.Rows

	for rows.Next() {

		var latitude float64
		var longitude float64
		var timestamp float64

		err = rows.Scan(&latitude, &longitude, &timestamp)
		if err != nil {
			panic(err)
		}

		if timestamp > 9999999999.999999 {
			timestamp = timestamp / 1000
		}

		fmt.Printf("%3v | %3v | %3f\n", latitude, longitude, timestamp)
	}

	http.Handle("/", ApiHandler{db, indexHandler})
	log.Printf("Listening on :%s", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}
