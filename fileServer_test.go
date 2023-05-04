package main

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"os"
	"testing"
)

// fServer := newFileServer()
// go fServer.handleNewFiles()
// fServer.initialize()
func TestNewFileServer(t *testing.T) {
	result := NewFileServer()

	if len(result.newFile) == 0 && len(result.Files) == 0 && len(result.URLs) == 0 {
		t.Log("The file server has been initialized correctly")
	} else {
		t.Error("The file server has not been initialized correctly")
	}
}

func TestNewFile(t *testing.T) {
	fServer := NewFileServer()
	fServer.Initialize()

	// Put a file into the file cache
	url, err := fServer.GenerateURL()
	if err != nil {
		t.Error(err)
	}

	err = os.WriteFile(TmpDir+"DAT_"+url, []byte("This is a test file"), 0644)
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(TmpDir + "DAT_" + url)

	uuid := sha256.Sum256([]byte("This is a test file"))
	fServer.Push(File{
		URL:  url,
		UUID: uuid,
	})

	// Now see if making a new file from the url.
	file, err := fServer.NewFile(url)
	if file == nil {
		t.Error(err)
	}

	var fileExists *FileExists
	if errors.Is(err, fileExists) {
		if bytes.Equal(uuid[:], file.UUID[:]) {
			t.Log("Success!")
		} else {
			t.Log(file.URL)
			t.Error("The file did not return with the correct uuid")
		}
	} else {
		if bytes.Equal(uuid[:], file.UUID[:]) {
			t.Log("Success!")
		} else {
			t.Log(file.URL)
			t.Error("The file did not return with the correct uuid")
		}
	}

	//////////////

	// Now check the case where the file does not exist and we want a blank new file
	// Put the file in the temp store
	url, err = fServer.GenerateURL()
	if err != nil {
		t.Error(err)
	}

	// Put the file in the temp store
	_, err = os.Create(TmpDir + "DAT_" + url)
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(TmpDir + "DAT_" + url)

	file, err = fServer.NewFile(url)
	if err != nil {
		t.Error(err)
	}
	if file.URL == "" {
		t.Log("Success!")
	} else {
		t.Error("The file did not return with the correct uuid")
	}
}

func TestGenerateURL(t *testing.T) {
	fServer := NewFileServer()
	fServer.Initialize()

	url, err := fServer.GenerateURL()
	if err != nil {
		t.Error(err)
		return
	}
	if len(url) != 16 {
		t.Error("Incorrect url length")
		return
	}

	t.Log("Generated URL correctly!")
}

func TestExists(t *testing.T) {
	fServer := NewFileServer()
	fServer.Initialize()

	// Put onto the cache
	uuid := sha256.Sum256([]byte("This is a test file"))
	fServer.Push(File{
		UUID: uuid,
	})

	// Check response
	uuid2 := sha256.Sum256([]byte("This is a test file"))
	if fServer.Exists(uuid2[:]) {
		t.Log("Successfully found file on server")
	} else {
		t.Error("Should have found file but didnt")
	}

	uuid3 := sha256.Sum256([]byte("This is ANOTHER test file"))
	if !fServer.Exists(uuid3[:]) {
		t.Log("Successfully did not find file on server")
	} else {
		t.Error("Should have not found file")
	}
}

func TestFindByURL(t *testing.T) {
	fServer := NewFileServer()
	fServer.Initialize()

	url, err := fServer.GenerateURL()
	if err != nil {
		t.Error(err)
	}

	fServer.Push(File{
		URL: url,
	})

	//Now we find
	found := fServer.FindByURL(url)
	if found.URL == url {
		t.Log("Successfully found file by URL")
	} else {
		t.Error("Did not find file successfully with URL")
	}
}
