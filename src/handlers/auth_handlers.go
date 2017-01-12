package handlers

import (
	"auth"
	"context"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"

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

	w.WriteHeader(StatusSuccess)
	w.Header().Set(ContentType, ApplicationJSON)
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
	res := jsonResponse{Result: fmt.Sprintf("%s successfully signed up for 2Do service!", u.Username)}

	v, err := json.Marshal(res)
	if err != nil {
		log.Println(fmt.Sprintf("Failure to sign up user (%s): %s", u.Username, err))
		InternalErrorHandler(w, r, "Failure to sign up: Internal Error")
		return
	}

	w.WriteHeader(StatusCreation)
	w.Header().Set(ContentType, ApplicationJSON)
	w.Write(v)
}

func ValidatePath(page http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := auth.GetTokenFromAuthHeader(r)
		if err != nil {
			BadRequestHandler(w, r, err.Error())
			return
		}

		var c auth.Claims
		token, err := jwt.ParseWithClaims(tokenString, &c, auth.KeyFunc)
		if err != nil {
			log.Printf("ValidateToken: ParseWithCalims failure: %s\n", err.Error())
			NotFoundHandler(w, r, "")
			return
		}

		// Check if user with id exists
		id := c.UserId

		uds := models.NewUserDataStore()
		log.Println(uds.GetDB())
		user, err := uds.GetUserById(id)
		if err != nil {
			log.Printf("ValidateToken: GetUserById Failed for: %v reason: %v", id, err.Error())
			NotFoundHandler(w, r, "")
			return
		}

		if user.Blocked {
			log.Printf("ValidateToken: Blocked user attempted access: %v\n", user)
			UnauthorizedHandler(w, r, "")
			return
		}

		if claims, ok := token.Claims.(*auth.Claims); ok && token.Valid {
			ctx := context.WithValue(r.Context(), auth.ClaimsKey, *claims)
			page(w, r.WithContext(ctx))
		} else {
			NotFoundHandler(w, r, "")
		}
	})
}
