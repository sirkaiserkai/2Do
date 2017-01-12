package models

import (
	"gopkg.in/mgo.v2"
	"log"
	"mdb"
	"testing"
)

// SETUP AND TEARDOWN METHODS //
func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func (uds UserDataStore) setup() {

}

func udsTeardown(uds UserDataStore) {
	defer uds.Close()

	session, err := mgo.Dial(mdb.Hostname)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	err = session.DB(uds.d.Database).C(uds.d.Collection).DropCollection()
	if err != nil {
		log.Fatal("Error in udsTeardown: %s", err.Error())
	}

}

func TestInsertUser(t *testing.T) {
	// Setup test
	uds := NewUserDataStore()
	uds.d.Collection = "2Do_TestInsertUser_Collection"
	uds.setup()

	defer udsTeardown(uds)

	// Main test content
	u := NewUser()
	err := uds.InsertUser(u)
	if err != nil {
		t.Error(err)
	}
}

func TestGetUserById(t *testing.T) {
	// Setup test
	uds := NewUserDataStore()
	uds.d.Collection = "2Do_TestGetUserById_Collection"
	uds.setup()

	defer udsTeardown(uds)

	u := NewUser()
	err := uds.InsertUser(u)
	if err != nil {
		t.Fatal(err)
	}

	// Main test content
	usr, err := uds.GetUserById(u.Id.Hex())
	if err != nil {
		t.Error(err)
	}

	if usr.Id.Hex() != u.Id.Hex() {
		t.Error("User ids do not match")
	}
}

func TestGetUserByName(t *testing.T) {
	// Setup test
	uds := NewUserDataStore()
	uds.d.Collection = "2Do_TestGetUserByName_Collection"
	uds.setup()

	defer udsTeardown(uds)

	u := NewUser()
	u.Username = "Some Dude"
	err := uds.InsertUser(u)
	if err != nil {
		t.Fatal(err)
	}

	// Main test content
	usr, err := uds.GetUserByName(u.Username)
	if err != nil {
		t.Error(err)
	}

	if usr.Username != u.Username {
		t.Error("Usernames do not match")
	}

}

func TestModifyUser(t *testing.T) {
	// Setup test
	uds := NewUserDataStore()
	uds.d.Collection = "2Do_TestModifyUser_Collection"
	uds.setup()

	defer udsTeardown(uds)
	u := NewUser()
	u.Username = "Some Dude"
	err := uds.InsertUser(u)
	if err != nil {
		t.Fatal(err)
	}

	// Main test content
	change := make(map[string]interface{})
	change["username"] = "DiffName"
	err = uds.ModifyUser(u.Id.Hex(), change)
	if err != nil {
		t.Error(err)
	}

	usr, err := uds.GetUserById(u.Id.Hex())
	if err != nil {
		t.Fatal(err)
	}

	if usr.Username != "DiffName" {
		t.Error("User's name is incorrect")
	}

}

func TestDeleteUser(t *testing.T) {
	// Setup test
	uds := NewUserDataStore()
	uds.d.Collection = "2Do_TestDeleteUser_Collection"
	uds.setup()

	defer udsTeardown(uds)

	u := NewUser()
	err := uds.InsertUser(u)
	if err != nil {
		t.Fatal(err)
	}

	// Main test content
	err = uds.DeleteUser(u.Id.Hex())
	if err != nil {
		t.Error(err)
	}

	usr, err := uds.GetUserById(u.Id.Hex())
	if err == nil || usr != nil {
		t.Error("Did not fail in getting the user")
	}

}
