package main

import (
	"crypto/sha256"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const contentFolder = "./serve/"

// const dbName = "db.gob"

func hash(fileName string) []byte {
	hash := sha256.New()
	f, err := os.Open(contentFolder + fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if _, err = io.Copy(hash, f); err != nil {
		log.Fatal(err)
	}

	return hash.Sum(nil)
}

type File struct {
	FileName      string
	UUID          []byte
	Size          int64
	Uploaded      time.Time
	UserUploaded  string
	DownloadCount int
}

type FileServer struct {
	files         []File
	addChannel    chan *File
	deleteChannel chan *File
}

func NewFileServer() *FileServer {
	fs := &FileServer{
		addChannel:    make(chan *File, 10),
		deleteChannel: make(chan *File, 10),
	}

	fs.addChannel <- &File{
		FileName: "hello.txt",
		UUID:     hash("hello.txt"),
		Uploaded: time.Now(),
	}

	fs.addChannel <- &File{
		FileName: "what.txt",
		UUID:     hash("what.txt"),
		Uploaded: time.Now(),
	}

	return fs
}

func (f *FileServer) Find(value string) ([]File, []int) {
	var files []File
	var indicies []int

	for k, file := range f.files {
		if file.FileName == value || string(file.UUID) == value {
			files = append(files, file)
			indicies = append(indicies, k)
		}
	}

	return files, indicies
}

func (f *FileServer) ListFiles(rw http.ResponseWriter, req *http.Request) {
	data, err := json.Marshal(f.files)
	if err != nil {
		rw.WriteHeader(500)
		rw.Write([]byte(err.Error()))
	}

	_, err = io.WriteString(rw, string(data))
	if err != nil {
		log.Println(err)
	}
}

func (f *FileServer) Start() {
	// Control the addition and removal of files to a channel
	go func() {
		for {
			select {
			case add := <-f.addChannel:
				f.files = append(f.files, *add)

			case delete := <-f.deleteChannel:
				_, toRemove := f.Find(string(delete.UUID))
				for _, v := range toRemove {
					f.files = append(f.files[:v], f.files[v+1:]...)
				}
			}
		}
	}()
}
