package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/radovskyb/watcher"
	"github.com/sontags/env"
	"github.com/unprofession-al/noip"
)

const (
	unset = "UNSET"
)

var listener, logging, dir, user, pass, host, interval, watch string

var watchChan chan string

func logger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
		if logging != unset {
			log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		}
	})
}

type InjectorMiddleware struct {
	handler http.Handler
}

func (m *InjectorMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rec := httptest.NewRecorder()
	m.handler.ServeHTTP(rec, r)
	for k, v := range rec.Header() {
		w.Header()[k] = v
	}
	out := bytes.Replace(rec.Body.Bytes(), []byte("</head>"), []byte(`    <script>
		(function() {
			var conn = new WebSocket("ws://127.0.0.1:8989/ws");
			conn.onclose = function(evt) {
				console.log('Connection closed');
			}
			conn.onmessage = function(evt) {
				console.log(evt);
				conn.close();
				location.reload(true);
			}
		})();
	</script></head>`), -1)
	w.Header().Set("Content-Length", strconv.Itoa(len(out)))
	w.Write(out)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}

	go func() {
		for {
			desc, ok := <-watchChan

			if !ok {
				return
			}
			if err := ws.WriteMessage(websocket.TextMessage, []byte(desc)); err != nil {
				return
			}

		}
	}()
}

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
