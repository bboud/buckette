package main

import (
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

// func TestPush(t *testing.T) {
// 	fServer := newFileServer()

// 	f := File{
// 		FileName:      "file",
// 		Uploaded:      time.Date(2000, 11, 11, 22, 22, 22, 22, time.FixedZone("EST", 0)),
// 		Size:          0,
// 		DownloadCount: 0,
// 		ContentType:   "content",
// 	}

// 	// // Processing the request
// 	// _, err := f.HandleUploadPart()
// 	// if err != nil {
// 	// 	t.Error(err)
// 	// }

// 	// fServer.Push(f)

// 	// result := fServer.Files[]
// }
