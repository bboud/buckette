package main

import (
	"crypto/sha256"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strings"
	"time"
)

const (
	FileStoreDir   = "/data/filer/store/files/"
	RecordStoreDir = "/data/filer/store/records/"
	TmpDir         = "/data/filer/tmp/"
	bufferSize     = 8 * 1024
)

type FileContent []byte

type File struct {
	FileName      string
	RecordHash    string
	FileHash      [32]byte
	Size          int64
	ContentType   string
	Uploaded      time.Time
	UserUploaded  string
	DownloadCount int
	rName         string
}

func (f *File) Push(part *multipart.Part, recordHash string) ([]byte, error) {
	f.RecordHash = recordHash
	f.rName = strings.ReplaceAll(f.RecordHash, "/", "")

	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	record, err := os.OpenFile(homedir+TmpDir+f.rName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	defer record.Close()

	file, err := os.OpenFile(homedir+TmpDir+"DAT_"+f.rName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	//Record writer
	rData, err := json.Marshal(f)
	if err != nil {
		cleanTmp(f.rName)
		return nil, err
	}

	record.Write(rData)

	// buffer := make([]byte, bufferSize)
	// for {
	// 	n, err := part.Read(buffer)
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	if err != nil {
	// 		cleanTmp(f.rName)
	// 		return nil, err
	// 	}

	// 	buffer = buffer[:n]
	// 	_, err = file.Write(buffer)
	// 	if err != nil {
	// 		cleanTmp(f.rName)
	// 		return nil, err
	// 	}
	// }
	content, err := io.ReadAll(part)
	if err != nil {
		log.Fatal(err)
	}
	file.Write(content)

	f.finalize(file)

	response, err := json.Marshal(f)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (f *File) Get() {

}

func (f *File) finalize(toHash *os.File) {

	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	// Hash the file contents
	contents, err := io.ReadAll(toHash)
	if err != nil {
		log.Fatal("Unable to read temp file")
		cleanTmp(f.rName)
		return
	}

	f.FileHash = hasher(contents)

	err = os.Rename(homedir+TmpDir+f.rName, homedir+RecordStoreDir+f.rName)
	if err != nil {
		log.Fatal(err)
	}
	err = os.Rename(homedir+TmpDir+"DAT_"+f.rName, homedir+FileStoreDir+f.rName)
	if err != nil {
		log.Fatal(err)
	}
}

func cleanTmp(recordHash string) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	if err := os.Remove(homedir + TmpDir + recordHash); err != nil {
		log.Println(err)
		log.Fatal("Remove file brokie")
	}
	if err := os.Remove(homedir + TmpDir + "DAT_" + recordHash); err != nil {
		log.Println(err)
		log.Fatal("Remove file brokie")
	}
}

func hasher(value []byte) [32]byte {
	return sha256.Sum256(value)
}
