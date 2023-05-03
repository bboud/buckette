package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

// This request is reconstructed from an actual valid request
func postRequest() error {
	file, err := os.Open("tmp.txt")
	if err != nil {
		return err
	}

	// Send as if a browser were sending the request
	client := http.Client{}
	// resp, err := client.Post("http://localhost:8080/upl", "multipart/form-data; boundary=---------------------------190042104612722886532900493499", file)
	// if err != nil {
	// 	return err
	// }
	// defer resp.Body.Close()

	// File-Hash: MmU5OTc1ODU0ODk3MmE4ZTg4MjJhZDQ3ZmExMDE3ZmY3MmYwNmYzZmY2YTAxNjg1MWY0NWMzOTg3MzJiYzUwYw==
	// File-Size: 14
	// File-Name: tmp.txt
	// File-Type: text/plain

	request, err := http.NewRequest("POST", "http://localhost:8080/upl", file)
	if err != nil {
		return err
	}
	request.Header.Add("File-Hash", "MmU5OTc1ODU0ODk3MmE4ZTg4MjJhZDQ3ZmExMDE3ZmY3MmYwNmYzZmY2YTAxNjg1MWY0NWMzOTg3MzJiYzUwYw==")
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

	// Expect:
	response := UploadResponse{
		Duplicate: false,
		URL:       "", // Not able to get this as it is randomly generated. It should never be empty in response, but it will be here.
		File: File{
			FileName:    "tmp.txt",
			Size:        14,
			ContentType: "text/plain",
		},
	}

	response.File.UUID = sha256.Sum256([]byte("this is a test"))
	response.File.uuidName = encodeToString(response.File.UUID[:])
	////////////

	var f *File
	err = json.Unmarshal(content, &f)
	if err != nil {
		return err
	}

	// We need to check if the response content contains the correct information
	if f.FileName != "tmp.txt" {
		return errors.New("file names do not match")
	}

	// Best we can do is verify the block comes back with most of the right information

	return nil
}

func TestHandleUpload(t *testing.T) {
	fServer := NewFileServer()
	fServer.Initialize()

	m := http.NewServeMux()
	server := http.Server{
		Addr:    ":8080",
		Handler: m,
	}

	go func() {
		// Sleep for a couple seconds to let the listener start
		time.Sleep(2 * time.Second)

		err := postRequest()
		if err != nil {
			t.Error(err)
		}

		if err := server.Shutdown(context.Background()); err != nil {
			t.Error(err)
		}

		// t.Log("Response is formed correctly")
	}()

	m.HandleFunc("/upl", fServer.HandleUpload)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		t.Error(err)
	}

	time.Sleep(2 * time.Second) // Block for enough time for test go routine can finish
}
