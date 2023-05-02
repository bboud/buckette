package main

import (
	"fmt"
	"net/http"
	"os"
)

func (fServer *FileServer) HandleDownload(rw http.ResponseWriter, req *http.Request) {
	LogConnection(req)

	file := fServer.FindByURL(req.RequestURI[1:])
	if file != nil {
		rw.Header().Set("Content-Type", file.ContentType)
		rw.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", file.FileName))

		fileName := encodeToString(file.UUID[:])
		content, err := os.ReadFile(FileStoreDir + fileName)
		if err != nil {
			LogWarning(
				"Unable to read file "+file.FileName,
				"Attempting to send file to requester",
				err,
			)
			rw.WriteHeader(500)
			return
		}
		rw.Write(content)
		return
	}

	// Serve main website
	http.FileServer(http.Dir(SiteFiles)).ServeHTTP(rw, req)
}
