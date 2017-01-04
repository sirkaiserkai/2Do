package main

import (
	"github.com/gorilla/mux"
	"handlers"
	"log"
	"net/http"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {

	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/", handlers.HomeHandler)
	api.HandleFunc("/todos", handlers.TodosHandler).Methods("GET", "POST")
	api.HandleFunc("/todos/{id}", handlers.TodoHandler).Methods("GET", "PUT", "DELETE")

	srv := &http.Server{
		Handler:      r,
		Addr:         "localhost:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
