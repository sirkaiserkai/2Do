package mdb

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

const hostname = "mongodb://localhost"
const databaseName = "2DoDB"

var masterSession *mgo.Session
var NotValidObjIndexError = errors.New("Not a valid ObjectIdHex!")

func init() {
	var err error
	masterSession, err = mgo.Dial(hostname)
	if err != nil {
		panic(err)
	}
}

// TODO: Add DB as field
type DataStore struct {
	session    *mgo.Session
	Database   string
	Collection string
}

func NewDataStore() DataStore {
	d := DataStore{}
	d.session = masterSession.Copy()
	d.Database = databaseName // Assigned the default database name
	return d
}

func (d *DataStore) GetAllObjects() ([]bson.Raw, error) {
	var results []bson.Raw

	err := d.session.DB(d.Database).C(d.Collection).Find(nil).All(&results)
	if err != nil {
		return nil, err
	}

	return results, err
}

func (d *DataStore) GetObjectById(id string) (*bson.Raw, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, errors.New("id: " + id + "Is not a valid ObjectIdHex!")
	}

	var raw bson.Raw
	oid := bson.ObjectIdHex(id)

	err := d.session.DB(d.Database).C(d.Collection).FindId(oid).One(&raw)
	if err != nil {
		return nil, err
	}

	return &raw, nil
}

func (d *DataStore) InsertObject(obj interface{}) error {
	return d.session.DB(d.Database).C(d.Collection).Insert(obj)
}

func (d *DataStore) ModifyObjectForId(id string, change map[string]interface{}) error {
	if !bson.IsObjectIdHex(id) {
		return NotValidObjIndexError
	}
	oid := bson.ObjectIdHex(id)

	selector := bson.M{"_id": oid}
	err := d.session.DB(d.Database).C(d.Collection).Update(selector, change)
	if err != nil {
		log.Println("ModifyObjectForId: " + err.Error())
		return err
	}

	return nil
}
