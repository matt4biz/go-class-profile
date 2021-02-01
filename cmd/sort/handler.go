package main

import (
	"log"
	"net/http"
	"strconv"
)

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func getLoop(r *http.Request) int {
	loop := 1

	if i, err := strconv.Atoi(r.FormValue("loop")); err == nil {
		loop = i - 1
	}

	return loop
}

func getDelay(r *http.Request) int {
	delay := 8

	if i, err := strconv.Atoi(r.FormValue("delay")); err == nil {
		delay = i
	}

	return delay
}
