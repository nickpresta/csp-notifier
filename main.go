package main

import (
	"net/http"
	"os"
	"runtime"
	"strconv"

	"labix.org/v2/mgo"

	"github.com/nickpresta/csp-notifier/handlers"
)

func setupMaxProcs() {
	envMaxProcs := os.Getenv("GOMAXPROCS")
	maxProcs, err := strconv.Atoi(envMaxProcs)
	if err != nil {
		maxProcs = runtime.NumCPU()
	}
	runtime.GOMAXPROCS(maxProcs)
}

func main() {
	setupMaxProcs()

	connect := os.Getenv("MONGO_CONNECT")
	if connect == "" {
		panic("Need MONGO_CONNECT environment variable defined in the form of: mongodb://myuser:mypass@localhost:40001/mydb")
	}

	session, err := mgo.Dial(connect)
	if err != nil {
		panic("Cannot connect to Mongo using: " + connect)
	}
	defer session.Close()

	http.HandleFunc("/", handlers.MainHandler)
	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		handlers.AddHandler(w, r, session)
	})
	http.HandleFunc("/reports", func(w http.ResponseWriter, r *http.Request) {
		handlers.ReportHandler(w, r, session)
	})
	if err := http.ListenAndServe(":9001", http.DefaultServeMux); err != nil {
		panic(err)
	}
}
