package main

import (
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type UploadResponse struct {
	URL        string
	Duplicate  bool
	File       File
	StatusCode int
}

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
			data, err := fServer.HandleUploadPart(rw, p, req)
			if err != nil {
				LogWarning(
					"Unable to complete request",
					"Attempting to get data to send to the response writer",
					err,
				)
				rw.WriteHeader(400)
				return
			}

			rw.WriteHeader(200)
			rw.Write(data)
		}
	}
}

func (fServer *FileServer) HandleUploadPart(rw http.ResponseWriter, p *multipart.Part, req *http.Request) ([]byte, error) {

	var uploadResponse = UploadResponse{
		Duplicate: false,
	}

	fileName := req.Header.Get("File-Name")
	fileSize, err := strconv.ParseInt(req.Header.Get("File-Size"), 10, 64)
	if err != nil {
		LogWarning(
			"Unable to get file size",
			"Attempting to parse file size from header",
			err,
		)

		return nil, err
	}
	contentType := req.Header.Get("File-Type")

	if fileName == "" || fileSize <= 0 || contentType == "" {
		LogWarning(
			"FileName, FileSize, or ContentType is empty",
			"Attempting to parse file name/size/content-type from header",
			errors.New("no file desriptors in header"),
		)
		rw.WriteHeader(400)
		return nil, err
	}

	t := time.Now()

	uploadResponse.File = File{
		FileName:      fileName,
		Uploaded:      time.Now(),
		Size:          fileSize,
		DownloadCount: 0,
		ContentType:   contentType,
	}

	err = fServer.GenerateURL(&uploadResponse.File)
	if err != nil {
		LogWarning(
			"Unable to generate URL",
			"Handling of the part upload",
			err,
		)
		rw.WriteHeader(500)
		return nil, err
	}

	if err := uploadResponse.File.copyToFileSystem(p); err != nil {
		uploadResponse.File.cleanTmp(false)
		LogWarning(
			"Unable to copy files to filesystem",
			"Handling of the part upload",
			err,
		)
		rw.WriteHeader(500)
		return nil, err
	}

	// Check Hash
	// Most expensive call
	err = fServer.CheckHash(&uploadResponse.File)
	var success = true

	//Drop this through for clean up stage
	tempHash := uploadResponse.File.tmpHash

	var fileExists *FileExists
	if errors.As(err, &fileExists) {
		uploadResponse = UploadResponse{
			URL:       fileExists.File.URL,
			Duplicate: true,
			File:      fileExists.File,
		}
		uploadResponse.File.tmpHash = tempHash
		success = false
	} else if err != nil {
		uploadResponse.File.cleanTmp(false)
		LogWarning(
			"Unable to check hash for file",
			"Handling of the part upload for "+uploadResponse.File.tmpHash,
			err,
		)
		rw.WriteHeader(500)
		return nil, err
	}

	if err := uploadResponse.File.cleanTmp(success); err != nil {
		LogFatal(
			"Unable to move files into database store!",
			"Handling of the part upload for "+uploadResponse.File.tmpHash,
			err,
		)
		rw.WriteHeader(500)
		return nil, err
	}

	if success {
		// Save a copy
		fServer.Push(uploadResponse.File)

		LogSuccess("New file uploaded in " + time.Since(t).String() + "! 🍃")
	}

	data, err := json.Marshal(uploadResponse)
	if err != nil {
		LogWarning(
			"Unable to marshal upload response",
			"Returning data from the upload handler",
			err,
		)
		rw.WriteHeader(500)
		return nil, err
	}

	return data, nil
}
