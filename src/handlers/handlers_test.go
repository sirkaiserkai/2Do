package handlers

import (
	"auth"
	"gopkg.in/mgo.v2"
	"log"
	"mdb"
	"models"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	testDB = "2Do_handlers_test_DB"
)

// Holy setup function, Batman!
func handlersSetup(method, route string) (*http.Request, *httptest.ResponseRecorder, models.UserDataStore, models.TodoDataStore, models.User) {
	mdb.DatabaseName = testDB

	req, err := http.NewRequest(method, route, nil)
	if err != nil {
		log.Fatal(err)
	}

	rr := httptest.NewRecorder()

	uds := models.NewUserDataStore()
	tds := models.NewTodoDataStore()

	uds.SetDB(testDB)
	tds.SetDB(testDB)

	u := models.NewUser()
	log.Println(u.Id.Hex())
	err = uds.InsertUser(u)
	if err != nil {
		log.Fatal(err)
	}

	token, err := auth.CreateToken(u)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	return req, rr, uds, tds, u
}

func udsTeardown(uds models.UserDataStore) {
	defer uds.Close()

	session, err := mgo.Dial(mdb.Hostname)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	err = session.DB(uds.GetDB()).C(uds.GetCollection()).DropCollection()
	if err != nil {
		log.Fatal("Error in dataStoreTeardown for uds: %s", err.Error())
	}

}

func tdsTeardown(tds models.TodoDataStore) {
	defer tds.Close()

	session, err := mgo.Dial(mdb.Hostname)
	if err != nil {
		log.Fatal(err)
	}
	err = session.DB(tds.GetDB()).C(tds.GetCollection()).DropCollection()
	if err != nil {
		log.Fatal("Error in dataStoreTeardown for tds: %s", err.Error())
	}

}

func TestTodosHandler0(t *testing.T) {
	req, rr, uds, _, _ := handlersSetup("GET", "api/todos")
	defer udsTeardown(uds)

	var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		TodosHandler(w, r)
	}
	ValidatePath(handler).ServeHTTP(rr, req)

	testStatus(StatusSuccess, rr, t)

	expected := "[]"
	testBody(expected, rr, t)

}

/*func TestTodosHandler1(t *testing.T) {
	req, rr, uds, tds, u := handlersSetup("GET", "api/todos")
	defer udsTeardown(uds)
	defer tdsTeardown(tds)

	t0 := models.NewTodo()
	t0.Ownerid = u.Id.Hex()
	t1 := models.NewTodo()
	t1.Ownerid = u.Id.Hex()

	tds.InsertTodo(t0)
	tds.InsertTodo(t1)

	var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		TodosHandler(w, r)
	}
	ValidatePath(handler).ServeHTTP(rr, req)

	testStatus(StatusSuccess, rr, t)

	expected := "[]"
	testBody(expected, rr, t)
}*/
