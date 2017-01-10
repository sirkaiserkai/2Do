package mdb

import (
	"errors"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

const hostname = "mongodb://localhost"
const databaseName = "2DoDB"

var masterSession *mgo.Session
var NotValidObjIndexError = "Id is not a valid ObjectIdHex: %s"
var NotFoundError = mgo.ErrNotFound

func init() {
	var err error
	masterSession, err = mgo.Dial(hostname)
	if err != nil {
		panic(err)
	}
}

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

func (d *DataStore) GetObjectById(i interface{}) (*bson.Raw, error) {

	id, ok := i.(string)
	if !ok {
		return nil, errors.New("Provided id is not of type string")
	}

	if !bson.IsObjectIdHex(id) {
		return nil, fmt.Errorf(NotValidObjIndexError, id)
	}

	var raw bson.Raw
	oid := bson.ObjectIdHex(id)

	err := d.session.DB(d.Database).C(d.Collection).FindId(oid).One(&raw)
	if err != nil {
		return nil, err
	}

	return &raw, nil
}

func (d *DataStore) GetObjectForQuery(query interface{}) (*bson.Raw, error) {
	q, ok := query.(bson.M)
	if !ok {
		return nil, errors.New("Invalid query structure must be bson.M")
	}

	var raw bson.Raw
	err := d.session.DB(d.Database).C(d.Collection).Find(q).One(&raw)
	if err != nil {
		return nil, err
	}

	return &raw, nil
}

func (d *DataStore) GetObjectsForQuery(query interface{}) ([]bson.Raw, error) {
	q, ok := query.(bson.M)
	if !ok {
		return nil, errors.New("Invalid query structure must be bson.M")
	}

	var results []bson.Raw
	err := d.session.DB(d.Database).C(d.Collection).Find(q).One(&results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (d *DataStore) InsertObject(obj interface{}) error {
	return d.session.DB(d.Database).C(d.Collection).Insert(obj)
}

func (d *DataStore) ModifyObjectForId(id string, change map[string]interface{}) error {
	if !bson.IsObjectIdHex(id) {
		return fmt.Errorf(NotValidObjIndexError, id)
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

func (d *DataStore) DeleteObjectForId(id string) error {
	if !bson.IsObjectIdHex(id) {
		return fmt.Errorf(NotValidObjIndexError, id)
	}

	oid := bson.ObjectIdHex(id)

	return d.session.DB(d.Database).C(d.Collection).RemoveId(oid)
}
