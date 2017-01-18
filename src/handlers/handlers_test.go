package handlers

import (
	"auth"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const (
	testDB = "2Do_handlers_test_DB"
)

func init() {
	models.USER_STORE_TYPE = models.Test
	models.TODO_STORE_TYPE = models.Test
}

// Holy setup function, Batman!
func handlersSetup(method, route, body string) (*http.Request, *httptest.ResponseRecorder) {
	//models.USER_STORE_TYPE = models.TestDataStoreType
	buf := bytes.NewBufferString(body)
	req, err := http.NewRequest(method, route, buf)
	if err != nil {
		log.Fatal(err)
	}

	rr := httptest.NewRecorder()

	return req, rr
}

func TestTodosHandler0(t *testing.T) {
	req, rr := handlersSetup("GET", "api/todos", "")

	u := models.NewUser()
	tus := models.NewUserStorage()
	tus.InsertUser(u)
	token, err := auth.CreateToken(u)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		TodosHandler(w, r)
	}

	ValidatePath(handler).ServeHTTP(rr, req)

	testStatus(StatusSuccess, rr, t)

	expected := "[]"
	testBody(expected, rr, t)
}

func TestTodosHandler1(t *testing.T) {
	req, rr := handlersSetup("GET", "api/todos", "")

	u := models.NewUser()
	tus := models.NewUserStorage()
	tus.InsertUser(u)
	token, err := auth.CreateToken(u)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	exampleTime := time.Now()
	t0 := models.NewTodo()
	t0.Ownerid = u.Id.Hex()
	t0.Created = exampleTime
	t0.Due = exampleTime
	tds := models.NewTodoStorage()
	tds.InsertTodo(t0)

	var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		TodosHandler(w, r)
	}

	ValidatePath(handler).ServeHTTP(rr, req)

	testStatus(StatusSuccess, rr, t)

	/*expected := fmt.Sprintf("[{\"id\":\"%s\","+
		"\"title\":\"%s\","+
		"\"note\":\"%s\""+
		"\"created_date\":\"%v\""+
		"\"due_date\":\"%v\""+
		"\"completed\":%v}]", t0.Id.String(), t0.Title, t0.Note,
		t0.Created.String(), t0.Due.String(), t0.Completed)
	expected = strings.Replace(expected, " ", "", -1)*/

	b, err := json.Marshal(t0)
	if err != nil {
		log.Fatal(err)
	}

	testBody(fmt.Sprintf("[%s]", string(b)), rr, t)

	tds.DeleteTodo(t0.Id.Hex(), u.Id.Hex())
	tus.DeleteUser(u.Id.Hex())
}

func TestTodosHandler2(t *testing.T) {
	u := models.NewUser()
	tus := models.NewUserStorage()
	tus.InsertUser(u)
	token, err := auth.CreateToken(u)
	if err != nil {
		log.Fatal(err)
	}

	exampleTime := time.Now()
	t0 := models.NewTodo()
	t0.Ownerid = u.Id.Hex()
	t0.Created = exampleTime
	t0.Due = exampleTime

	b, err := json.Marshal(t0)
	if err != nil {
		log.Fatal(err)
	}

	body := string(b)
	req, rr := handlersSetup("POST", "api/todos", body)

	req.Header.Set("Authorization", "Bearer "+token)

	var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		TodosHandler(w, r)
	}

	ValidatePath(handler).ServeHTTP(rr, req)

	testStatus(StatusCreation, rr, t)

	res := jsonResponse{
		Result: fmt.Sprintf("Successfully created 2Do: %s", t0.Id.String()),
		Data:   t0,
	}

	msg, err := json.Marshal(res)
	if err != nil {
		log.Fatal(err)
	}

	testBody(string(msg), rr, t)

	tds := models.NewTodoDataStore()
	tds.DeleteTodo(t0.Id.Hex(), u.Id.Hex())
	tus.DeleteUser(u.Id.Hex())
}
