package mdb

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"gopkg.in/mgo.v2/bson"
)

const testDB = "2DoDB"

var (
	ts0         = TestStruct{Id: bson.NewObjectId(), Val0: "Word to your mother"}
	ts1         = TestStruct{Id: bson.NewObjectId(), Val0: "Hello, world!", Val1: 1234}
	ts2         = TestStruct{Id: bson.NewObjectId(), Val0: "Testing Forever", Val1: 9876}
	testStructs = []TestStruct{ts0, ts1, ts2}
)

// SETUP AND TEARDOWN METHODS //

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile) // Flags make log statements print code line number
}

// getSetup reassigns the database to be the test database
func (d DataStore) setup() {
	d.Database = testDB
}

// getSetup is a test setup method that inserts the three TestStructs
// into the database. Note: This is dependent on the InsertObject func
// and calling setup is redundant
func (d DataStore) getSetup() {
	d.setup()
	for _, ts := range testStructs {
		if err := d.InsertObject(ts); err != nil {
			log.Fatal(err)
		}
	}
}

// teardown deletes the collection of the datastore
func teardown(d DataStore) {
	// Deletes the collection
	defer d.session.Close()
	err := d.session.DB(d.Database).C(d.Collection).DropCollection()
	if err != nil {
		log.Fatal("Error in teardown: %s", err.Error())
	}
}

// END OF SETUP AND TEST METHODS //

// TESTSTRUCT INTERFACE AND INTERFACE METHODS //
type TestStruct struct {
	Id   bson.ObjectId `bson:"_id,omitempty"`
	Val0 string        `bson:"value0"`
	Val1 int           `bson:"value1"`
}

func (t0 TestStruct) Equal(t1 TestStruct) bool {
	if t0.Id != t1.Id {
		return false
	}

	if t0.Val0 != t1.Val0 {
		return false
	}

	if t0.Val1 != t1.Val1 {
		return false
	}

	return true
}

func (t TestStruct) String() string {
	return fmt.Sprintf("TestStruct: %v\tVal0: %s\tVal1: %d", t.Id, t.Val0, t.Val1)
}

// Shitty runtime but it's just for testing so fuck it.
func setComparison(s0, s1 []TestStruct) bool {
	if len(s0) != len(s1) {
		return false
	}

	for _, t0 := range s0 {
		found := false
		for _, t1 := range s1 {
			if t0.Equal(t1) {
				found = true
			}
		}

		if !found {
			return false
		}
	}

	return true
}

// END OF TESTSTRUCT INTERFACE AND INTERFACE METHODS //

// TEST FUNCTIONS //

func TestInsertObject(t *testing.T) {
	// Test Setup
	d := NewDataStore()
	d.Collection = "2Do_TestInsertObject_Collection"
	d.setup()
	defer teardown(d)

	// Main test content
	if err := d.InsertObject(ts0); err != nil {
		t.Error(err)
	}

	if err := d.InsertObject(ts1); err != nil {
		t.Error(err)
	}

	if err := d.InsertObject(ts2); err != nil {
		t.Error(err)
	}
}

func TestGetObjects(t *testing.T) {
	// Test Setup
	d := NewDataStore()
	d.Collection = "2Do_TestGetObjects_Collection"
	d.getSetup()

	defer teardown(d)

	// Main test content
	objs, err := d.GetAllObjects()
	if err != nil {
		t.Error(err)
	}

	ts := make([]TestStruct, 0)
	for _, o := range objs {
		testStruct := TestStruct{}
		err := o.Unmarshal(&testStruct)
		if err != nil {
			t.Error(err)
		}

		ts = append(ts, testStruct)

	}

	if !setComparison(ts, testStructs) {
		t.Error("Test structs not equal")
	}
}

func TestGetObjectById(t *testing.T) {
	// Test setup
	d := NewDataStore()
	d.Collection = "2Do_TestGetObjectById_Collection"
	d.getSetup()
	defer teardown(d)

	// Main test content
	var id = ts1.Id.Hex()
	obj, err := d.GetObjectById(id)
	if err != nil {
		t.Error(err)
	}

	var tstStrct TestStruct

	if err := obj.Unmarshal(&tstStrct); err != nil {
		t.Error(err)
	}

}

func TestModifyObjectForId(t *testing.T) {
	// Test setup
	d := NewDataStore()
	d.Collection = "2Do_TestGetObjectById_Collection"
	d.getSetup()

	defer teardown(d)

	// Main test content
	id := ts0.Id.Hex()
	change := bson.M{"value0": "Updated Value"}
	err := d.ModifyObjectForId(id, change)
	if err != nil {
		t.Error(err)
	}

	obj, err := d.GetObjectById(id)
	if err != nil {
		t.Fatal(err)
	}

	var tstStrct TestStruct

	if err := obj.Unmarshal(&tstStrct); err != nil {
		t.Fatal(err)
	}

	if strings.Compare(tstStrct.Val0, "Updated Value") != 0 {
		t.Error("Val0 != \"Updated Value\"")
		t.Log("Val0: " + tstStrct.Val0)
	}

	err = d.ModifyObjectForId("123", change)
	if err == nil {
		t.Error("Did not return NotValidObjIndexError")
	}
}

func TestDeleteObjectForId(t *testing.T) {
	// Test setup
	d := NewDataStore()
	d.Collection = "2Do_TestGetObjectById_Collection"
	d.getSetup()

	defer teardown(d)

	// Main test content
	id := ts0.Id.Hex()
	err := d.DeleteObjectForId(id)
	if err != nil {
		t.Error(err)
	}

	id = "abc"
	err = d.DeleteObjectForId(id)
	if err == nil {
		t.Error(err)
	}
}

// END OF TEST FUNCTIONS //
