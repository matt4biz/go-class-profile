package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const url = "https://jsonplaceholder.typicode.com"

type todo struct {
	UserID    int    `json:"userID"`
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var mark = map[bool]string{
	false: " ",
	true:  "x",
}

func handler(w http.ResponseWriter, r *http.Request) {
	// we do it this way so we can ensure we've created the leak
	// because we're not using the default client with pooling

	req, _ := http.NewRequest("GET", url+"/todos/"+r.URL.Path[1:], nil)
	tr := &http.Transport{}
	cli := &http.Client{Transport: tr}
	resp, err := cli.Do(req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
	}

	// right here we need defer resp.Body.Close()
	// without which we will leak goroutines & fds

	if resp.StatusCode != http.StatusOK {
		http.NotFound(w, r)
		return
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var item todo

	if err := json.Unmarshal(body, &item); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "[%s] %d - %s\n", mark[item.Completed], item.ID, item.Title)
	queries.Inc()
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	// we don't add pprof; it's done for us automatically

	http.HandleFunc("/", handler)
	http.Handle("/metrics", promhttp.Handler())

	log.Printf("Listening on port %s", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

var queries = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "all_queries",
	Help: "How many queries we've received.",
})

func init() {
	prometheus.MustRegister(queries)
}
