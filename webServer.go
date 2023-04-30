package main

import (
	"log"
	"net/http"
)

func startServer() {
	fs := newFileServer()
	go fs.start()
	fs.initialize()

	LogPrint("Registering web server handlers 📝")

	http.HandleFunc("/", root)
	http.HandleFunc("/dld", download)
	http.HandleFunc("/upl", fs.upload)

	LogSucess("Web server is ready for ducks! 🕸️ 🦆")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
