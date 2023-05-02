package main

import (
	"log"
	"net/http"
)

const (
	SiteFiles = "./frontend/dist"
)

func main() {
	LogPrint("Welcome to buckette version 0.0.1, thanks for stopin' by! 👋")
	fServer := NewFileServer()
	fServer.Initialize()

	LogPrint("Registering web server handlers 📝")

	// Root handler will parse download URL's and serve web content
	http.HandleFunc("/", fServer.HandleDownload)

	// API handlers
	http.HandleFunc("/upl", fServer.HandleUpload)
	http.HandleFunc("/ls/", fServer.HandleLS)

	LogSuccess("Web server is ready! Is that what ducks walk on? 🕸️ 🦆")
	log.Fatal(http.ListenAndServe(":8080", nil)) //Blocker
}
