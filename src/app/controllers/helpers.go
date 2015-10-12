package controllers

import (
	"encoding/json"
	"log"
	"net/http"
)

type errorResponse struct {
	ErrCode    int    `json:"error_code"`
	ErrMessage string `json:"error"`
}

// sendError sends and logs the given error.
func sendError(errorCode int, err error, w http.ResponseWriter) {
	w.WriteHeader(errorCode)
	log.Printf("Error: Returning status code %d with error message %s.\n", errorCode, err)
	resp := errorResponse{ErrCode: errorCode, ErrMessage: err.Error()}
	json.NewEncoder(w).Encode(resp)
}
