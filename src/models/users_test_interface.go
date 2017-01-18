package models

import (
	"log"
)

var usersMap = make(map[string]User)

// the UserStorage interface
type TestUserStorage struct {
	users *map[string]User
}

func newTestUserStorage() *TestUserStorage {
	t := TestUserStorage{}
	t.users = &usersMap
	return &t
}
func (tus *TestUserStorage) Close() {
	log.Println("Closing TestUserStorage")
}

func (tus *TestUserStorage) GetUserById(id string) (*User, error) {
	users := *tus.users
	if usr, ok := users[id]; ok {
		return &usr, nil
	}

	return nil, ErrUserNotFound
}
func (tus *TestUserStorage) GetUserByName(name string) (*User, error) {
	if tus.users == nil {
		log.Fatal("tus.users nil: map must be made first.")
	}
	users := *tus.users

	for _, usr := range users {
		if usr.Username == name {
			return &usr, nil
		}
	}

	return nil, ErrUserNotFound

}
func (tus *TestUserStorage) InsertUser(u User) error {
	(*tus.users)[u.Id.Hex()] = u
	return nil
}
func (tus *TestUserStorage) ModifyUser(id string, change map[string]interface{}) error {
	return nil
}
func (tus *TestUserStorage) DeleteUser(id string) error {
	users := (*tus.users)
	_, ok := users[id]
	if !ok {
		return ErrUserNotFound
	}

	delete(users, id)
	return nil
}
