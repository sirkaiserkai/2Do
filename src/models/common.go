package models

var TODO_STORE_TYPE StoreType = Regular
var USER_STORE_TYPE StoreType = Regular

// Used to set the the store type for testing purposes.
type StoreType int

const (
	Regular StoreType = 0
	Test    StoreType = 1
)
