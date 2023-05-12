package fileserver

import (
	"encoding/json"
	"net/http"

	"buckette/logger"
)

func (fServer *FileServer) HandleLS(rw http.ResponseWriter, req *http.Request) {
	logger.Connection(req)

	response, err := json.Marshal(fServer.Files)
	if err != nil {
		logger.Fatal(
			"Unable to marshal files cache to json",
			"fServer.HandleLS",
			err,
		)
		rw.WriteHeader(500)
		return
	}

	rw.Write(response)

}
