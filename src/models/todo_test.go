package models

import (
	"testing"

	"gopkg.in/mgo.v2"

	"log"
	"mdb"
)

const testDB = "2DoDB"

// SETUP AND TEARDOWN METHODS //

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func (tds TodoDataStore) setup() {
	tds.d.Database = testDB
}

func tdsTeardown(tds TodoDataStore) {
	defer tds.Close()

	session, err := mgo.Dial(mdb.Hostname)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	err = session.DB(tds.d.Database).C(tds.d.Collection).DropCollection()
	if err != nil {
		log.Fatal("Error in tdsTeardown: %s", err.Error())
	}
}

func setComparison(ts0, ts1 []Todo) bool {
	if len(ts0) != len(ts1) {
		return false
	}

	for _, t0 := range ts0 {
		found := false
		for _, t1 := range ts1 {
			if t0 == t1 {
				found = true
			}
		}

		if !found {
			return false
		}
	}

	return true
}

// END OF SETUP AND TEARDOWN METHODS //

func TestInsertTodo(t *testing.T) {
	// Test setup
	tds := NewTodoDataStore()
	tds.d.Collection = "2Do_TestInsert2Do_Collection"
	tds.setup()

	defer tdsTeardown(tds)

	// Main test content
	t0 := Todo{}
	err := tds.InsertTodo(t0)
	if err != nil {
		t.Error(err)
	}
}

func TestGetAllTodos(t *testing.T) {
	// Test setup
	tds := NewTodoDataStore()
	tds.d.Collection = "2Do_TestGetAllTodos_Collection"
	tds.setup()

	defer tdsTeardown(tds)

	t0 := NewTodo()
	t0.Title = "Example 0"
	t1 := NewTodo()
	t1.Title = "Example 1"
	t1.Note = "Note Example 1"
	t2 := NewTodo()
	t2.Title = "Example 2"
	t2.Note = "Note Example 2"
	t2.Ownerid = "12345"

	tds.InsertTodo(t0)
	tds.InsertTodo(t1)
	tds.InsertTodo(t2)

	// Main test content
	ts, err := tds.GetAllTodos()
	if err != nil {
		t.Error(err)
	}

	if len(ts) != 3 {
		t.Error("Not correct number of Todos")
	}

}

func TestGetTodoById(t *testing.T) {
	// Test setup
	tds := NewTodoDataStore()
	tds.d.Collection = "2Do_TestGetTodoById_Collection"
	tds.setup()

	defer tdsTeardown(tds)

	t0 := NewTodo()
	t0.Title = "Example 0"

	tds.InsertTodo(t0)

	// Main test content
	t_0, err := tds.GetTodoById(t0.Id.Hex())
	if err != nil {
		t.Fatal(err)
	}

	if *t_0 != t0 {
		t.Error("t_0 not equal to t0")
	}
}

func TestGetTodosForUserId(t *testing.T) {
	// Test setup
	tds := NewTodoDataStore()
	tds.d.Collection = "2Do_TestGetTodosForUserId_Collection"
	tds.setup()

	defer tdsTeardown(tds)

	ownerId := "123456"
	t0 := NewTodo()
	t0.Ownerid = ownerId

	t1 := NewTodo()
	t1.Ownerid = ownerId

	tds.InsertTodo(t0)
	tds.InsertTodo(t1)

	ts0 := []Todo{t0, t1}
	// Main test content
	ts1, err := tds.GetTodosForUserId(ownerId)
	if err != nil {
		t.Fatal(err)
	}
	if !setComparison(ts0, ts1) {
		t.Error("Sets not equal")
	}

	ts1, err = tds.GetTodosForUserId("abcde")
	if err != nil {
		t.Error(err)
	}

	if len(ts1) != 0 {
		t.Error("Should not find any todos for user id")
	}
}

func TestModifyTodo(t *testing.T) {
	// Test setup
	tds := NewTodoDataStore()
	tds.d.Collection = "2Do_TestModifyTodo_Collection"
	tds.setup()

	defer tdsTeardown(tds)

	ownerId := "12345"
	t0 := NewTodo()
	t0.Title = "Hello, world!"
	t0.Ownerid = ownerId

	tds.InsertTodo(t0)
	// Main test content
	changes := make(map[string]interface{})
	changes["title"] = "Changed Title"
	changes["note"] = "Example Note"
	err := tds.ModifyTodo(t0.Id.Hex(), ownerId, changes)
	if err != nil {
		t.Error(err)
	}

	err = tds.ModifyTodo("1234", ownerId, changes)
	if err == nil {
		t.Error("Error: Should be not found error")
	}
}

func TestDeleteTodo(t *testing.T) {
	// Test setup
	tds := NewTodoDataStore()
	tds.d.Collection = "2Do_TestDeleteTodo_Collection"
	tds.setup()

	defer tdsTeardown(tds)

	userId := "12345"
	t0 := NewTodo()
	t0.Ownerid = userId
	tds.InsertTodo(t0)

	// Main test content
	err := tds.DeleteTodo(t0.Id.Hex(), userId)
	if err != nil {
		t.Error(err)
	}

}
