package models

import (
	"errors"
	"gopkg.in/mgo.v2/bson"
	"mdb"
	"time"
)

const TodoCollection = "todos"

var TodoConvertError = errors.New("Failed to convert to Todo type")

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

type TodoDataStore struct {
	d mdb.DataStore
}

func NewTodoDataStore() TodoDataStore {
	tds := TodoDataStore{}
	tds.d = mdb.NewDataStore()
	tds.d.Collection = TodoCollection
	return tds
}

func (tds TodoDataStore) Close() {
	tds.d.Close()
}

func (tds *TodoDataStore) SetDB(db string) {
	tds.d.Database = db
}

func (tds TodoDataStore) GetDB() string {
	return tds.d.Database
}

func (tds *TodoDataStore) SetCollection(coll string) {
	tds.d.Collection = coll
}

func (tds TodoDataStore) GetCollection() string {
	return tds.d.Collection
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

// TODO: Remove ability to change ownerid field
func (tds TodoDataStore) ModifyTodo(todoId, userId string, changes map[string]interface{}) error {
	params := make(map[string]string)
	params["id"] = todoId
	params["ownerid"] = userId
	return tds.d.ModifyObjectForId(params, changes)
}

func (tds TodoDataStore) DeleteTodo(id, userId string) error {
	m := make(map[string]string)
	m["id"] = id
	m["ownerid"] = userId

	return tds.d.DeleteObjectForSelector(m)
}
