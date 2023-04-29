package main

import (
	"log"
	"os"
)

type FileServerRequest struct {
	file     string
	response chan *File
	err      error
}

type FileServer struct {
	Files     []File
	QueueSize int
	newFile   chan *File
	request   chan *FileServerRequest
}

func newFileServer() *FileServer {
	return &FileServer{
		QueueSize: 0,
		newFile:   make(chan *File, 5),
		request:   make(chan *FileServerRequest),
	}
}

func (fs *FileServer) start() {
	var err error
	_, exists := os.Stat(FileStoreDir)
	if os.IsNotExist(exists) {
		err = os.MkdirAll(FileStoreDir, 0755)
	}

	_, exists = os.Stat(RecordStoreDir)
	if os.IsNotExist(exists) {
		err = os.MkdirAll(RecordStoreDir, 0755)
	}

	_, exists = os.Stat(TmpDir)
	if os.IsNotExist(exists) {
		err = os.MkdirAll(TmpDir, 0755)
	}

	if err != nil {
		log.Fatal("Cannot create database directories")
	}

	go fs.handleNewFiles()
	//go fs.handleRequests()

	// Block function exit
	select {}
}

func (fs *FileServer) push(f *File) {
	fs.QueueSize += 1
	fs.newFile <- f
}

func (fs *FileServer) handleNewFiles() {
	for f := range fs.newFile {
		fs.Files = append(fs.Files, *f)
		fs.QueueSize -= 1
	}
}
