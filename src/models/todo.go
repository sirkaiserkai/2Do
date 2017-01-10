package models

import (
	"errors"
	"gopkg.in/mgo.v2/bson"
	"mdb"
	"time"
	//"log"
)

const TodoCollection = "todos"

var TodoConvertError = errors.New("Failed to convert to Todo type")

type Todo struct {
	Id bson.ObjectId `json:"id" bson:"_id,omitempty"` // MongoDBId
	// TODO: Add id which would be exposed and hide the mongodb id
	// see http://stackoverflow.com/a/13740114/2812587 for more info
	Title   string    `json:"title" bson:"title"`
	Note    string    `json:"note" bson:"note"`
	Created time.Time `json:"created_date" bson:"created_date"`
	Due     time.Time `json:"due_date" bson:"due_date"`
	Ownerid string    `json:"-" bson:"ownerid"`
}

func NewTodo() Todo {
	t := Todo{}
	t.Id = bson.NewObjectId()
	return t
}

type TodoDataStore struct {
	d mdb.DataStore
}

func NewTodoDataStore() TodoDataStore {
	tds := TodoDataStore{}
	tds.d = mdb.NewDataStore()
	tds.d.Collection = TodoCollection
	return tds
}

func (tds TodoDataStore) GetAllTodos() ([]Todo, error) {

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

func (tds TodoDataStore) GetTodoById(id string) (*Todo, error) {
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

func (tds TodoDataStore) GetTodosForUserId(id string) ([]Todo, error) {
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

func (tds TodoDataStore) InsertTodo(t Todo) error {
	return tds.d.InsertObject(t)
}

func (tds TodoDataStore) ModifyTodo(id string, change map[string]interface{}) error {
	return tds.d.ModifyObjectForId(id, change)
}

func (tds TodoDataStore) DeleteTodo(id string) error {
	return tds.d.DeleteObjectForId(id)
}
