package main

import (
	"log"
	"net/http"
)

func logger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
		if c.logging {
			log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		}
	})
}
