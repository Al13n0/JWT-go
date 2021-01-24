package utils

import (
	"JWT-go/models"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
)

//ResponseJSOsend json respone
func ResponseJSON(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
}

//RespondWithError responde with error json message to user
func RespondWithError(w http.ResponseWriter, status int, error models.Error) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(error)
}

//GenerateToken is the function to generate a JWT token in take in input a user struct and return a token and an error
func GenerateToken(user models.User) (string, error) {
	//var err Error
	var secret = os.Getenv("SECRET")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"iss":   "testiss",
	})
	//sign the token
	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		log.Fatal(err)

	}
	return tokenString, nil
}
