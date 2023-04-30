package main

type UploadResponse struct {
	URL       string
	Duplicate bool
	File      File
}
