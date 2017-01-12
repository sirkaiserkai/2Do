package mdb

import (
	"errors"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

const Hostname = "mongodb://localhost"

var DatabaseName = "2DoDB"
var masterSession *mgo.Session

var NotFoundError = mgo.ErrNotFound

type NotValidObjIndexError error

func notValidObjIndexError(id string) NotValidObjIndexError {
	return fmt.Errorf("Id is not a valid ObjectIdHex: %s", id)
}

func init() {
	var err error
	masterSession, err = mgo.Dial(Hostname)
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
	d.Database = DatabaseName // Assigned the default database name
	return d
}

func (d DataStore) Close() {
	d.session.Close()
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
		return nil, notValidObjIndexError(id)
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
	err := d.session.DB(d.Database).C(d.Collection).Find(q).All(&results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (d *DataStore) InsertObject(obj interface{}) error {
	return d.session.DB(d.Database).C(d.Collection).Insert(obj)
}

func (d *DataStore) ModifyObjectForId(params map[string]string, change map[string]interface{}) error {

	selector := bson.M{}
	for k, v := range params {
		if k == "id" {
			id := params["id"]
			if !bson.IsObjectIdHex(id) {
				return notValidObjIndexError(id)
			}
			oid := bson.ObjectIdHex(id)
			selector["_id"] = oid
		} else {
			selector[k] = v
		}
	}

	err := d.session.DB(d.Database).C(d.Collection).Update(selector, change)
	if err != nil {
		log.Println("ModifyObjectForId: " + err.Error())
		return err
	}

	return nil
}

func (d *DataStore) DeleteObjectForSelector(params map[string]string) error {

	selector := bson.M{}
	for k, v := range params {
		if k == "id" {
			id := params["id"]
			if !bson.IsObjectIdHex(id) {
				return notValidObjIndexError(id)
			}
			oid := bson.ObjectIdHex(id)
			selector["_id"] = oid
		} else {
			selector[k] = v
		}
	}

	return d.session.DB(d.Database).C(d.Collection).Remove(selector)
}
