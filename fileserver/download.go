package fileserver

import (
	"fmt"
	"net/http"
	"os"

	"buckette/logger"
)

const (
	SiteFiles = "./frontend/dist"
)

func (fServer *FileServer) HandleDownload(rw http.ResponseWriter, req *http.Request) {
	logger.Connection(req)

	file := fServer.findByURL(req.RequestURI[1:])
	if file != nil {
		rw.Header().Set("Content-Type", file.ContentType)
		rw.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", file.FileName))

		content, err := os.ReadFile(FileStoreDir + file.UUID)
		if err != nil {
			logger.Warning(
				"Unable to read file "+file.FileName,
				"fileserver.HandleDownload",
				err,
			)
			rw.WriteHeader(500)
			return
		}
		file.DownloadCount += 1
		rw.Write(content)
		return
	}

	// Serve main website
	http.FileServer(http.Dir(SiteFiles)).ServeHTTP(rw, req)
}
