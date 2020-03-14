package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	http.HandleFunc("/", errorHandler("./pages"))

	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.ListenAndServe(fmt.Sprintf(":8080"), nil)
}

func errorHandler(path string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		errCode := r.Header.Get("X-Code")
		code, err := strconv.Atoi(errCode)
		if err != nil {
			code = 404
		}
		w.WriteHeader(code)

		scode := strconv.Itoa(code)
		file := fmt.Sprintf("%v/%cxx%v", path, scode[0], ".html")
		f, err := os.Open(file)
		if err != nil {
			log.Printf("Unexpected error opening file: %v", err)
			http.NotFound(w, r)
			return
		}
		defer f.Close()
		io.Copy(w, f)
	}
}
