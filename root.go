package main

import "net/http"

func root(rw http.ResponseWriter, req *http.Request) {
	LogConnection(req)
	http.FileServer(http.Dir("./frontend/dist")).ServeHTTP(rw, req)
}
