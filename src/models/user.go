package models

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"mdb"
)

const UserCollection = "users"

var ErrUserNotFound = mdb.NotFoundError

type User struct {
	Id       bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Username string        `json:"username" bson:"username"`
	Password string        `json:"-" bson:"password"` // hashed password
	Blocked  bool          `json:"-" bson:"blocked"`
}

func (u User) String() string {
	return fmt.Sprintf("User: [ '%s', '%s' ]", u.Id.Hex(), u.Username)
}

func NewUser() User {
	u := User{}
	u.Id = bson.NewObjectId()
	return u
}

type UserDataStore struct {
	d mdb.DataStore
}

func NewUserDataStore() UserDataStore {
	uds := UserDataStore{}
	uds.d = mdb.NewDataStore()
	uds.d.Collection = UserCollection
	return uds
}

func (uds *UserDataStore) GetUserById(id string) (*User, error) {
	/*u := User{}

	raw, err := uds.d.GetObjectById(id)
	if err != nil {
		return nil, err
	}

	err = raw.Unmarshal(&u)
	if err != nil {
		return nil, err
	}

	return &u, nil*/
	return getUser(id, uds.d.GetObjectById)
}

func (uds *UserDataStore) GetUserByName(name string) (*User, error) {
	q := bson.M{"username": name}
	return getUser(q, uds.d.GetObjectForQuery)
}

// getUser is a wrapper method that handles converting the bson raw result
// to a User type.
func getUser(param interface{}, queryFunc func(interface{}) (*bson.Raw, error)) (*User, error) {
	u := User{}

	raw, err := queryFunc(param)
	if err != nil {
		return nil, err
	}

	err = raw.Unmarshal(&u)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (uds *UserDataStore) InsertUser(u User) error {
	return uds.d.InsertObject(u)
}

func (uds *UserDataStore) ModifyUser(id string, change map[string]interface{}) error {
	return uds.d.ModifyObjectForId(id, change)
}

func (uds *UserDataStore) DeleteUser(id string, change map[string]interface{}) error {
	return uds.d.DeleteObjectForId(id)
}
