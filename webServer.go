package main

import (
	"log"
	"net/http"
)

func startServer() {
	http.Handle("/", http.FileServer(http.Dir("./html")))
	http.HandleFunc("/dld", download)
	http.HandleFunc("/upl", upload)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
