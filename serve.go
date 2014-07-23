package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8989"
	}

	log.Println("listening on port", port)

	log.Fatal(http.ListenAndServe(":"+port, http.FileServer(http.Dir("."))))
}
