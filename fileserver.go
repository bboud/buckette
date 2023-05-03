package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
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
func (fServer *FileServer) Exists(fileHash []byte) bool {
	for _, file := range fServer.Files {
		if bytes.Equal(file.UUID[:], fileHash) {
			return true
		}
	}

	return false
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

func (fServer *FileServer) GenerateURL() (string, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		LogFatal(
			"Unable to generate unique URL",
			"Generation of URL",
			err)
		return "", err
	}

	hash := hex.EncodeToString(b)
	if fServer.FindByURL(hash) != nil {
		return fServer.GenerateURL()
	}

	return hash, nil
}

type FileExists struct{}

func (e *FileExists) Error() string {
	return "file already exists"
}

func (fServer *FileServer) NewFile(tmpHash string) (*File, error) {
	// Hash the file contents
	contents, err := os.ReadFile(TmpDir + "DAT_" + tmpHash)
	if err != nil {
		return nil, err
	}

	uuid := sha256.Sum256(contents)
	if fServer.Exists(uuid[:]) {
		return fServer.FindByURL(tmpHash), &FileExists{}
	}

	// Set the UUID for the new file!
	return &File{
		UUID: uuid,
	}, nil
}
