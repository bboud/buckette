package main

import (
	"log"
	"net/http"
)

func startServer() {
	fServer := newFileServer()
	go fServer.handleNewFiles()
	fServer.initialize()

	LogPrint("Registering web server handlers 📝")

	// Root handler will parse download URL's and serve web content
	http.HandleFunc("/", fServer.download)

	// API handlers
	http.HandleFunc("/upl", fServer.upload)
	http.HandleFunc("/ls/", fServer.ls)

	LogSucess("Web server is ready! Is that what ducks walk on? 🕸️ 🦆")
	log.Fatal(http.ListenAndServe(":8080", nil)) //Blocker
}
