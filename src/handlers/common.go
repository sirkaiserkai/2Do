package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const (
	ContentType     = "Content-Type"
	ApplicationJSON = "application/json"
)

const (
	StatusSuccess       = 200
	StatusCreation      = 201
	StatusBadRequest    = 400
	StatusUnauthorized  = 401
	StatusNotFound      = 404
	StatusInternalError = 500
)

// GENERIC REQUEST HANDLERS //
func NotFoundHandler(w http.ResponseWriter, r *http.Request, errMsg string) {
	//return func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(StatusNotFound)
	w.Header().Set(ContentType, ApplicationJSON)
	v := jsonResponse{ErrorMessage: errMsg}

	msg, err := json.Marshal(v)
	if err != nil || errMsg == "" {
		w.Write([]byte("{ \"error_message\": \"Failure to retrieve objects\"}"))
		if err != nil {
			log.Println(fmt.Sprintf("Failure to marshal jsonResponse: %v", err))
		}
		return
	}

	w.Write(msg)
	//}
}

func InternalErrorHandler(w http.ResponseWriter, r *http.Request, errMsg string) {
	w.WriteHeader(StatusInternalError)
	w.Header().Set(ContentType, ApplicationJSON)
	v := jsonResponse{ErrorMessage: errMsg}

	msg, err := json.Marshal(v)
	if err != nil || errMsg == "" {
		w.Write([]byte("{ \"error_message\": \"The server encountered an unexpected condition which prevented it from fulfilling the request.\"}"))
		if err != nil {
			log.Println(fmt.Sprintf("Failure to marshal jsonResponse: %v", err))
		}
		return
	}

	w.Write(msg)
}

func UnauthorizedHandler(w http.ResponseWriter, r *http.Request, errMsg string) {
	w.WriteHeader(StatusUnauthorized)
	w.Header().Set("WWW-Authenticate", "Bearer")
	w.Header().Set(ContentType, ApplicationJSON)
	v := jsonResponse{ErrorMessage: errMsg}

	msg, err := json.Marshal(v)
	if err != nil || errMsg == "" {
		w.Write([]byte("{ \"error_message\": \"The request requires user authentication.\"}"))
		if err != nil {
			log.Println(fmt.Sprintf("Failure to marshal jsonResponse: %v", err))
		}
		return
	}

	w.Write(msg)
}

func BadRequestHandler(w http.ResponseWriter, r *http.Request, errMsg string) {
	w.WriteHeader(StatusBadRequest)
	w.Header().Set(ContentType, ApplicationJSON)
	v := jsonResponse{ErrorMessage: errMsg}

	msg, err := json.Marshal(v)
	if err != nil || errMsg == "" {
		w.Write([]byte("{ \"error_message\": \"Bad Request.\"}"))
		if err != nil {
			log.Println(fmt.Sprintf("Failure to marshal jsonResponse: %v", err))
		}
		return
	}

	w.Write(msg)
}
