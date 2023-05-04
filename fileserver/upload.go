package fileserver

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

	"github.com/bboud/buckette/logger"
)

func (fServer *FileServer) HandleUpload(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	logger.LogConnection(req)

	// Only do disk lookups if we are not at record max
	if fServer.RecordsCount >= MaxRecords {
		logger.LogWarning(
			"This shouldn't have happened.",
			"filerserver.HandleUpload",
			errors.New("the server has reached max records"),
		)
		rw.WriteHeader(500)
		return
	}

	multipartReader, err := req.MultipartReader()
	if err != nil {
		logger.LogWarning(
			"Multipart reader has failed",
			"filerserver.HandleUpload",
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
				logger.LogWarning(
					"Unable to read next part",
					"filerserver.HandleUpload",
					err,
				)
				rw.WriteHeader(400)
				return
			}

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
		logger.LogWarning(
			"Unable to get file size",
			"filerserver.HandleUploadPart",
			err,
		)
		return
	}
	contentType := req.Header.Get("File-Type")

	if fileName == "" || fileSize <= 0 || contentType == "" {
		logger.LogWarning(
			"FileName, FileSize, or ContentType is empty",
			"filerserver.HandleUploadPart",
			errors.New("no file desriptors in header"),
		)
		rw.WriteHeader(400)
		return
	}

	//t := time.Now()

	url, err := fServer.generateURL()
	if err != nil {
		logger.LogWarning(
			"Unable to generate URL",
			"filerserver.HandleUploadPart",
			err,
		)
		rw.WriteHeader(500)
		return
	}

	if err := copyToFileSystem(p, url); err != nil {
		logger.LogWarning(
			"Unable to copy files to filesystem",
			"filerserver.HandleUploadPart",
			err,
		)
		rw.WriteHeader(500)
		return
	}
	defer cleanTmp(url)

	// Now see if making a new file from the url.
	file, err := fServer.NewFile(url)
	if file == nil {
		logger.LogWarning(
			"Unable to create new file",
			"filerserver.HandleUploadPart",
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
		logger.LogWarning(
			"Unable to marshal the response data",
			"filerserver.HandleUploadPart",
			err,
		)
		rw.WriteHeader(500)
		return
	}

	rw.Write(data)
}

func copyToFileSystem(reader io.Reader, url string) error {
	file, err := os.OpenFile(TmpDir+"DAT_"+url, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	//defer file.Close()

	_, err = io.Copy(file, reader)
	if err != nil {
		return err
	}

	return nil
}

func writeRecord(f *File) error {
	uuid := f.UUID

	record, err := os.OpenFile(RecordStoreDir+uuid, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer record.Close()

	//Record writer
	rData, err := json.Marshal(f)
	if err != nil {
		return err
	}

	record.Write(rData)

	// We want to store the files using their hash for faster lookup on disk
	err = os.Rename(TmpDir+"DAT_"+f.URL, FileStoreDir+uuid)
	if err != nil {
		return err
	}
	return nil
}

func cleanTmp(url string) error {
	if err := os.Remove(TmpDir + url); err != nil {
		if !os.IsNotExist(err) {
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
