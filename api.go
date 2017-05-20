package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type apiError struct {
	Tag     string `json:"-"`
	Error   error  `json:"-"`
	Message string `json:"error"`
	Code    int    `json:"code"`
}

type ApiHandler struct {
	DB      *sql.DB
	Handler func(w http.ResponseWriter, r *http.Request, db *sql.DB) *apiError
}

func (api ApiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Headers for all responses
	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	// Initialize API handler
	err := api.Handler(w, r, api.DB)
	if err != nil {

		log.Printf("AAA %s %s %s [%s] %s", r.RemoteAddr, r.Method, r.URL, err.Tag, err.Error)

		w.WriteHeader(err.Code)

		resp := json.NewEncoder(w)
		err_json := resp.Encode(err)
		if err_json != nil {
			log.Println("Encode JSON for error response was failed.")
			return
		}
		return
	}

	log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
}

type Location struct {
	Latitude  float64 `json:latitude`
	Longitude float64 `json:longitude`
	Timestamp float64 `json:timestamp`
}

type Response1 struct {
	Location []string
}

func indexHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) *apiError {
	if r.URL.Path != "/" {
		return &apiError{
			"indexHandler url",
			errors.New("Not found!"),
			"Not found",
			http.StatusNotFound,
		}
	}

	err := db.Ping()
	if err != nil {
		return &apiError{
			"indexHandler db ping",
			err,
			"internal server error",
			http.StatusInternalServerError,
		}
	}

	return nil
}
