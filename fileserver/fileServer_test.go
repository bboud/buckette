package fileserver

import (
	"crypto/sha256"
	"errors"
	"os"
	"strings"
	"testing"
)

func TestNewFileServer(t *testing.T) {
	result := NewFileServer()

	if result.Files != nil && result.URLs != nil {
		t.Log("The file server has been initialized correctly")
	} else {
		t.Error("The file server has not been initialized correctly")
	}
}

func TestNewFile(t *testing.T) {

	var err error

	// If it exists, clear it
	err = os.RemoveAll("./buckette-data")
	if !errors.Is(err, os.ErrNotExist) && err != nil {
		t.Log(err)
	}

	makeDataStore()

	fServer := NewFileServer()
	fServer.initialize()

	// Put a file into the file cache
	url, err := fServer.generateURL()
	if err != nil {
		t.Error(err)
	}

	err = os.WriteFile(TmpDir+"DAT_"+url, []byte("This is a test file"), 0644)
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(TmpDir + "DAT_" + url)

	uuid := sha256.Sum256([]byte("This is a test file"))
	fServer.push(&File{
		URL:  url,
		UUID: encodeToString(uuid[:]),
	})

	// Now see if making a new file from the url.
	file, err := fServer.NewFile(url)
	if file == nil {
		t.Error(err)
	}

	var fileExists *FileExists
	if errors.Is(err, fileExists) {
		// bytes.Equal(uuid[:], file.UUID.UUIDByte[:])
		if strings.Compare(encodeToString(uuid[:]), file.UUID) == 0 {
			t.Log("Success!")
		} else {
			t.Log(file.URL)
			t.Error("The file did not return with the correct uuid")
		}
	} else {
		//bytes.Equal(uuid[:], file.UUID.UUIDByte[:])
		if strings.Compare(file.UUID, encodeToString(uuid[:])) == 0 {
			t.Log("Success!")
		} else {
			t.Log(file.URL)
			t.Error("The file did not return with the correct uuid")
		}
	}

	//////////////

	// Now check the case where the file does not exist and we want a blank new file
	// Put the file in the temp store
	url, err = fServer.generateURL()
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
	fServer.initialize()

	url, err := fServer.generateURL()
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
	fServer.initialize()

	// Put onto the cache
	uuid := sha256.Sum256([]byte("This is a test file"))
	fServer.push(&File{
		UUID: encodeToString(uuid[:]),
	})

	// Check response
	uuid2 := sha256.Sum256([]byte("This is a test file"))
	if fServer.exists(encodeToString(uuid2[:])) {
		t.Log("Successfully found file on server")
	} else {
		t.Error("Should have found file but didnt")
	}

	uuid3 := sha256.Sum256([]byte("This is ANOTHER test file"))
	if !fServer.exists(encodeToString(uuid3[:])) {
		t.Log("Successfully did not find file on server")
	} else {
		t.Error("Should have not found file")
	}
}

func TestFindByURL(t *testing.T) {
	fServer := NewFileServer()
	fServer.initialize()

	url, err := fServer.generateURL()
	if err != nil {
		t.Error(err)
	}

	fServer.push(&File{
		URL: url,
	})

	//Now we find
	found := fServer.findByURL(url)
	if found.URL == url {
		t.Log("Successfully found file by URL")
	} else {
		t.Error("Did not find file successfully with URL")
	}
}

func TestMakeDataStore(t *testing.T) {

	var notExists bool
	var err error

	// If it exists, clear it
	err = os.RemoveAll("./buckette-data")
	if err != nil {
		t.Log(err)
	}

	err = makeDataStore()
	if err != nil {
		t.Error(err)
	}

	_, err = os.Stat(FileStoreDir)
	notExists = os.IsNotExist(err)
	if notExists {
		t.Error(err)

	} else {
		t.Log("File store dir exists")
	}

	_, err = os.Stat(RecordStoreDir)
	notExists = os.IsNotExist(err)
	if notExists {
		t.Error(err)

	} else {
		t.Log("File store dir exists")
	}

	_, err = os.Stat(TmpDir)
	notExists = os.IsNotExist(err)
	if notExists {
		t.Error(err)
	} else {
		t.Log("File store dir exists")
	}
}

func TestLoadFromDisk(t *testing.T) {
	TestHandleUpload(t)

	fServer := NewFileServer()

	err := fServer.loadFromDisk()
	if err != nil {
		t.Error(err)
	}

	file := fServer.findByUUID("Lpl1hUiXKo6IIq1H+hAX*3Lwbz*2oBaFH0XDmHMrxQw")
	if file == nil {
		t.Error("File came back nil and thus was not found")
	} else {
		t.Log("Successfully loaded from disk")
	}
}

func TestInitialize(t *testing.T) {
	TestMakeDataStore(t)

	TestLoadFromDisk(t)
}

func TestFindByUUID(t *testing.T) {
	fServer := NewFileServer()

	fServer.push(&File{
		UUID: "not an actual uuid",
	})

	file := fServer.findByUUID("not an actual uuid")
	if file == nil {
		t.Error("Was not able to correctly find the uuid")
	}
}
