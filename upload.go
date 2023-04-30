package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func (fServer *FileServer) upload(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	LogConnection(req)

	multipartReader, err := req.MultipartReader()
	if err != nil {
		log.Printf("ERROR: %v | Called by: %s", err, "MultipartReader")
		return
	}
	if strings.HasPrefix(req.Header.Get("Content-Type"), "multipart/") {
		for {
			p, err := multipartReader.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("ERROR: %v | Called by: %s", err, "NextPart")
				break
			}

			fmt.Println(p.Header)
			// Handle the file saving in another goroutine
			f := &File{
				FileName:      "testing",
				Uploaded:      time.Now(),
				DownloadCount: 0,
			}

			response, err := f.HandleUploadPart(p, fServer)
			// Don't need to handle here
			if err != nil {
				continue
			}
			rw.WriteHeader(200)
			rw.Write(response)
		}
	}
}
