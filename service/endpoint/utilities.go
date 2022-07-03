package endpoint

import (
	"encoding/json"
	"log"
	"net/http"
)

type errorResponse struct {
	Error string `json:"error"`
}

func processError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := &errorResponse{message}

	writeObj(w, response)
}

func logError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func writeObj(w http.ResponseWriter, obj any) {
	marshal, err := json.Marshal(obj)
	if err != nil {
		logError(err)
		processError(w, err.Error(), 500)
		return
	}

	_, err = w.Write(marshal)
	if err != nil {
		logError(err)
		processError(w, err.Error(), 500)
		return
	}
}
