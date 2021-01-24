package controllers

import (
	"JWT-go/models"
	"JWT-go/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type Controller struct {
}

func (c Controller) Signup(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		var error models.Error
		json.NewDecoder(r.Body).Decode(&user)

		//Input validation for email and password
		if user.Email == "" {
			error.Message = "email is missing"
			utils.RespondWithError(w, http.StatusBadRequest, error)
			return
		}
		if user.Password == "" {
			error.Message = "password is missing"
			utils.RespondWithError(w, http.StatusBadRequest, error)
			return
		}
		//hash the password
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
		if err != nil {
			log.Print(err)
		}
		user.Password = string(hash)

		stmt := "insert into users (email,password) values ($1,$2) RETURNING id"
		err = db.QueryRow(stmt, user.Email, user.Password).Scan(&user.ID)

		if err != nil {
			log.Print(err)
			error.Message = "Server error"
			utils.RespondWithError(w, http.StatusInternalServerError, error)
			return
		}

		user.Password = "" //done for not returing password to user in the response
		w.Header().Set("Content-Type", "application/json")
		utils.ResponseJSON(w, user)
	}
}

func (c Controller) Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		var error models.Error
		var jwt models.JWT

		json.NewDecoder(r.Body).Decode(&user)

		//Input validation for email and password
		if user.Email == "" {
			error.Message = "email is missing"
			utils.RespondWithError(w, http.StatusBadRequest, error)
			return
		}
		if user.Password == "" {
			error.Message = "password is missing"
			utils.RespondWithError(w, http.StatusBadRequest, error)
			return
		}
		//password received from the user post request
		password := user.Password

		//check if user exists inside our DB
		err := db.QueryRow("select * from users where email = $1", user.Email).Scan(&user.ID, &user.Email, &user.Password)
		if err != nil {
			if err == sql.ErrNoRows {
				error.Message = "The user does not exist"
				utils.RespondWithError(w, http.StatusBadRequest, error)
				return
			} else {
				log.Fatal(err)
			}
		}
		hashedPassword := user.Password
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			error.Message = "Invalid Password"
			utils.RespondWithError(w, http.StatusUnauthorized, error)
			return
		}
		token, err := utils.GenerateToken(user)
		if err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusOK)
		jwt.Token = token
		utils.ResponseJSON(w, jwt)
	}
}

func (c Controller) TokenVerifyMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var errorObject models.Error
		//return the token containes in thte user request
		authHeader := r.Header.Get("Authorization")
		//log.Print(authHeader)
		bearerToken := strings.Split(authHeader, " ")

		log.Print(bearerToken)
		if len(bearerToken) == 2 {
			authToken := bearerToken[1]

			//parse jwt token and verify algorithm
			token, error := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*(jwt.SigningMethodHMAC)); !ok {
					return nil, fmt.Errorf("there was an error")
				}
				return []byte(os.Getenv("SECRET")), nil
			})
			if error != nil {
				errorObject.Message = error.Error()
				utils.RespondWithError(w, http.StatusUnauthorized, errorObject)
				return
			}
			if token.Valid {
				next.ServeHTTP(w, r)
			} else {
				errorObject.Message = error.Error()
				utils.RespondWithError(w, http.StatusUnauthorized, errorObject)
				return
			}
		} else {
			errorObject.Message = "Invalid token"
			utils.RespondWithError(w, http.StatusUnauthorized, errorObject)
			return
		}
	})
}
