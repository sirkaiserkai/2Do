package main

import (
	"github.com/gorilla/mux"
	"handlers"
	"logger"

	"log"
	"net/http"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile) // Loggers print the execution line. SUPER clutch everyone should use this for debugging
}

const (
	homeRoute = "/"

	loginRoute  = "/login"
	signUpRoute = "/signup"

	todosRoute = "/todos"
	todoRoute  = "/todos/{id}"

	usrAccntRoute = "/account"
)

const addr = "localhost:8000"

func main() {

	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()

	homeHandler := logger.Logger(handlers.ValidatePath(handlers.HomeHandler), homeRoute)
	todosHandler := logger.Logger(handlers.ValidatePath(handlers.TodosHandler), todosRoute)
	todoHandler := logger.Logger(handlers.ValidatePath(handlers.TodoGetHandler), todoRoute)

	signUpHandler := logger.Logger(handlers.SignUpHandler, signUpRoute)
	logInHandler := logger.Logger(handlers.LogInHandler, loginRoute)

	api.HandleFunc(homeRoute, homeHandler).Methods("GET")
	api.HandleFunc(todosRoute, todosHandler).Methods("GET", "POST")
	api.HandleFunc(todoRoute, todoHandler).Methods("GET", "PUT", "DELETE")

	api.HandleFunc(signUpRoute, signUpHandler).Methods("POST")
	api.HandleFunc(loginRoute, logInHandler).Methods("POST")

	srv := &http.Server{
		Handler:      r,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Hosting on %s", addr)
	log.Fatal(srv.ListenAndServe())
}
