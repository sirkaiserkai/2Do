package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// GENERIC REQUEST HANDLERS //
func NotFoundHandler(w http.ResponseWriter, r *http.Request, errMsg string) {
	w.WriteHeader(404)
	w.Header().Set("Content-Type", "application/json")
	v := jsonResponse{errorMessage: errMsg}

	msg, err := json.Marshal(v)
	if err != nil || errMsg == "" {
		w.Write([]byte("{ \"error_message\"; \"Failure to retrieve objects\"}"))
		if errMsg != "" {
			log.Println(fmt.Sprintf("Failure to marshal jsonResponse: %v", err))
		}
		return
	}

	w.Write(msg)
}

func InternalErrorHandler(w http.ResponseWriter, r *http.Request, errMsg string) {
	w.WriteHeader(500)
	w.Header().Set("Content-Type", "application/json")
	v := jsonResponse{errorMessage: errMsg}

	msg, err := json.Marshal(v)
	if err != nil || errMsg == "" {
		w.Write([]byte("{ \"error_message\"; \"The server encountered an unexpected condition which prevented it from fulfilling the request.\"}"))
		if errMsg != "" {
			log.Println(fmt.Sprintf("Failure to marshal jsonResponse: %v", err))
		}
		return
	}

	w.Write(msg)
}

func UnauthorizedHandler(w http.ResponseWriter, r *http.Request, errMsg string) {
	w.WriteHeader(401)
	w.Header().Set("WWW-Authenticate", "Bearer")
	w.Header().Set("Content-Type", "application/json")
	v := jsonResponse{errorMessage: errMsg}

	msg, err := json.Marshal(v)
	if err != nil || errMsg == "" {
		w.Write([]byte("{ \"error_message\"; \"The request requires user authentication.\"}"))
		if errMsg != "" {
			log.Println(fmt.Sprintf("Failure to marshal jsonResponse: %v", err))
		}
		return
	}

	w.Write(msg)
}

func BadRequestHandler(w http.ResponseWriter, r *http.Request, errMsg string) {
	w.WriteHeader(400)
	w.Header().Set("Content-Type", "application/json")
	v := jsonResponse{errorMessage: errMsg}

	msg, err := json.Marshal(v)
	if err != nil || errMsg == "" {
		w.Write([]byte("{ \"error_message\"; \"Bad Request.\"}"))
		if errMsg != "" {
			log.Println(fmt.Sprintf("Failure to marshal jsonResponse: %v", err))
		}
		return
	}

	w.Write(msg)
}
