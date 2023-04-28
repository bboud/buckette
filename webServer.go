package main

import (
	"log"
	"net/http"
)

func startServer() {
	fileServer := NewFileServer()
	fileServer.Start()

	http.Handle("/", http.FileServer(http.Dir("./frontend/dist")))
	http.Handle("/fs", http.FileServer(http.Dir("./serve")))
	http.HandleFunc("/download/", fileServer.download)
	http.HandleFunc("/upl", fileServer.upload)
	http.HandleFunc("/ls", fileServer.ListFiles)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
