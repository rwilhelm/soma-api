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

type Location struct {
	Latitude  float64 `json:latitude`
	Longitude float64 `json:longitude`
	Timestamp float64 `json:timestamp`
}

type Response1 struct {
	Location []string
}

type api struct {
	db *sql.DB
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

	defer db.Close()

	api := api{db: db}

	http.HandleFunc("/hello", api.handleHello)

	//http.handleFunc("/upload", s.handleUpload)
	//http.handle

	//hello := helloWorldHandler()
	//http.Handle("/hello", hello)

	//allLocations := allLocationsHandler()
	//http.Handle("/api/v2/loc", allLocations)

	log.Printf("Listening on :%s", PORT)
	http.ListenAndServe(":"+PORT, nil)

	//if err != nil && err != http.ErrServerClosed {
	//	log.Fatal(err)
	//}
}

func (a *api) handleHello(w http.ResponseWriter, r *http.Request) {
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

func helloWorldHandler() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	}
	return http.HandlerFunc(fn)
}

//func dbTestHandler() http.Handler {
//	fn := func(w http.ResponseWriter, r *http.Request, db *sql.DB) http.Handler {
//		err := db.Ping()
//		if err != nil {
//			log.Fatal("DB ERROR")
//		}
//		log.Println("DB SUCCESS")
//		w.Write([]byte("hello world"))
//	}
//	return http.HandlerFunc(fn)
//}

func allLocationsHandler() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

	}
	return http.HandlerFunc(fn)
}
