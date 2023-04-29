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

func (f *FileServer) upload(rw http.ResponseWriter, req *http.Request) {
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

	// outFile, err := os.OpenFile("mine.txt", os.O_CREATE|os.O_RDWR, 0644)
	// if err != nil {
	// 	log.Printf("ERROR: %v | Called by: %s", err, "OpenFile")
	// 	return
	// }
	// defer outFile.Close()

	// //stride := 4098
	// content, err := ioutil.ReadAll(req.Body)
	// if err != nil {
	// 	log.Printf("ERROR: %v | Called by: %s", err, "ReadAll")
	// }
	// defer req.Body.Close()

	// writer := bufio.NewWriter(outFile)
	// writer.Write(content)

	// filename := "test"

	// file := File{
	// 	FileName:      filename,
	// 	UUID:          hash(filename),
	// 	Size:          req.ContentLength,
	// 	Uploaded:      time.Now(),
	// 	DownloadCount: 0,
	// }

	// f.addChannel <- &file
}
