package main

import (
	"github.com/gorilla/mux"
	"handlers"
	"log"
	"logger"
	"net/http"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile) // Loggers print the execution line
}

const (
	homeRoute  = "/"
	todosRoute = "/todos"
	todoRoute  = "/todos/{id}"
)

func main() {

	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()

	homeHandler := logger.Logger(handlers.HomeHandler, homeRoute)
	todosHandler := logger.Logger(handlers.TodosHandler, todosRoute)
	todoHandler := logger.Logger(handlers.TodosHandler, todoRoute)

	api.HandleFunc(homeRoute, homeHandler).Methods("GET")
	api.HandleFunc(todosRoute, todosHandler).Methods("GET", "POST")
	api.HandleFunc(todoRoute, todoHandler).Methods("GET", "PUT", "DELETE")

	srv := &http.Server{
		Handler:      r,
		Addr:         "localhost:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
