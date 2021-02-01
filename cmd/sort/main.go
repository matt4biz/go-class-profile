package main

import (
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	var speed string

	flag.StringVar(&speed, "speed", "slow", "painting speed")
	flag.Parse()

	switch speed {
	case "faster":
		paintSquare = paintSquareFast
	case "fastest":
		paintSquare = paintSquareFastest
	case "slow":
		// nop
	default:
		log.Fatal("unknown speed:", speed)
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "8081"
	}

	log.Printf("Speed %q", speed)
	log.Printf("Listening on port %s", port)

	router := mux.NewRouter()

	router.Use(logMiddleware)
	router.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)

	router.HandleFunc("/insert", insertHandler).Methods(http.MethodGet)
	router.HandleFunc("/qsort", qsortHigh).Methods(http.MethodGet)
	router.HandleFunc("/qsortm", qsortMiddle).Methods(http.MethodGet)
	router.HandleFunc("/qsort3", qsortMedian).Methods(http.MethodGet)
	router.HandleFunc("/qsorti", qsortInsert).Methods(http.MethodGet)
	router.HandleFunc("/qsortf", qsortFlag).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":"+port, router))
}
