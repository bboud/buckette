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
	FileStoreDir   = "./buckette-data/files/"
	RecordStoreDir = "./buckette-data/records/"
	TmpDir         = "./buckette-data/tmp/"
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

	record, err := os.OpenFile(TmpDir+f.tmpHash, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		LogWarning(
			"Unable to create file "+TmpDir+f.tmpHash,
			"Handling upload part for "+f.tmpHash,
			err,
		)
		return nil, err
	}
	defer record.Close()

	file, err := os.OpenFile(TmpDir+"DAT_"+f.tmpHash, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		LogWarning(
			"Unable to create file "+TmpDir+"DAT_"+f.tmpHash,
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

	defer cleanTmp(f.tmpHash)

	if err := f.finalize(fServer); err != nil {
		LogWarning(
			"Duplicate value for "+f.uuidName,
			"Returning response for duplicate value",
			err)

		duplicate := &UploadResponse{
			URL:       fServer.Files[f.UUID].URL,
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

func (f *File) finalize(fServer *FileServer) *ErrFileExists {
	// Hash the file contents
	contents, err := os.ReadFile(TmpDir + "DAT_" + f.tmpHash)
	if err != nil {
		log.Fatal("Unable to read temp file")
		log.Println("Unable to clean up after failed temp file read")
	}

	f.UUID = hasher(contents)
	if err := fServer.exists(f.UUID[:]); err != nil {
		return err
	}
	f.uuidName = encodeToString(f.UUID[:])

	err = os.Rename(TmpDir+f.tmpHash, RecordStoreDir+f.uuidName)
	if err != nil {
		LogFatal(
			"Unable to move temporary file to file store",
			"Moving file from temp to final storage for "+f.tmpHash,
			err)
	}

	// We want to store the files using their hash for faster lookup on disk
	err = os.Rename(TmpDir+"DAT_"+f.tmpHash, FileStoreDir+f.uuidName)
	if err != nil {
		LogFatal(
			"Unable to move temporary file to file store",
			"Moving file from temp to final storage for "+f.tmpHash,
			err)
	}

	return nil
}

func cleanTmp(tmpHash string) {
	if err := os.Remove(TmpDir + tmpHash); err != nil {
		if !os.IsNotExist(err) {
			LogWarning(
				"Unable to clean up after "+tmpHash,
				"Cleaning the temporary file directory",
				err)
		}
	}
	if err := os.Remove(TmpDir + "DAT_" + tmpHash); err != nil {
		if !os.IsNotExist(err) {
			LogWarning(
				"Unable to clean up after "+tmpHash,
				"Cleaning the temporary file directory",
				err)
		}
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
