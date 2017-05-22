package main

import "net/http"

func BasicAuth(pass http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		username, password, _ := r.BasicAuth()

		if username != "XXXUSERNAMEXXX" || password != "bla" {
			http.Error(w, "authorization failed", http.StatusUnauthorized)
			return
		}
		pass(w, r)
	}
}
