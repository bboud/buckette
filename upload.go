package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func (fServer *FileServer) HandleUpload(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	LogConnection(req)

	// Only do disk lookups if we are not at record max
	if fServer.RecordsCount >= MaxRecords {
		LogWarning(
			"This shouldn't have happened.",
			"Attempting to upload a file",
			errors.New("the server has reached max records"),
		)
		rw.WriteHeader(500)
		return
	}

	multipartReader, err := req.MultipartReader()
	if err != nil {
		LogWarning(
			"Multipart reader has failed",
			"Attempting to read multipart form",
			err,
		)
		return
	}
	if strings.HasPrefix(req.Header.Get("Content-Type"), "multipart/") {
		for {
			p, err := multipartReader.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				LogWarning(
					"Unable to read next part",
					"Attempting to read next part of the multipart form",
					err,
				)
				rw.WriteHeader(400)
				return
			}

			// Move this transmision black into the function
			fServer.HandleUploadPart(rw, p, req)
		}
	}
}

type UploadResponse struct {
	Duplicate bool
	File      *File
}

func (fServer *FileServer) HandleUploadPart(rw http.ResponseWriter, p *multipart.Part, req *http.Request) {
	fileName := req.Header.Get("File-Name")
	fileSize, err := strconv.ParseInt(req.Header.Get("File-Size"), 10, 64)
	if err != nil {
		LogWarning(
			"Unable to get file size",
			"Attempting to parse file size from header",
			err,
		)
		return
	}
	contentType := req.Header.Get("File-Type")

	if fileName == "" || fileSize <= 0 || contentType == "" {
		LogWarning(
			"FileName, FileSize, or ContentType is empty",
			"Attempting to parse file name/size/content-type from header",
			errors.New("no file desriptors in header"),
		)
		rw.WriteHeader(400)
		return
	}

	//t := time.Now()

	url, err := fServer.GenerateURL()
	if err != nil {
		LogWarning(
			"Unable to generate URL",
			"Handling of the part upload",
			err,
		)
		rw.WriteHeader(500)
		return
	}

	if err := copyToFileSystem(p, url); err != nil {
		LogWarning(
			"Unable to copy files to filesystem",
			"Handling of the part upload",
			err,
		)
		rw.WriteHeader(500)
		return
	}
	defer cleanTmp(url)

	// Now see if making a new file from the url.
	file, err := fServer.NewFile(url)
	if file == nil {
		LogWarning(
			"Unable to create new file",
			"Attempting to generate a new file for the response",
			err,
		)
		return
	}

	var fileExists *FileExists
	var uplResponse UploadResponse
	uplResponse.File = file
	uplResponse.Duplicate = errors.Is(err, fileExists)
	if !uplResponse.Duplicate {
		file.FileName = fileName
		file.URL = url
		file.Size = fileSize
		file.ContentType = contentType
		file.Uploaded = time.Now()
		writeRecord(file)
	}

	data, err := json.Marshal(uplResponse)
	if err != nil {
		LogWarning(
			"Unable to marshal the response data",
			"HandleUploadPart",
			err,
		)
		rw.WriteHeader(500)
		return
	}

	rw.Write(data)
}

const (
	FileStoreDir   = "./buckette-data/files/"
	RecordStoreDir = "./buckette-data/records/"
	TmpDir         = "./buckette-data/tmp/"
)

func copyToFileSystem(reader io.Reader, url string) error {
	file, err := os.OpenFile(TmpDir+"DAT_"+url, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		LogWarning(
			"Unable to create file "+TmpDir+"DAT_"+url,
			"Handling upload part for "+url,
			err,
		)
		return err
	}
	//defer file.Close()

	_, err = io.Copy(file, reader)
	if err != nil {
		LogWarning(
			"Unable to copy data to file "+TmpDir+"DAT_"+url,
			"Handling upload part for "+url,
			err,
		)
		return err
	}

	return nil
}

func writeRecord(f *File) error {
	uuidName := encodeToString(f.UUID[:])

	record, err := os.OpenFile(RecordStoreDir+uuidName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		LogWarning(
			"Unable to write record",
			"Handling upload of part",
			err,
		)
		return err
	}
	defer record.Close()

	//Record writer
	rData, err := json.Marshal(f)
	if err != nil {
		LogFatal(
			"Unable to marshal json",
			"Handling upload of part",
			err)
	}

	record.Write(rData)

	// We want to store the files using their hash for faster lookup on disk
	err = os.Rename(TmpDir+"DAT_"+f.URL, FileStoreDir+uuidName)
	if err != nil {
		LogFatal(
			"Unable to move temporary file to file store",
			"Moving file from temp to final storage for "+f.URL,
			err)
		return err
	}
	return nil
}

func cleanTmp(url string) error {
	if err := os.Remove(TmpDir + url); err != nil {
		if !os.IsNotExist(err) {
			LogWarning(
				"Unable to clean up after "+url,
				"Cleaning the temporary file directory",
				err)
			return err
		}
	}
	return nil
}

func encodeToString(v []byte) string {
	vName := base64.RawStdEncoding.EncodeToString(v)
	vName = strings.ReplaceAll(vName, "/", "*")
	return vName
}
