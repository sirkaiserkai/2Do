package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"models"
	"net/http"
)

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
	tds := models.NewTodoDataStore()
	ts, err := tds.GetAllTodos()
	if err != nil {
		w.Write([]byte("Failed to get Todos"))
		log.Println(err)
		return
	}

	data, err := json.Marshal(ts)
	if err != nil {
		w.Write([]byte("Failed to get Todos"))
		log.Println(err)
		return
	}

	w.Write(data)

}

// TodosPostHandler is the handler function so that a user
// can insert new todos.
func TodosPostHandler(w http.ResponseWriter, r *http.Request) {
	t := models.NewTodo()

	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		w.Write([]byte("Body incorrect"))
		log.Println(err)
		return
	}

	tds := models.NewTodoDataStore()

	err = tds.InsertTodo(t)
	if err != nil {
		w.Write([]byte("Failed to Add"))
		log.Println(err)
		return
	}

	log.Println("Todo: " + t.Id.String() + " created")
	w.Write([]byte("Todo: " + t.Id.String() + " created"))
}

// TodoGetHandler is the handler function in order to retrieve a
// specific todo with an ID.
func TodoGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	tds := models.NewTodoDataStore()

	t, err := tds.GetTodoById(id)
	if err != nil {
		w.Write([]byte("Failure to get todo"))
		log.Println("TodoHandler Error: " + err.Error())
		return
	}

	data, err := json.Marshal(t)
	if err != nil {
		w.Write([]byte("Failure to get todo"))
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

	m := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		panic(err)
	}

	log.Println(m)
	tds := models.NewTodoDataStore()
	err = tds.ModifyTodo(id, m)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("Todo not found"))
		return
	}

	w.WriteHeader(200)
	w.Write([]byte("Todo: " + id + " successfully modified"))
}

func TodoDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	tds := models.NewTodoDataStore()
	err := tds.DeleteTodo(id)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("Todo not found"))
		return
	}

	w.WriteHeader(200)
	w.Write([]byte("Todo: " + id + " successfully deleted"))
}
