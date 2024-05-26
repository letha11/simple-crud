package api

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

type ErrorResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

type NoDataResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

type GenericSuccessReponse[T any] struct {
	StatusCode int `json:"status_code"`
	Data       T   `json:"data"`
}

func writeError(w http.ResponseWriter, message string, code int) {
	err := ErrorResponse{
		StatusCode: code,
		Message:    message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(err)
}

var (
	RequestErrorHandler = func(w http.ResponseWriter, err error, code int) {
		writeError(w, err.Error(), code)
	}
	InternalErrorHandler = func(w http.ResponseWriter, err any) {
		logrus.Error(err)
		writeError(w, "An Unexpected Error Occured.", http.StatusInternalServerError)
	}
)

func writeSuccessResponse(w http.ResponseWriter, code int, response interface{}) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

var (
	GenericResponseHandler = func(w http.ResponseWriter, code int, data any) {
		resp := GenericSuccessReponse[any]{StatusCode: code, Data: data}
		writeSuccessResponse(w, code, resp)
	}
	NoDataResponseHandler = func(w http.ResponseWriter, code int, message string) {
		writeSuccessResponse(w, code, NoDataResponse{StatusCode: code, Message: message})
	}
)
