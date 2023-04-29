package main

import (
	"log"
	"net/http"
)

func startServer() {
	fs := newFileServer()
	go fs.start()

	http.Handle("/", http.FileServer(http.Dir("./html")))
	http.HandleFunc("/dld", download)
	http.HandleFunc("/upl", fs.upload)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
