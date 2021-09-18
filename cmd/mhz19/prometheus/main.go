package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kebhr/mhz19"
)

func main() {
	m := mhz19.MHZ19{}
	if err := m.Connect(); err != nil {
		fmt.Printf("{\"error\": \"%s\"}", err)
		os.Exit(0)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, "/metrics", http.StatusMovedPermanently); });
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) { co2Metrics(m, w, r); });
	log.Fatal(http.ListenAndServe(":9101", nil))

}

func co2Metrics(m mhz19.MHZ19, w http.ResponseWriter, r *http.Request) {
	v, err := m.ReadCO2()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, fmt.Sprintf("error: \"%s\"", err))
	} else {
		fmt.Fprintf(w, fmt.Sprintf("co2: %d", v))
	}
}
