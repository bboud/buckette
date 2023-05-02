package main

import (
	"net/http"
	"strings"
)

func (fServer *FileServer) HandleLS(rw http.ResponseWriter, req *http.Request) {
	LogConnection(req)
	arguments := strings.Split(req.RequestURI, "/")
	arguments = arguments[2:]
}
