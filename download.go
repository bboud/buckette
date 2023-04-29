package main

import (
	"fmt"
	"net/http"
)

func download(rw http.ResponseWriter, req *http.Request) {
	//URL routing ugh
	fmt.Println(req.URL)
	fmt.Println(req.RequestURI)
}
