package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"time"
)

func init() {
	log.SetFlags(log.LUTC | log.Ldate | log.Ltime)
}

type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *statusResponseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &statusResponseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)
		log.Printf("%s <- [%v %s] %s %s in %v\n",
			r.RemoteAddr,
			rw.status, http.StatusText(rw.status),
			r.Method, r.URL,
			time.Since(start),
		)
	})
}

func main() {
	http.Handle("/", helloHandler())
	server := http.Server{
		Addr:           ":8000",
		Handler:        logMiddleware(http.DefaultServeMux),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Println("Serving requests on http://127.0.0.1:8000")
	log.Fatal(server.ListenAndServe())
}

func helloHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})
}
