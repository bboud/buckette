package main

import (
	"bufio"
	"compress/lzw"
	"crypto/sha256"
	"encoding/base32"
	"io"
	"log"
	"mime/multipart"
	"os"
	"time"
)

const (
	FileStoreDir   = "/home/brand0n/data/filer/store/files/"
	RecordStoreDir = "/home/brand0n/data/filer/store/records/"
	bufferSize     = 4096
)

type FileServerError struct {
	ErrorString string
}

func (err *FileServerError) Error() string {
	return err.ErrorString
}

type FileContent []byte

type File struct {
	FileName      string
	UUID          string
	Size          int64
	ContentType   string
	Uploaded      time.Time
	UserUploaded  string
	DownloadCount int
}

func hash(data []byte) string {
	h := sha256.Sum256(data)
	return base32.HexEncoding.EncodeToString(h[:])
}

func (f *File) Push(part *multipart.Part) error {
	bufReader := bufio.NewReaderSize(part, bufferSize)

	hashable, err := bufReader.Peek(256)
	if err != nil {
		return err
	}

	f.UUID = hash(hashable)

	record, err := os.OpenFile(RecordStoreDir+string(f.UUID), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer record.Close()
	file, err := os.OpenFile(FileStoreDir+string(f.UUID), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Good until here!

	// Record writer
	// rData, err := json.Marshal(f)
	// if err != nil {
	// 	return err
	// }

	compressor := lzw.NewWriter(file, lzw.LSB, 8)

	buffer := make([]byte, bufferSize)
	rw := bufio.NewReadWriter(bufio.NewReader(part), bufio.NewWriter(compressor))
	for {
		n, err := rw.Read(buffer)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		buffer = buffer[:n]
		//fmt.Println(buffer)
		rw.Write(buffer)
		rw.Flush()
	}
}

func (f *File) Get() {

}

// bad bad rewrite this
func initFileServer() {
	_, err := os.Stat(FileStoreDir)
	_, err2 := os.Stat(RecordStoreDir)
	if os.IsNotExist(err) && os.IsNotExist(err2) {
		err = os.MkdirAll(FileStoreDir, 0755)
		err2 = os.MkdirAll(RecordStoreDir, 0755)
		if err != nil || err2 != nil {
			log.Println(err)
			log.Println(err2)
		}
		return
	}
	if err != nil {
		log.Println("It is possible only one of the directories exist for the file server")
		log.Fatal(err)
	}
}
