package main

import (
	"net/http"
)

func BasicAuth(pass http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		username, password, _ := r.BasicAuth()

		const realm = ""

		if username != BAUSER || password != BAPASS {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			http.Error(w, "authorization failed", http.StatusUnauthorized)
			return
		}

		pass(w, r)
	}
}
