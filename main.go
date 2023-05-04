package main

import (
	"log"
	"net/http"

	"github.com/bboud/buckette/fileserver"
	"github.com/bboud/buckette/logger"
)

func main() {
	logger.LogPrint("Welcome to buckette version 0.0.1, thanks for stopin' by! 👋")
	fServer := fileserver.NewFileServer()

	logger.LogPrint("Registering web server handlers 📝")

	// Root handler will parse download URL's and serve web content
	http.HandleFunc("/", fServer.HandleDownload)

	// API handlers
	http.HandleFunc("/upl", fServer.HandleUpload)
	http.HandleFunc("/ls/", fServer.HandleLS)

	logger.LogSuccess("Web server is ready! Is that what ducks walk on? 🕸️ 🦆")
	log.Fatal(http.ListenAndServe(":8080", nil)) //Blocker
}
