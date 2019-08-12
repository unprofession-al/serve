package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/gorilla/websocket"
	"github.com/radovskyb/watcher"
)

var watchChan chan string

type InjectorMiddleware struct{}

func (m *InjectorMiddleware) Wrap(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/ws") {
			rec := httptest.NewRecorder()
			next.ServeHTTP(rec, r)

			for k, v := range rec.Header() {
				w.Header()[k] = v
			}
			out := rec.Body.Bytes()

			if strings.HasSuffix(strings.ToLower(r.RequestURI), ".md") {
				extensions := parser.CommonExtensions | parser.AutoHeadingIDs
				parser := parser.NewWithExtensions(extensions)

				htmlFlags := html.CommonFlags | html.HrefTargetBlank
				opts := html.RendererOptions{Flags: htmlFlags}
				renderer := html.NewRenderer(opts)

				data := htmlData{
					Body:  string(markdown.ToHTML(out, parser, renderer)),
					Title: r.URL.Path,
				}

				tmpl := template.Must(template.New("page").Parse(htmlScaffold))
				var rendered bytes.Buffer
				err := tmpl.Execute(&rendered, data)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("template cound not be rendered"))
					return
				}

				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				out = rendered.Bytes()
			}

			out = bytes.Replace(out, []byte("</head>"), []byte(`    <script>
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
		} else {
			next.ServeHTTP(w, r)
		}
	}

	return http.HandlerFunc(fn)
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

func watchFiles(dir string) {
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
