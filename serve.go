package main

import (
	"log"
	"net/http"
	"os"
)

func Log(handler http.Handler, logging string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if logging != "" {
			log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		}
		handler.ServeHTTP(w, r)
	})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8989"
	}
	logging := os.Getenv("LOG")
	log.Println("listening on port", port)
	log.Fatal(http.ListenAndServe(":"+port, Log(http.FileServer(http.Dir(".")), logging)))
}
