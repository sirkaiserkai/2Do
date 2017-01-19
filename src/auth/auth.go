package auth

import (
	"config"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"log"
	"models"
	"net/http"
	"strings"
	"time"
)

var secret = config.GetConfig().Secret

const ClaimsKey = 0

type Claims struct {
	UserId string `json:"userId"`
	//Username string `json:"username"`
	jwt.StandardClaims
}

func CreateToken(u models.User) (string, error) {
	expireToken := time.Now().Add(time.Hour * 48).Unix()

	claims := Claims{
		u.Id.Hex(),
		//u.Username,
		jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    config.GetConfig().JWTIssuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		panic(err) // TODO: Handle error more eloquently
	}

	return signedToken, nil
}

// Ripped from: https://github.com/auth0/go-jwt-middleware/blob/f3f7de3b9e394e3af3b88e1b9457f6f71d1ae0ac/jwtmiddleware.go
func GetTokenFromAuthHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", nil // No error, just no token
	}

	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 ||
		strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", fmt.Errorf("Authorization header format must be Bearer {token}")
	}

	return authHeaderParts[1], nil
}

// keyFunc is the callback method used by the jwt parser to provide the
// correct key to parse unverified tokens.
// (See: https://godoc.org/github.com/dgrijalva/jwt-go#Keyfunc )
func KeyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected Signing method")
	}

	return []byte(secret), nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func PasswordEquality(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		log.Printf("PasswordEquality: Failure to login err: %s\n", err.Error())
		return false
	}

	return true
}

/*
func ForbiddenReqHandler(message string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf("Forbidden Request: %s", message)))
		return
	}
}*/

/*
func UnauthorizedReqHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(401)
	w.Write([]byte("Unauthorized Access"))
}*/

func GetClaims(r *http.Request) (*Claims, error) {
	claims, ok := r.Context().Value(ClaimsKey).(Claims)
	if !ok {
		err := fmt.Errorf("Failed to retrieve Claims")
		log.Printf("GetClaims: %s\n", err.Error())
		return nil, err
	}

	return &claims, nil
}
