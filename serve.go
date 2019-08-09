package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/radovskyb/watcher"
	"github.com/sontags/env"
	"github.com/unprofession-al/noip"
)

const (
	unset = "UNSET"
)

var listener, logging, dir, user, pass, host, interval, watch string

func init() {
	env.Var(&listener, "LISTENER", "127.0.0.1:8989", "Listener that should be binded")
	env.Var(&logging, "LOG", unset, "If not empty, log output will be written to STDOUT")
	env.Var(&dir, "DIR", ".", "Directory that should be served")
	env.Var(&user, "NOIP_USER", unset, "User to access no-ip")
	env.Var(&pass, "NOIP_PASS", unset, "Password to access no-ip")
	env.Var(&host, "NOIP_HOST", unset, "Hostname to update via no-ip")
	env.Var(&interval, "NOIP_INTERVAL", unset, "Interval to update no-ip")
	env.Var(&watch, "WATCH", unset, "Watch and refresh")
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

	if watch != unset {
		watchChan = make(chan string)
		w := watcher.New()

		go func() {
			for {
				select {
				case event := <-w.Event:
					watchChan <- event.String()
					// fmt.Println(event) // Print the event's info.
				case err := <-w.Error:
					log.Fatal(err)
				case <-w.Closed:
					return
				}
			}
		}()

		if err := w.AddRecursive(dir); err != nil {
			log.Fatalln(err)
		}

		go func() {
			if err := w.Start(time.Millisecond * 100); err != nil {
				log.Fatalln(err)
			}
		}()
	}

	log.Printf("Listening at http://%s\n", listener)

	http.Handle("/", logger(&InjectorMiddleware{http.FileServer(http.Dir(dir))}))
	http.HandleFunc("/ws", serveWs)
	log.Fatal(http.ListenAndServe(listener, nil))
}
