package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/spf13/pflag"
	"github.com/unprofession-al/noip"
)

const (
	unset = "UNSET"
)

func main() {
	pflag.Parse()

	if c.noip.user != "" && c.noip.pass != "" && c.noip.host != "" {
		cli := noip.New(c.noip.user, c.noip.pass, c.noip.host, "", "sontags.serve/v1.0 daniel.menet")
		cli.Run(c.noip.interval, true)
		log.Println("NO-IP is now managed...")
	}

	if !c.noWatch {
		watchFiles(c.dir)
	}

	log.Printf("Listening at http://%s\n", c.listener)

	r := mux.NewRouter().StrictSlash(true)
	r.Path("/ws").HandlerFunc(serveWs)
	r.PathPrefix("/render").HandlerFunc(MarkdownHandler)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(c.dir)))
	r.Path("/ws").HandlerFunc(serveWs)
	in := InjectorMiddleware{}
	chain := alice.New(logger, in.Wrap).Then(r)
	log.Fatal(http.ListenAndServe(c.listener, chain))
}
