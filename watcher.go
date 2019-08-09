package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/gorilla/websocket"
)

var watchChan chan string

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
		function sleep(ms) {
			return new Promise(resolve => setTimeout(resolve, ms));
		}
		(function() {
			var conn = new WebSocket("ws://127.0.0.1:8989/ws");
			conn.onclose = function(evt) {
				console.log('Connection closed');
			}
			conn.onmessage = async function(evt) {
				await sleep(300);
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
