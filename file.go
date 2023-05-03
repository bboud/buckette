package main

import (
	"time"
)

type File struct {
	FileName      string
	UUID          [32]byte
	URL           string
	Size          int64
	ContentType   string
	Uploaded      time.Time
	UserUploaded  string
	DownloadCount int
}
