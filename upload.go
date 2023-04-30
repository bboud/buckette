package main

import (
	"io"
	"log"
	"net/http"
	"strconv"
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

			fileName := req.Header.Get("File-Name")
			fileSize, err := strconv.ParseInt(req.Header.Get("File-Size"), 10, 64)
			if err != nil {
				LogWarning(
					"Unable to get file size",
					"Attempting to parse file size from header",
					err,
				)
			}
			//hash := req.Header.Get("File-Hash")

			// Handle the file saving in another goroutine
			f := &File{
				FileName:      fileName,
				Uploaded:      time.Now(),
				Size:          fileSize,
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
