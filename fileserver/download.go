package fileserver

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bboud/buckette/logger"
)

const (
	SiteFiles = "./frontend/dist"
)

func (fServer *FileServer) HandleDownload(rw http.ResponseWriter, req *http.Request) {
	logger.LogConnection(req)

	file := fServer.findByURL(req.RequestURI[1:])
	if file != nil {
		rw.Header().Set("Content-Type", file.ContentType)
		rw.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", file.FileName))

		content, err := os.ReadFile(FileStoreDir + file.UUID)
		if err != nil {
			logger.LogWarning(
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
