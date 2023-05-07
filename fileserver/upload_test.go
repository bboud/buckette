package fileserver

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

func TestCopyToFileSystem(t *testing.T) {
	fServer := NewFileServer()
	fServer.initialize()

	url, err := fServer.generateURL()
	if err != nil {
		t.Error(err)
		return
	}

	err = copyToFileSystem(strings.NewReader("This is a string!"), url)
	if err != nil {
		t.Error(err)
		return
	}

	// Read the file
	fileRead, err := os.ReadFile(TmpDir + "DAT_" + url)
	if err != nil {
		t.Error(err)
		return
	}
	if string(fileRead) == "This is a string!" {
		t.Log("Successfully copied data to file system")
	}
}

func TestWriteRecord(t *testing.T) {
	fServer := NewFileServer()
	fServer.initialize()

	url, err := fServer.generateURL()
	if err != nil {
		t.Error(err)
	}

	uuid := sha256.Sum256([]byte("Some data"))

	file := File{
		UUID: encodeToString(uuid[:]),
		URL:  url,
	}

	uuidName := file.UUID

	err = copyToFileSystem(strings.NewReader("This is a string!"), url)
	if err != nil {
		t.Error(err)
	}

	err = fServer.writeRecord(&file)
	if err != nil {
		t.Error(err)
	}

	_, err = os.Stat(FileStoreDir + uuidName)
	if !errors.Is(err, os.ErrNotExist) {
		t.Log("File was successfully moved")
	} else {
		t.Error(err)
		return
	}

	_, err = os.Stat(RecordStoreDir + uuidName)
	if !errors.Is(err, os.ErrNotExist) {
		t.Log("File was successfully moved")
	} else {
		t.Error(err)
		return
	}

	_, err = os.Stat(TmpDir + url)
	if errors.Is(err, os.ErrNotExist) {
		t.Log("File was successfully removed")
	} else {
		t.Error(err)
		return
	}
}

func TestCleanTmp(t *testing.T) {
	fServer := NewFileServer()
	fServer.initialize()

	url, _ := fServer.generateURL()

	err := os.WriteFile(TmpDir+"DAT_"+url, []byte("Some stuff"), 0644)
	if err != nil {
		t.Error(err)
	}

	cleanTmp(url)

	_, err = os.Stat(TmpDir + url)
	if errors.Is(err, os.ErrNotExist) {
		t.Log("File was successfully removed")
		return
	} else {
		t.Error("File was unable to be removed")
	}
}

// This request is reconstructed from an actual valid request
func postRequest(fServer *FileServer) error {
	file, err := os.Open("tmp_request.txt")
	if err != nil {
		return err
	}

	// Send as if a browser were sending the request
	client := http.Client{}
	request, err := http.NewRequest("POST", "http://localhost:8080/upl", file)
	if err != nil {
		return err
	}

	url, err := fServer.generateURL()
	if err != nil {
		return err
	}

	request.Header.Add("File-Size", "14")
	request.Header.Add("File-Name", "tmp.txt")
	request.Header.Add("File-Type", "text/plain")
	request.Header.Add("Content-Type", "multipart/form-data; boundary=---------------------------190042104612722886532900493499")

	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	uuid := sha256.Sum256([]byte("this is a test"))

	// Expect:
	expectedResponse := UploadResponse{
		Duplicate: false,
		File: &File{
			UUID:        encodeToString(uuid[:]),
			FileName:    "tmp.txt",
			Size:        14,
			ContentType: "text/plain",
			URL:         url,
		},
	}
	////////////

	var uplResponse *UploadResponse
	err = json.Unmarshal(content, &uplResponse)
	if err != nil {
		return err
	}

	// // We need to check if the response content contains the correct information
	if strings.Compare(uplResponse.File.FileName, expectedResponse.File.FileName) != 0 {
		return errors.New("file names do not match")
	}

	//!bytes.Equal(expectedResponse.File.UUID.UUIDByte[:], uplResponse.File.UUID.UUIDByte[:])
	if strings.Compare(expectedResponse.File.UUID, uplResponse.File.UUID) != 0 {
		return errors.New("uuids do not match")
	}

	if uplResponse.File.Size != expectedResponse.File.Size {
		return errors.New("file sizes do not match")
	}

	if strings.Compare(uplResponse.File.ContentType, expectedResponse.File.ContentType) != 0 {
		return errors.New("file content types do not match")
	}

	// Best we can do is verify the block comes back with most of the right information

	return nil
}

func TestHandleUpload(t *testing.T) {

	var err error

	// If it exists, clear it
	err = os.RemoveAll("./buckette-data")
	if !errors.Is(err, os.ErrNotExist) && err != nil {
		t.Log(err)
	}

	fServer := NewFileServer()
	fServer.initialize()

	m := http.NewServeMux()
	server := http.Server{
		Addr:    ":8080",
		Handler: m,
	}

	go func() {
		// Sleep for a couple seconds to let the listener start
		time.Sleep(2 * time.Second)

		var fileExists *FileExists
		err := postRequest(fServer)
		if errors.Is(err, fileExists) {
			t.Log("File already exists, good!")
		} else if !errors.Is(err, fileExists) && err != nil {
			t.Error(err)
		}

		if err := server.Shutdown(context.Background()); err != nil {
			t.Error(err)
		}

		t.Log("Response is formed correctly")
	}()

	m.HandleFunc("/upl", fServer.HandleUpload)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		t.Error(err)
	}

	time.Sleep(2 * time.Second) // Block for enough time for test go routine can finish
}
