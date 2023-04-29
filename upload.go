package main

import (
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"
)

func upload(rw http.ResponseWriter, req *http.Request) {
	mediaType, params, err := mime.ParseMediaType(req.Header.Get("Content-Type"))
	if err != nil {
		log.Fatal(err)
	}
	if strings.HasPrefix(mediaType, "multipart/") {
		mr := multipart.NewReader(req.Body, params["boundary"])
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("ERROR: %v | Called by: %s", err, "NextPart")
				break
			}
			slurp, err := io.ReadAll(p)
			if err != nil {
				log.Printf("ERROR: %v | Called by: %s", err, "ReadAll(Slurp)")
				break
			}
			fmt.Printf("#####Header#####\n %s \n\n", p.Header)
			fmt.Printf("#####Part#####\n %s \n\n", slurp)
		}
	}
}
