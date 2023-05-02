package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type UploadResponse struct {
	URL       string
	Duplicate bool
	File      File
}

func (fServer *FileServer) HandleUpload(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	LogConnection(req)

	// Only do disk lookups if we are at record max
	if fServer.RecordsCount <= MaxRecords {
		LogWarning(
			"This shouldn't have happened.",
			"Attempting to ",
			errors.New("the server has reached max records"),
		)
		rw.WriteHeader(500)
		return
	}

	multipartReader, err := req.MultipartReader()
	if err != nil {
		log.Printf("ERROR: %v | Called by: %s", err, "MultipartReader")
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
				rw.WriteHeader(500)
				return
			}

			fileName := req.Header.Get("File-Name")
			fileSize, err := strconv.ParseInt(req.Header.Get("File-Size"), 10, 64)
			if err != nil {
				LogWarning(
					"Unable to get file size",
					"Attempting to parse file size from header",
					err,
				)
				rw.WriteHeader(500)
				return
			}
			contentType := req.Header.Get("File-Type")
			hash := req.Header.Get("File-Hash")
			url := fServer.Exists([]byte(hash))

			if url != "" {
				uploadResponse := UploadResponse{
					Duplicate: true,
					URL:       url,
				}
				data, err := json.Marshal(uploadResponse)
				if err != nil {
					LogWarning(
						"Unable to marshal json response for "+url,
						"Upload handler",
						err,
					)
					rw.WriteHeader(500)
					return
				}
				rw.WriteHeader(200)
				rw.Write(data)
			} else {
				uploadResponse := UploadResponse{
					Duplicate: true,
					URL:       url,
				}
				// Handle the file saving in another goroutine
				uploadResponse.File = File{
					FileName:      fileName,
					Uploaded:      time.Now(),
					Size:          fileSize,
					DownloadCount: 0,
					ContentType:   contentType,
				}

				// This maybe shouldn't be monolithic
				err = fServer.HandleUploadPart(&uploadResponse.File, p)
				if err != nil {
					LogWarning(
						"Unable to marshal json response for "+url,
						"Upload handler",
						err,
					)
					rw.WriteHeader(500)
					return
				}

				data, err := json.Marshal(uploadResponse)
				if err != nil {
					LogWarning(
						"Unable to marshal json response for "+url,
						"Upload handler",
						err,
					)
					rw.WriteHeader(500)
					return
				}

				rw.WriteHeader(200)
				rw.Write(data)
			}
		}
	}
}

func (fServer *FileServer) HandleUploadPart(f *File, part *multipart.Part) error {
	t := time.Now()

	err := fServer.GenerateURL(f)
	if err != nil {
		LogWarning(
			"Unable to generate URL",
			"Handling of the part upload for "+f.tmpHash,
			err,
		)
		return err
	}

	if err := f.copyToFileSystem(part); err != nil {
		f.cleanTmp(false)
		LogWarning(
			"Unable to copy files to filesystem",
			"Handling of the part upload for "+f.tmpHash,
			err,
		)
		return err
	}

	// Check Hash
	// Most expensive call
	if err := fServer.CheckHash(f); err != nil {
		f.cleanTmp(false)
		LogWarning(
			"Unable to check hash for file",
			"Handling of the part upload for "+f.tmpHash,
			err,
		)
		return err
	}

	if err := f.cleanTmp(true); err != nil {
		LogFatal(
			"Unable to move files into database store!",
			"Handling of the part upload for "+f.tmpHash,
			err,
		)
		return err
	}

	// Save a copy
	fServer.Push(*f)

	LogSuccess("New file uploaded in " + time.Since(t).String() + "! 🍃")

	return nil
}
