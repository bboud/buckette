package main

import (
	"crypto/sha256"
	"errors"
	"os"
	"strings"
	"testing"
)

func TestCopyToFileSystem(t *testing.T) {
	fServer := NewFileServer()
	fServer.Initialize()

	url, err := fServer.GenerateURL()
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
	fServer.Initialize()

	url, err := fServer.GenerateURL()
	if err != nil {
		t.Error(err)
	}

	file := File{
		UUID: sha256.Sum256([]byte("Some data")),
		URL:  url,
	}

	uuidName := encodeToString(file.UUID[:])

	err = copyToFileSystem(strings.NewReader("This is a string!"), url)
	if err != nil {
		t.Error(err)
	}

	err = writeRecord(&file)
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
	fServer.Initialize()

	url, _ := fServer.GenerateURL()

	err := os.WriteFile(TmpDir+url, []byte("Some stuff"), 0644)
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
