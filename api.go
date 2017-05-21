package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

//type apiHandler struct {}

type apiError struct {
	Tag     string `json:"-"`
	Error   error  `json:"-"`
	Message string `json:"error"`
	Code    int    `json:"code"`
}

type api struct {
	DB      *sql.DB
	Handler func(w http.ResponseWriter, r *http.Request, db *sql.DB) *apiError
}

func (api apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

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

func indexHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) *apiError {

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

func ExampleStripPrefix() {
	// To serve a directory on disk (/tmp) under an alternate URL path
	// (/tmpfiles/), use StripPrefix to modify the request URL's path before the
	// FileServer sees it:
	http.Handle("/tmpfiles/", http.StripPrefix("/tmpfiles/", http.FileServer(http.Dir("/tmp"))))
}
