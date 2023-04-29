package main

import (
	"bufio"
	"compress/lzw"
	"crypto/sha256"
	"encoding/json"
	"os"
	"time"
)

const (
	FileStoreDir   = "/data/filer/store/files/"
	RecordStoreDir = "/data/filer/store/records/"
	bufferSize     = 4096
)

type FileServerError struct {
	ErrorString string
}

type FileContent []byte

type File struct {
	FileName      string
	UUID          []byte
	Size          int64
	ContentType   string
	Uploaded      time.Time
	UserUploaded  string
	DownloadCount int
}

func hash(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}

func (f *File) Push(content *FileContent) error {
	f.UUID = hash(*content)
	record, err := os.OpenFile(RecordStoreDir+string(f.UUID), os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		return err
	}
	defer record.Close()
	file, err := os.OpenFile(FileStoreDir+string(f.UUID), os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		return err
	}
	defer file.Close()

	rData, err := json.Marshal(f)
	if err != nil {
		return err
	}

	record.Write(rData)

	bufWritter := bufio.NewWriterSize(file, bufferSize)
	// Do the compression!
	compressWriter := lzw.NewWriter(bufWritter, lzw.LSB, 8)
	compressWriter.Write(*content)

	return nil
}

func (f *File) Get() {

}
