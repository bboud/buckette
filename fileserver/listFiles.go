package fileserver

import (
	"net/http"
	"strings"

	"github.com/bboud/buckette/logger"
)

func (fServer *FileServer) HandleLS(rw http.ResponseWriter, req *http.Request) {
	logger.LogConnection(req)
	arguments := strings.Split(req.RequestURI, "/")
	arguments = arguments[2:]
}
