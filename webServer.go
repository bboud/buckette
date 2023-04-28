package main

import (
	"log"
	"net/http"
)

func startServer() {
	fileServer := NewFileServer()
	fileServer.Start()

	http.Handle("/", http.FileServer(http.Dir("./html")))
	http.HandleFunc("/f/", fileServer.download)
	http.HandleFunc("/upload", fileServer.upload)
	http.HandleFunc("/ls", fileServer.ListFiles)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
