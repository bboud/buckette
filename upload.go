package main

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func (f *FileServer) upload(rw http.ResponseWriter, req *http.Request) {
	inFile, header, err := req.FormFile("file")
	if err == http.ErrNotMultipart || err != nil {
		log.Println(err)
		return
	}
	defer inFile.Close()

	outFile, err := os.OpenFile(header.Filename, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Println(err)
		return
	}
	defer outFile.Close()

	stride := 4098
	reader := bufio.NewReader(inFile)
	writer := bufio.NewWriter(outFile)
	buffer := make([]byte, 0, stride)

	for {
		n, err := io.ReadFull(reader, buffer[:cap(buffer)])
		buffer = buffer[:n]
		writer.Write(buffer)
		if err == io.EOF {
			break
		}
		if err == io.ErrUnexpectedEOF {
			log.Fatal(err)
		}
	}

	filename := header.Filename

	file := File{
		FileName:      filename,
		UUID:          hash(filename),
		Size:          header.Size,
		Uploaded:      time.Now(),
		DownloadCount: 0,
	}

	f.addChannel <- &file
}
