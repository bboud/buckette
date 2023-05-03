package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"os"
)

const MaxRecords = 100000

type FileServer struct {
	Files        map[[32]byte]File
	URLs         map[string][32]byte
	RecordsCount int64
	newFile      chan *File
}

func NewFileServer() *FileServer {
	return &FileServer{
		newFile: make(chan *File),
		Files:   make(map[[32]byte]File),
		URLs:    make(map[string][32]byte),
	}
}

func makeDataStore() {
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
		LogFatal(
			"Cannot create database directories",
			"Initialization of database",
			err)
	}
}

func (fServer *FileServer) LoadFromDisk() {
	dir, err := os.ReadDir(RecordStoreDir)
	if err != nil {
		LogFatal(
			"Unable to load records into cache from disk",
			"Initialization of database",
			err)
	}

	for _, file := range dir {
		if file.IsDir() {
			continue
		}

		if fServer.RecordsCount >= MaxRecords {
			LogWarning(
				"Max records reached",
				"Initialization of database",
				errors.New("no need for action"),
			)
			return
		}

		recordData, err := os.ReadFile(RecordStoreDir + file.Name())
		if err != nil {
			LogFatal(
				"Unable to read record "+file.Name(),
				"Initialization of database",
				err)
		}

		var record File
		err = json.Unmarshal(recordData, &record)
		if err != nil {
			LogFatal(
				"Unable to unmarshal record for "+file.Name(),
				"Initialization of database",
				err)
		}
		record.uuidName = encodeToString(record.UUID[:])
		fServer.Push(record)
	}
}

func (fServer *FileServer) Initialize() {

	LogPrint("Initializing file server! 🗄️")

	LogPrint("Checking if data directories exist 🗃️")
	makeDataStore()

	LogPrint("Loading all records into cache from disk 🏋️")
	fServer.LoadFromDisk()

	LogSuccess("Fileserver is ready! 👻")
}

func (fServer *FileServer) Push(f File) {
	fServer.Files[f.UUID] = f
	fServer.URLs[f.URL] = f.UUID
	fServer.RecordsCount += 1
}

// This will need a condition to check the disk once past max records in ram
// but this can be done later..
func (fServer *FileServer) Exists(fileHash []byte) string {
	for _, file := range fServer.Files {
		if bytes.Equal(file.UUID[:], fileHash) {
			return file.URL
		}
	}

	return ""
}

// Finds the record
func (fServer *FileServer) FindByURL(url string) *File {
	uuid, ok := fServer.URLs[url]
	if !ok {
		return nil
	}
	file := fServer.Files[uuid]
	return &file
}

// Rename to capital
func (fServer *FileServer) GenerateURL(f *File) error {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		LogFatal(
			"Unable to generate unique URL",
			"Generation of URL",
			err)
		return err
	}

	hash := hex.EncodeToString(b)
	if fServer.FindByURL(hash) != nil {
		return fServer.GenerateURL(f)
	}

	f.tmpHash = hash
	f.URL = hash

	return nil
}

type FileExists struct {
	File File
}

func (e *FileExists) Error() string {
	return "file already exists"
}

func (fServer *FileServer) CheckHash(f *File) error {
	// Hash the file contents
	contents, err := os.ReadFile(TmpDir + "DAT_" + f.tmpHash)
	if err != nil {
		log.Fatal("Unable to read temp file")
		log.Println("Unable to clean up after failed temp file read")
	}

	f.UUID = sha256.Sum256(contents)
	f.uuidName = encodeToString(f.UUID[:])
	if fServer.Exists(f.UUID[:]) == "" {
		return &FileExists{
			File: fServer.Files[f.UUID],
		}
	}

	return nil
}
