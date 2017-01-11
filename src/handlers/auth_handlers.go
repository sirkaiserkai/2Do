package handlers

import (
	"auth"
	"encoding/json"
	"fmt"
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
		BadRequestHandler(w, r, "Failure to log in")
		log.Println("Failure to Log In: %s\n", err)
		return
	}

	if !auth.PasswordEquality(u.Password, password) {
		BadRequestHandler(w, r, "Failure to log in")
		log.Println("Failure to Log In: Passwords not equal")
		return
	}

	token, err := auth.CreateToken(*u)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(200)
	w.Write([]byte(
		"{'result': 'Successfully logged into 2Do'," +
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
		BadRequestHandler(w, r, "Password must be above 6 characters in length")
		return
	}

	uds := models.NewUserDataStore()

	sameUsernameUser, err := uds.GetUserByName(username)
	if err != nil {
		if err != models.ErrUserNotFound {
			log.Printf("SignUpHandler: GetUserByName: Failure to get user: %s\n", err)
			InternalErrorHandler(w, r, "Failure to sign up: Internal Error")
			return
		}
	}

	if sameUsernameUser != nil {
		log.Printf("SignUpHandler: User with user already exists for username: %s\n", err.Error())
		BadRequestHandler(w, r, "User already exists with username")
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
		InternalErrorHandler(w, r, "Failure to sign up: Internal Error")
	}

	log.Printf("User successfully created: %v\n", u)
	res := jsonResponse{result: fmt.Sprintf("%s successfully signed up for 2Do service!", u.Username)}

	v, err := json.Marshal(res)
	if err != nil {
		log.Println(fmt.Sprintf("Failure to sign up user (%s): %s", u.Username, err))
		InternalErrorHandler(w, r, "Failure to sign up: Internal Error")
		return
	}

	w.WriteHeader(201)
	w.Write(v)
}
