package main

import (
	"log"
	"net/http"

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

	if c.watch {
		watchFiles(c.dir)
	}

	log.Printf("Listening at http://%s\n", c.listener)

	http.Handle("/", logger(&InjectorMiddleware{http.FileServer(http.Dir(c.dir))}))
	http.HandleFunc("/ws", serveWs)
	log.Fatal(http.ListenAndServe(c.listener, nil))
}
