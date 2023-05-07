package fileserver

import (
	"net/http"
	"strings"

	"buckette/logger"
)

func (fServer *FileServer) HandleLS(rw http.ResponseWriter, req *http.Request) {
	logger.Connection(req)
	arguments := strings.Split(req.RequestURI, "/")
	arguments = arguments[2:]
}
