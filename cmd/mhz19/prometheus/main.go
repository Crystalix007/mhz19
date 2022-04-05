package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/kebhr/mhz19"
)

type mhz_res struct {
	mut   sync.Mutex
	count int
	err   error
}

func main() {
	m := mhz19.MHZ19{}
	if err := m.Connect(); err != nil {
		fmt.Printf("{\"error\": \"%s\"}", err)
		os.Exit(1)
	}

	co2Value := mhz_res{
		mut: sync.Mutex{},
		count: 0,
		err: nil,
	}

	go updateCO2Metrics(m, &co2Value)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/metrics", http.StatusMovedPermanently)
	})
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) { co2Metrics(&co2Value, w, r) })
	log.Fatal(http.ListenAndServe(":9101", nil))
}

func updateCO2Metrics(m mhz19.MHZ19, co2Value *mhz_res) {
	for true {
		co2, err := m.ReadCO2()

		co2Value.mut.Lock()
		co2Value.count = co2
		co2Value.err = err
		co2Value.mut.Unlock()

		if err != nil {
			return
		}

		time.Sleep(2 * time.Second)
	}
}

func co2Metrics(co2Value *mhz_res, w http.ResponseWriter, r *http.Request) {
	co2Value.mut.Lock()
	defer co2Value.mut.Unlock()

	if co2Value.err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error \"%s\"\n", co2Value.err)
	} else {
		fmt.Fprintf(w, "co2 %d\n", co2Value.count)
	}
}
