package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"runtime/debug"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/", sourceCodeHandler)
	mux.HandleFunc("/panic/", panicDemo)
	mux.HandleFunc("/panic-after/", panicAfterDemo)
	mux.HandleFunc("/", hello)
	log.Fatal(http.ListenAndServe(":3000", recoverMx(mux, true)))
}

func sourceCodeHandler(w http.ResponseWriter, r *http.Request) {
	path := r.FormValue("path")
	file, err := os.Open(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	io.Copy(w, file)
}

func recoverMx(app http.Handler, dev bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				stack := debug.Stack()
				log.Println(string(stack))
				w.WriteHeader(http.StatusInternalServerError)
				if dev {
					fmt.Fprintf(w, "<h1>panic: %s</h1><pre>%s</pre>", err, string(stack))
				} else {
					http.Error(w, "Something went wrong :(", http.StatusInternalServerError)
				}
			}
		}()

		nw := &responseWriter{ResponseWriter: w}
		app.ServeHTTP(nw, r)
		nw.flush()
	}
}

// type ResponseWriter interface {
// 	Header() Header
// 	Write([]byte) (int, error)
// 	WriteHeader(statusCode int)
// }

type responseWriter struct {
	http.ResponseWriter
	writes [][]byte
	status int
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.writes = append(rw.writes, b)
	return len(b), nil
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
}

func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := rw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("the ResponseWriter dose not support the Hijacker interface")
	}
	return hijacker.Hijack()
}

func (rw *responseWriter) Flush() {
	flusher, ok := rw.ResponseWriter.(http.Flusher)
	if !ok {
		return
	}
	flusher.Flush()
}

func (rw *responseWriter) flush() error {
	if rw.status != 0 {
		rw.ResponseWriter.WriteHeader(rw.status)
	}
	for _, w := range rw.writes {
		_, err := rw.ResponseWriter.Write(w)
		if err != nil {
			return err
		}
	}
	return nil
}

func panicDemo(w http.ResponseWriter, r *http.Request) {
	funcThatPanics()
}

func panicAfterDemo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello!</h1>")
	funcThatPanics()
}

func funcThatPanics() {
	panic("Oh no!")
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Hello!</h1>")
}
