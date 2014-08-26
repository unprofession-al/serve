package main

import (
	"log"
	"net/http"

	"github.com/sontags/env"
)

var port, logging string

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if logging != "" {
			log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		}
		handler.ServeHTTP(w, r)
	})
}

func init() {
	env.Var(&port, "PORT", "8989", "Port that is binded")
	env.Var(&logging, "LOG", "", "If not empty, log output will be written to STDOUT")
}

func main() {
	env.Parse("S")
	log.Println("listening on port", port)
	log.Fatal(http.ListenAndServe(":"+port, Log(http.FileServer(http.Dir(".")))))
}
