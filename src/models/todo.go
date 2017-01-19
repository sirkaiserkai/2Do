package models

import (
	"errors"
	"gopkg.in/mgo.v2/bson"
	"mdb"
	"time"
)

const TodoCollection = "todos"

var TodoConvertError = errors.New("Failed to convert to Todo type")
var TodoNotFoundError = mdb.NotFoundError

type Todo struct {
	Id bson.ObjectId `json:"id" bson:"_id,omitempty"` // MongoDBId
	// TODO: Add id which would be exposed and hide the mongodb id
	// see http://stackoverflow.com/a/13740114/2812587 for more info
	Title     string    `json:"title" bson:"title"`
	Note      string    `json:"note" bson:"note"`
	Created   time.Time `json:"created_date" bson:"created_date,omitempty"`
	Due       time.Time `json:"due_date" bson:"due_date,omitempty"`
	Ownerid   string    `json:"-" bson:"ownerid"`
	Completed bool      `json:"completed" bson:"completed"`
}

func NewTodo() Todo {
	t := Todo{}
	t.Id = bson.NewObjectId()
	return t
}

// TodoStorage is an interface which details the requirments
// to interface with retrieval and insertion of todos
// into long term storage.
type TodoStorage interface {
	Close()
	GetAllTodos() ([]Todo, error)
	GetTodoById(id string) (*Todo, error)
	GetTodosForUserId(id string) ([]Todo, error)
	InsertTodo(t Todo) error
	ModifyTodo(todoId, userId string, changes map[string]interface{}) error
	DeleteTodo(id, userId string) error
}

// NewTodoStorage is the abstracted function that returns
// a TodoStorage implementation depending on the value of
// the TODO_STORE_TYPE.
func NewTodoStorage() TodoStorage {
	switch TODO_STORE_TYPE {
	case Regular:
		return NewTodoDataStore()
	case Test:
		return newTestTodoStorage()
	}

	return NewTodoDataStore()
}

// TodoDataStore is a wrapper struct for DataStore.
// It implements the TodoStorage interface
type TodoDataStore struct {
	d mdb.DataStore
}

func NewTodoDataStore() *TodoDataStore {
	tds := TodoDataStore{}
	tds.d = mdb.NewDataStore()
	tds.d.Collection = TodoCollection
	return &tds
}

func (tds *TodoDataStore) Close() {
	tds.d.Close()
}

func (tds *TodoDataStore) SetDB(db string) {
	tds.d.Database = db
}

func (tds *TodoDataStore) GetDB() string {
	return tds.d.Database
}

func (tds *TodoDataStore) SetCollection(coll string) {
	tds.d.Collection = coll
}

func (tds *TodoDataStore) GetCollection() string {
	return tds.d.Collection
}

func (tds *TodoDataStore) GetAllTodos() ([]Todo, error) {

	ts := make([]Todo, 0)

	raws, err := tds.d.GetAllObjects()
	if err != nil {
		return nil, err
	}

	for _, raw := range raws {
		t := Todo{}
		err := raw.Unmarshal(&t)
		if err != nil {
			return nil, err
		}
		ts = append(ts, t)
	}

	return ts, nil
}

func (tds *TodoDataStore) GetTodoById(id string) (*Todo, error) {
	t := Todo{}

	raw, err := tds.d.GetObjectById(id)
	if err != nil {
		return nil, err
	}

	err = raw.Unmarshal(&t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (tds *TodoDataStore) GetTodosForUserId(id string) ([]Todo, error) {
	ts := make([]Todo, 0)
	query := bson.M{"ownerid": id}

	raws, err := tds.d.GetObjectsForQuery(query)
	if err != nil {
		return nil, err
	}

	for _, raw := range raws {
		t := Todo{}
		err := raw.Unmarshal(&t)
		if err != nil {
			return nil, err
		}
		ts = append(ts, t)
	}

	return ts, nil
}

func (tds *TodoDataStore) InsertTodo(t Todo) error {
	return tds.d.InsertObject(t)
}

func (tds *TodoDataStore) ModifyTodo(todoId, userId string, changes map[string]interface{}) error {
	params := make(map[string]string)
	params["id"] = todoId
	params["ownerid"] = userId

	// This is required because we're using the $set operator to replace values
	// of a specified field. It will create the field in the db lest we
	// remove it explicitly from the changes map. (TODO: Find a better way of
	// doing this)
	// https://docs.mongodb.com/manual/reference/operator/update/set/
	allowedKeys := []string{"title", "note", "due_date", "created_date"}
	for k, v := range changes {
		found := false

		for _, modifiableVal := range allowedKeys {
			if k == modifiableVal {
				found = true
				break
			}
		}

		if !found {
			delete(changes, k)
		}

		if _, ok := v.(string); !ok {
			return errors.New("Incorrect format for key: " + k)
		}
	}

	err := tds.d.ModifyObjectForId(params, changes)
	if err != nil {
		if err == mdb.NotFoundError {
			return TodoNotFoundError
		} else {
			return err
		}
	}

	return nil
}

func (tds *TodoDataStore) DeleteTodo(id, userId string) error {
	m := make(map[string]string)
	m["id"] = id
	m["ownerid"] = userId
	return tds.d.DeleteObjectForSelector(m)
}
