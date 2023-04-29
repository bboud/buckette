package main

import (
	"io"
	"log"
	"net/http"
	"strings"
)

// func processFile(p *multipart.Part) {
// 	// Needs to be buffered
// 	slurp, err := io.ReadAll(p)
// 	if err != nil {
// 		log.Printf("ERROR: %v | Called by: %s", err, "ReadAll(Slurp)")
// 		return
// 	}

// 	fmt.Printf("#####Header#####\n %s \n\n", p.Header)
// 	fmt.Printf("#####Part#####\n %s \n\n", slurp)
// }

func upload(rw http.ResponseWriter, req *http.Request) {

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
				FileName: "testing",
			}
			if err = f.Push(p); err != nil {
				log.Fatal(err)
			}
		}
	}
	rw.WriteHeader(200)
}
