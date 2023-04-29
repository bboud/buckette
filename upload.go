package main

import (
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// type File struct {
// 	Size          int64
// 	ContentType   string
// 	UserUploaded  string
// }

func (fs *FileServer) upload(rw http.ResponseWriter, req *http.Request) {
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

			// Handle the file saving in another goroutine
			f := &File{
				FileName:      "testing",
				Uploaded:      time.Now(),
				DownloadCount: 0,
			}

			response, err := f.Push(p)
			if err != nil {
				log.Println(err)
			}
			rw.WriteHeader(200)
			rw.Write(response)

			fs.push(f)
		}
	}
}
