package api

import (
	"encoding/json"
	"net/http"

	"github.com/simple-crud-go/internal/models"
	"github.com/sirupsen/logrus"
)

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type NoDataResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type RegisterSuccessResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

type GenericSuccessReponse[T any] struct {
	Status string `json:"status"`
	Data   T      `json:"data"`
}

func writeError(w http.ResponseWriter, message string, code int) {
	err := ErrorResponse{
		Status:  "error",
		Message: message,
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
		resp := GenericSuccessReponse[any]{Status: "success", Data: data}
		writeSuccessResponse(w, code, resp)
	}
	NoDataResponseHandler = func(w http.ResponseWriter, code int, message string) {
		writeSuccessResponse(w, code, NoDataResponse{Status: "success", Message: message})
	}
)
