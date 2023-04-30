package main

import (
	"crypto/sha256"
	"encoding/base64"
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
	UUID          [32]byte
	URL           string
	Size          int64
	ContentType   string
	Uploaded      time.Time
	UserUploaded  string
	DownloadCount int

	tmpHash  string
	uuidName string
}

func (f *File) HandleUploadPart(part *multipart.Part, fServer *FileServer) ([]byte, error) {
	t := time.Now()
	f.tmpHash = fServer.generateURL(8)
	f.URL = f.tmpHash

	homedir, err := os.UserHomeDir()
	if err != nil {
		LogFatal(
			"Unable to load user's home directory",
			"Handling upload part for "+f.tmpHash,
			err)
	}

	record, err := os.OpenFile(homedir+TmpDir+f.tmpHash, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		LogWarning(
			"Unable to create file "+homedir+TmpDir+f.tmpHash,
			"Handling upload part for "+f.tmpHash,
			err,
		)
		return nil, err
	}
	defer record.Close()

	file, err := os.OpenFile(homedir+TmpDir+"DAT_"+f.tmpHash, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		LogWarning(
			"Unable to create file "+homedir+TmpDir+"DAT_"+f.tmpHash,
			"Handling upload part for "+f.tmpHash,
			err,
		)
		return nil, err
	}
	defer file.Close()

	_, err = io.Copy(file, part)
	if err != nil {
		return nil, err
	}

	if err := f.finalize(fServer); err != nil {
		LogWarning(
			"Duplicate value for "+f.uuidName,
			"Returning response for duplicate value",
			err)

		duplicate := &UploadResponse{
			URL:       err.DuplicateFile,
			Duplicate: true,
		}

		//Form a new request
		response, err := json.Marshal(duplicate)
		if err != nil {
			LogWarning(
				"Unable to marshal json for "+f.tmpHash,
				"Returning response for duplicate value",
				err)
			return nil, err
		}

		return response, nil
	}

	//Record writer
	rData, err := json.Marshal(f)
	if err != nil {
		cleanTmp(f.tmpHash)
		LogFatal(
			"Unable to marshal json for "+f.tmpHash,
			"Writing the record to the record store",
			err)
	}

	record.Write(rData)

	fServer.push(f)

	LogSucess("New file uploaded in " + time.Since(t).String() + "! 🍃")

	upload := &UploadResponse{
		URL:       f.URL,
		Duplicate: false,
		File:      *f,
	}

	response, err := json.Marshal(upload)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (f *File) Get() {

}

func (f *File) finalize(fServer *FileServer) *ErrFileExists {

	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	// Hash the file contents
	contents, err := os.ReadFile(homedir + TmpDir + "DAT_" + f.tmpHash)
	if err != nil {
		log.Fatal("Unable to read temp file")
		cleanTmp(f.tmpHash)
		log.Println("Unable to clean up after failed temp file read")
	}

	f.UUID = hasher(contents)
	if err := fServer.exists(f.UUID[:]); err != nil {
		return err
	}
	f.uuidName = encodeToString(f.UUID[:])

	err = os.Rename(homedir+TmpDir+f.tmpHash, homedir+RecordStoreDir+f.uuidName)
	if err != nil {
		cleanTmp(f.tmpHash)
		LogFatal(
			"Unable to move temporary file to file store",
			"Moving file from temp to final storage for "+f.tmpHash,
			err)
	}

	// We want to store the files using their hash for faster lookup on disk
	err = os.Rename(homedir+TmpDir+"DAT_"+f.tmpHash, homedir+FileStoreDir+f.uuidName)
	if err != nil {
		cleanTmp(f.tmpHash)
		LogFatal(
			"Unable to move temporary file to file store",
			"Moving file from temp to final storage for "+f.tmpHash,
			err)
	}

	return nil
}

func cleanTmp(tmpHash string) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	if err := os.Remove(homedir + TmpDir + tmpHash); err != nil {
		log.Println(err)
		log.Fatal("Remove file brokie")
	}
	if err := os.Remove(homedir + TmpDir + "DAT_" + tmpHash); err != nil {
		log.Println(err)
		log.Fatal("Remove file brokie")
	}
}

func hasher(value []byte) [32]byte {
	return sha256.Sum256(value)
}

func encodeToString(v []byte) string {
	vName := base64.RawStdEncoding.EncodeToString(v)
	vName = strings.ReplaceAll(vName, "/", "*")
	return vName
}

func unSanitize(v string) string {
	v = strings.ReplaceAll(v, "*", "/")
	return v
}

func decodeFromString(v string) [32]byte {
	v = strings.ReplaceAll(v, "*", "/")
	s, err := base64.RawStdEncoding.DecodeString(v)
	if err != nil {
		LogFatal(
			"Unable to decode string into bytes",
			"String decoding",
			err)
	}
	return [32]byte(s)
}
