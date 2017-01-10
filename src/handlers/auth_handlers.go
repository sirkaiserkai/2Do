package handlers

import (
	"auth"
	"encoding/json"
	"log"
	"models"
	"net/http"
)

func LogInHandler(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		panic(err)
	}

	username, ok := m["username"].(string)
	if !ok {
		return
	}
	password, ok := m["password"].(string)
	if !ok {
		return
	}

	uds := models.NewUserDataStore()

	u, err := uds.GetUserByName(username)
	if err != nil {
		log.Println("Failure to Log In: %s\n", err)
		w.WriteHeader(400)
		w.Write([]byte("Failure to Log in"))
		return
	}

	if !auth.PasswordEquality(u.Password, password) {
		log.Println("Failure to Log In: Passwords not equal")
		w.WriteHeader(400)
		w.Write([]byte("Failure to Log in"))
		return
	}
	log.Println(u)
	token, err := auth.CreateToken(*u)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(200)
	w.Write([]byte(
		"{'message': 'Successfully logged into 2Do'," +
			" 'token': '" + token + "'}"))
}

func isPasswordLegit(password string) bool {
	if len(password) < 6 {
		return false
	}

	/*if password == "password" ||
		password == "123456" {
		return false
	}*/

	return true
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		panic(err)
	}

	username, ok := m["username"].(string)
	if !ok {
		return
	}
	password, ok := m["password"].(string)
	if !ok {
		return
	}

	if !isPasswordLegit(password) {
		w.WriteHeader(400)
		w.Write([]byte("Password must be above 6 characters in legnth"))
	}

	uds := models.NewUserDataStore()

	sameUsernameUser, err := uds.GetUserByName(username)
	if err != nil {
		if err != models.ErrUserNotFound {
			log.Printf("SignUpHandler: GetUserByName: Failure to get user: %s\n", err)
			w.WriteHeader(400)
			w.Write([]byte("Error: Failure to sign up"))
			return
		}
	}

	if sameUsernameUser != nil {
		log.Printf("SignUpHandler: User with username already exists\n")
		w.WriteHeader(400)
		w.Write([]byte("Failure to sign up: User with username already exists"))
		return
	}

	u := models.NewUser()
	u.Username = username
	u.Password, err = auth.HashPassword(password)
	if err != nil {
		panic(err)
	}

	err = uds.InsertUser(u)
	if err != nil {
		log.Printf("SignUpHadler: InsertUser Error: %s\n", err)
		w.WriteHeader(400)
		w.Write([]byte("Error: Failure to sign up"))
	}

	log.Printf("User successfully created: %v\n", u)
	w.WriteHeader(200)
	w.Write([]byte("Successfully signed up for 2Do service!"))
}
