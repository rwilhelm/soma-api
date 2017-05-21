package main

import (
	"database/sql"
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

type Location struct {
	Latitude  float64 `json:latitude`
	Longitude float64 `json:longitude`
	Timestamp float64 `json:timestamp`
}

type Response1 struct {
	Location []string
}
