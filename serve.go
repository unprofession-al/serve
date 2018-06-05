package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/sontags/env"
	"github.com/unprofession-al/noip"
)

const (
	unset = "UNSET"
)

var port, logging, dir, user, pass, host, interval string

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if logging != "" {
			log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		}
		handler.ServeHTTP(w, r)
	})
}

func init() {
	env.Var(&port, "PORT", "8989", "Port that should be binded")
	env.Var(&logging, "LOG", "", "If not empty, log output will be written to STDOUT")
	env.Var(&dir, "DIR", ".", "Directory that should be served")
	env.Var(&user, "NOIP_USER", unset, "User to access no-ip")
	env.Var(&pass, "NOIP_PASS", unset, "Password to access no-ip")
	env.Var(&host, "NOIP_HOST", unset, "Hostname to update via no-ip")
	env.Var(&interval, "NOIP_INTERVAL", unset, "Interval to update no-ip")
}

func main() {
	env.Parse("S", false)

	if user != unset && pass != unset && host != unset && interval != unset {
		inter, err := strconv.Atoi(interval)
		if err == nil {
			cli := noip.New(user, pass, host, "", "sontags.serve/v1.0 daniel.menet")
			cli.Run(inter, true)
			log.Println("NO-IP is now managed...")
		}
	}

	log.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(":"+port, Log(http.FileServer(http.Dir(dir)))))
}
