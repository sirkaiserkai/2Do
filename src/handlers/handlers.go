package handlers

import (
	"auth"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"models"
	"net/http"
)

// jsonResponse is the struct for almost all responses
type jsonResponse struct {
	result       string      `json:"result,omitempty"`
	errorMessage string      `json:"error_message,omitempty"`
	data         interface{} `json:"data,omitempty"`
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Nothing to see here move along"))
}

// TodosHandler is a handler function for the /api/todos endpoint
// it acts as a multiplexer to a respective http method handler.
func TodosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		TodosGetHandler(w, r)
	case "POST":
		TodosPostHandler(w, r)
	}
}

// TodoHandler is a handler function for the /api/todos/{id} endpoint
// it acts as a multiplexer to a respective http method handler.
func TodoHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		TodoGetHandler(w, r)
	case "PUT":
		TodoPutHandler(w, r)
	case "DELETE":
		TodoDeleteHandler(w, r)
	}
}

// TodosGetHandler is the handler function which returns all the
// respective todos for a user.
func TodosGetHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := auth.GetClaims(r)
	if err != nil {
		UnauthorizedHandler(w, r, "")
		return
	}

	tds := models.NewTodoDataStore()

	ts, err := tds.GetTodosForUserId(claims.UserId)
	if err != nil {
		NotFoundHandler(w, r, "No 2Dos found.")
		log.Println("Failed to get Todos: " + err.Error())
		return
	}

	data, err := json.Marshal(ts)
	if err != nil {
		InternalErrorHandler(w, r, "")
		log.Println("Failed to get Todos: " + err.Error())
		return
	}

	w.Write(data)
}

// TodosPostHandler is the handler function so that a user
// can insert new todos.
func TodosPostHandler(w http.ResponseWriter, r *http.Request) {
	t := models.NewTodo()
	claims, err := auth.GetClaims(r)
	if err != nil {
		UnauthorizedHandler(w, r, "")
		return
	}

	err = json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		BadRequestHandler(w, r, "Body format incorrect for 2Do. Try: { \"title\": \"Some Title\", \"note\": \"Example note\" }")
		log.Println(err)
		return
	}

	t.Ownerid = claims.UserId

	tds := models.NewTodoDataStore()
	err = tds.InsertTodo(t)
	if err != nil {
		InternalErrorHandler(w, r, "Failure to add 2Do")
		log.Println("Failure to add 2Do: " + err.Error())
		return
	}
	log.Println("2Do: " + t.Id.String() + " created")

	res := jsonResponse{
		result: fmt.Sprintf("Successfully created 2Do: %s", t.Id.String()),
		data:   t,
	}

	msg, err := json.Marshal(res)
	if err != nil {
		InternalErrorHandler(w, r, "")
		log.Println(fmt.Sprintf("TodosPostHandler: %s", err.Error()))
		return
	}

	w.WriteHeader(201)
	w.Write(msg)
}

// TodoGetHandler is the handler function in order to retrieve a
// specific todo with an ID.
func TodoGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	claims, err := auth.GetClaims(r)
	if err != nil {
		UnauthorizedHandler(w, r, "")
		log.Println("Authentication token invalid")
		return
	}

	tds := models.NewTodoDataStore()

	t, err := tds.GetTodoById(id)
	if err != nil {
		NotFoundHandler(w, r, fmt.Sprintf("Failed to retrieve 2Do with id: %s", id))
		log.Println(fmt.Sprintf("TodoGetHandler Failure: 2Do (%s) not found", id))
		return
	}

	if t.Ownerid != claims.UserId {
		NotFoundHandler(w, r, fmt.Sprintf("Failed to retrieve 2Do with id: %s", id))
		log.Println(fmt.Sprintf("TodoGetHandler Failure: t.Ownerid (%s) != claims.UserId (%s)", t.Ownerid, claims.UserId))
		return
	}

	data, err := json.Marshal(t)
	if err != nil {
		NotFoundHandler(w, r, fmt.Sprintf("Failed to retrieve 2Do with id: %s", id))
		log.Println("TodoHandler Error: " + err.Error())
		return
	}

	w.Write(data)
}

// TodoPutHandler is the handler function which allows a user
// to modify an existing todo.
func TodoPutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	claims, err := auth.GetClaims(r)
	if err != nil {
		UnauthorizedHandler(w, r, "")
		return
	}

	m := make(map[string]interface{})
	err = json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		InternalErrorHandler(w, r, "Failure to modify 2Do.")
		log.Println("Failure to modify 2Do: " + err.Error())
		return
	}

	tds := models.NewTodoDataStore()

	err = tds.ModifyTodo(id, claims.UserId, m)
	if err != nil {
		NotFoundHandler(w, r, "2Do not found.")
		return
	}

	res := jsonResponse{result: fmt.Sprintf("Successfully modified 2Do: %s", id)}
	msg, err := json.Marshal(res)
	if err != nil {
		InternalErrorHandler(w, r, "Failure to modify 2Do")
		log.Println("Failure to modify 2Do: " + err.Error())
		return
	}

	w.WriteHeader(200)
	w.Write(msg)
}

func TodoDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	claims, err := auth.GetClaims(r)
	if err != nil {
		UnauthorizedHandler(w, r, "")
		return
	}

	tds := models.NewTodoDataStore()
	err = tds.DeleteTodo(id, claims.UserId)
	if err != nil {
		NotFoundHandler(w, r, "2Do not found.")
		return
	}

	res := jsonResponse{result: fmt.Sprintf("Successfully deleted 2Do: %s", id)}
	msg, err := json.Marshal(res)
	if err != nil {
		InternalErrorHandler(w, r, "Failure to delete 2Do")
		log.Println("Failure to delete 2Do: " + err.Error())
		return
	}

	w.WriteHeader(200)
	w.Write(msg)
}
