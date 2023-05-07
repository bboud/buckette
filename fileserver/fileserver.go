package fileserver

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"buckette/logger"
)

const MaxRecords = 100000

const (
	FileStoreDir   = "./buckette-data/files/"
	RecordStoreDir = "./buckette-data/records/"
	TmpDir         = "./buckette-data/tmp/"
)

type File struct {
	FileName      string
	UUID          string
	URL           string
	Size          int64
	ContentType   string
	Uploaded      time.Time
	UserUploaded  string
	DownloadCount int
}

type FileServer struct {
	Files        map[string]*File
	URLs         map[string]string
	RecordsCount int64
}

func NewFileServer() *FileServer {

	fServer := &FileServer{
		Files: make(map[string]*File),
		URLs:  make(map[string]string),
	}

	err := fServer.initialize()
	if err != nil {
		logger.Fatal(
			"Unable to initialize file server",
			"fileserver.NewFileServer",
			err,
		)
		os.Exit(1)
	}

	return fServer
}

func makeDataStore() error {
	_, err := os.Stat(FileStoreDir)
	notExists := os.IsNotExist(err)
	if notExists {
		err = os.MkdirAll(FileStoreDir, 0755)
		if err != nil {
			return err
		}
	} else {
		return err
	}

	_, err = os.Stat(RecordStoreDir)
	notExists = os.IsNotExist(err)
	if notExists {
		err = os.MkdirAll(RecordStoreDir, 0755)
		if err != nil {
			return err
		}
	} else {
		return err
	}

	_, err = os.Stat(TmpDir)
	notExists = os.IsNotExist(err)
	if notExists {
		err = os.MkdirAll(TmpDir, 0755)
		if err != nil {
			return err
		}
	} else {
		return err
	}

	return nil
}

func (fServer *FileServer) loadFromDisk() error {
	dir, err := os.ReadDir(RecordStoreDir)
	if err != nil {
		return err
	}

	for _, file := range dir {
		if file.IsDir() {
			continue
		}

		if fServer.RecordsCount >= MaxRecords {
			return &MaxRecordsReached{}
		}

		recordData, err := os.ReadFile(RecordStoreDir + file.Name())
		if err != nil {
			return err
		}

		var record *File
		err = json.Unmarshal(recordData, &record)
		if err != nil {
			return err
		}
		fServer.push(record)
	}
	return nil
}

func (fServer *FileServer) initialize() error {

	logger.Print("Initializing file server! 🗄️")

	logger.Print("Checking if data directories exist 🗃️")
	err := makeDataStore()
	if err != nil {
		return err
	}

	logger.Print("Loading all records into cache from disk 🏋️")
	err = fServer.loadFromDisk()
	if err != nil {
		return err
	}

	logger.Success("Fileserver is ready! 👻")
	return nil
}

func (fServer *FileServer) push(f *File) {
	fServer.Files[f.UUID] = f
	fServer.URLs[f.URL] = f.UUID
	fServer.RecordsCount += 1
}

// This will need a condition to check the disk once past max records in ram
// but this can be done later..
func (fServer *FileServer) exists(fileHash string) bool {
	for _, file := range fServer.Files {
		if file.UUID == fileHash {
			return true
		}
	}
	fmt.Println("File doesn't exist")
	return false
}

// Finds the record
func (fServer *FileServer) findByURL(url string) *File {
	uuid, ok := fServer.URLs[url]
	if !ok {
		return nil
	}
	file := fServer.Files[uuid]
	return file
}

func (fServer *FileServer) findByUUID(uuid string) *File {
	file, ok := fServer.Files[uuid]
	if !ok {
		return nil
	} else {
		return file
	}
}

func (fServer *FileServer) generateURL() (string, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	hash := hex.EncodeToString(b)
	if fServer.findByURL(hash) != nil {
		return fServer.generateURL()
	}

	return hash, nil
}

func (fServer *FileServer) NewFile(tmpHash string) (*File, error) {
	// Hash the file contents
	contents, err := os.ReadFile(TmpDir + "DAT_" + tmpHash)
	if err != nil {
		return nil, err
	}

	hash := sha256.Sum256(contents)
	uuid := encodeToString(hash[:])
	if fServer.exists(uuid) {
		return fServer.findByUUID(uuid), &FileExists{}
	}

	// Set the UUID for the new file!
	return &File{
		UUID: uuid,
	}, nil
}
