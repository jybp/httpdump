package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/http/httputil"
	"strings"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%+v", err)
	}
}

func run() error {
	var port int
	var multipartform bool
	flag.IntVar(&port, "p", 8080, "Port")
	flag.BoolVar(&multipartform, "multipart", false, "Attempt to parse multipart/form-data requests")
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if multipartform && strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
			multipartFormHandler(w, r)
			return
		}
		dumpHandler(true)(w, r)
	})
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func dumpHandler(body bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := httputil.DumpRequest(r, body)
		if err != nil {
			errHandler(err)(w, r)
			return
		}
		log.Printf("%s", string(b))
	}
}

func multipartFormHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(math.MaxInt64 - int64(10<<20)); err != nil {
		errHandler(err)(w, r)
		return
	}
	dumpHandler(false)(w, r)
	fmt.Printf("Parsed multipart form:\n")
	for k, values := range r.PostForm {
		fmt.Printf("--- key: %s\n", k)
		for _, v := range values {
			fmt.Printf("--- value:\n%s\n", v)
		}
	}
}

func errHandler(err error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Error: %+v", err)
		log.Printf("Error: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
