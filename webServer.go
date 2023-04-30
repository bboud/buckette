package main

import (
	"log"
	"net/http"
)

func startServer() {
	fServer := newFileServer()
	go fServer.start()
	fServer.initialize()

	LogPrint("Registering web server handlers 📝")

	http.HandleFunc("/", fServer.router)
	http.HandleFunc("/upl", fServer.upload)

	LogSucess("Web server is ready for ducks! 🕸️ 🦆")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
