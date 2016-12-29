package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	// "mdb"
	"models"
	"net/http"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	//defer mdb.Session.Close()

	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/", HomeHandler)
	api.HandleFunc("/todos", TodosHandler).Methods("GET", "POST")
	api.HandleFunc("/todos/{id}", TodoHandler).Methods("GET", "PUT")

	srv := &http.Server{
		Handler:      r,
		Addr:         "localhost:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Nothing to see here move along"))
}

func TodosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		TodosGetHandler(w, r)
	case "POST":
		TodosPostHandler(w, r)
	}
}

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

func TodoHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		TodoGetHandler(w, r)
	case "PUT":
		TodoPutHandler(w, r)
	}
}

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

func TodoPutHandler(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//_ := vars["id"]

}
