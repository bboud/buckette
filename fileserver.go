package main

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"os"
)

const MaxRecords = 100000

type ErrFileExists struct {
	ErrorString   string
	DuplicateFile string
}

func (err *ErrFileExists) Error() string {
	return err.ErrorString
}

type FileServer struct {
	Files        map[[32]byte]File
	RecordsCount int64
	QueueSize    int
	newFile      chan *File
}

func newFileServer() *FileServer {
	return &FileServer{
		QueueSize: 0,
		newFile:   make(chan *File, 5),
		Files:     make(map[[32]byte]File),
	}
}

func (fServer *FileServer) initialize() {

	LogPrint("Initializing file server! 🗄️")

	homedir, err := os.UserHomeDir()
	if err != nil {
		LogFatal(
			"Unable to load user's home directory",
			"Initialization of database",
			err)
	}

	LogPrint("Checking if data directories exist 🗃️")
	_, exists := os.Stat(homedir + FileStoreDir)
	if os.IsNotExist(exists) {
		err = os.MkdirAll(homedir+FileStoreDir, 0755)
	}

	_, exists = os.Stat(homedir + RecordStoreDir)
	if os.IsNotExist(exists) {
		err = os.MkdirAll(homedir+RecordStoreDir, 0755)
	}

	_, exists = os.Stat(homedir + TmpDir)
	if os.IsNotExist(exists) {
		err = os.MkdirAll(homedir+TmpDir, 0755)
	}

	if err != nil {
		log.Fatal("Cannot create database directories")
	}

	LogPrint("Loading all records into cache from disk 🏋️")

	dir, err := os.ReadDir(homedir + RecordStoreDir)
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

		recordData, err := os.ReadFile(homedir + RecordStoreDir + file.Name())
		if err != nil {
			LogFatal(
				"Unable to read record "+file.Name(),
				"Initialization of database",
				err)
		}

		var record *File
		err = json.Unmarshal(recordData, &record)
		if err != nil {
			LogFatal(
				"Unable to unmarshal record for "+file.Name(),
				"Initialization of database",
				err)
		}
		record.uuidName = encodeToString(record.UUID[:])
		fServer.push(record)
	}
	LogSucess("Fileserver is ready! 👻")
}

func (fServer *FileServer) start() {
	go fServer.handleNewFiles()

	select {}
}

func (fServer *FileServer) handleNewFiles() {
	for f := range fServer.newFile {
		fServer.Files[f.UUID] = *f
		fServer.RecordsCount += 1
		fServer.QueueSize -= 1
	}
}

func (fServer *FileServer) push(f *File) {
	fServer.QueueSize += 1
	fServer.newFile <- f
}

func (fServer *FileServer) exists(fileHash []byte) *ErrFileExists {
	homedir, err := os.UserHomeDir()
	if err != nil {
		LogFatal(
			"Unable to load user's home directory",
			"Initialization of database",
			err)
	}

	for _, file := range fServer.Files {
		if bytes.Equal(file.UUID[:], fileHash) {
			return &ErrFileExists{
				ErrorString:   file.FileName + " already exists in the cache",
				DuplicateFile: file.tmpHash,
			}
		}
	}

	// Only do disk lookups if we are at record max
	if fServer.RecordsCount <= MaxRecords {
		return nil
	}

	dir, err := os.ReadDir(homedir + FileStoreDir)
	if err != nil {
		log.Fatalf("Unable to open %s due to %e\n", homedir+FileStoreDir, err)
	}

	// Now for expensive
	for _, file := range dir {
		if file.IsDir() {
			continue
		}
		fName := encodeToString(fileHash)

		if fName == unSanitize(file.Name()) {
			// OK now we need to get the associated record
			record, err := os.ReadFile(RecordStoreDir + fName)
			if err != nil {
				log.Fatalf("Unable to read record %s error: %e\n", fName, err)
			}

			returnFile := &File{}
			err = json.Unmarshal(record, &returnFile)
			if err != nil {
				log.Println("Unable to marshal json from record lookup")
			}

			returnFile.uuidName = encodeToString(returnFile.UUID[:])
			// Add it to cache
			fServer.newFile <- returnFile

			return &ErrFileExists{
				ErrorString:   fName + " already exists in the cache",
				DuplicateFile: returnFile.URL,
			}
		}
	}

	return nil
}

// Finds the record
func (fServer *FileServer) FindByURL(url string) *File {
	for _, v := range fServer.Files {
		if v.URL == url {
			return &v
		}
	}
	return nil
}

func (fServer *FileServer) generateURL(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		LogFatal(
			"Unable to generate unique URL",
			"Generation of URL",
			err)
		return ""
	}

	hash := hex.EncodeToString(b)
	if fServer.FindByURL(hash) != nil {
		return fServer.generateURL(length)
	}
	return hash
}
