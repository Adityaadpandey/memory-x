package users

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/adityaadpandey/memory-x/go-api/internal/dbclient"
	"github.com/adityaadpandey/memory-x/go-api/internal/types"
	"github.com/adityaadpandey/memory-x/go-api/internal/utils/jwttoken"
	response "github.com/adityaadpandey/memory-x/go-api/internal/utils/response"
	"github.com/adityaadpandey/memory-x/go-api/prisma/db"
)

func Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user types.User

		// Decode incoming JSON
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.Response{
				Status: response.Error,
				Error:  "Invalid request payload",
			})
			return
		}
		//  hashing the pass
		encyPass, err := jwttoken.HashPassword(user.Password)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.Response{
				Status: response.Error,
				Error:  "Failed to hash password",
			})
			return
		}
		// Create the user in the DB and capture the result
		createdUser, err := dbclient.PrismaClient.User.CreateOne(
			db.User.Name.Set(user.Name),
			db.User.Email.Set(user.Email),
			db.User.Password.Set(encyPass),
		).Exec(r.Context())
		if err != nil {
			log.Printf("Error creating user: %v", err)
			response.WriteJson(w, http.StatusInternalServerError, response.Response{
				Status: response.Error,
				Error:  err.Error(),
			})
			return
		}

		// Generate JWT token
		key, err := jwttoken.CreateToken(createdUser.ID, createdUser.Name)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.Response{
				Status: response.Error,
				Error:  "Failed to create token",
			})
			return
		}

		// Send success response
		response.WriteJson(w, http.StatusCreated, response.Response{
			Status:  response.Success,
			Error:   "",
			Message: key,
		})
	}
}

func Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user types.User

		// Decode incoming JSON
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.Response{
				Status: response.Error,
				Error:  "Invalid request payload",
			})
			return
		}

		// Find the user in the DB
		foundUser, err := dbclient.PrismaClient.User.FindUnique(
			db.User.Email.Equals(user.Email),
		).Exec(r.Context())
		if err != nil {
			response.WriteJson(w, http.StatusUnauthorized, response.Response{
				Status: response.Error,
				Error:  "Invalid email or password",
			})
			return
		}

		// Verify the password
		password, _ := foundUser.Password()
		err = jwttoken.ComparePasswords(password, user.Password)
		if err != nil {
			response.WriteJson(w, http.StatusUnauthorized, response.Response{
				Status: response.Error,
				Error:  "Invalid email or password",
			})
			return
		}

		key, err := jwttoken.CreateToken(foundUser.ID, foundUser.Name)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.Response{
				Status: response.Error,
				Error:  "Failed to create token",
			})
			return
		}

		response.WriteJson(w, http.StatusOK, response.Response{
			Status:  response.Success,
			Error:   "",
			Message: key,
		})
	}
}

func Verify() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// just taking the token from the header
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			response.WriteJson(w, http.StatusUnauthorized, response.Response{
				Status: response.Error,
				Error:  "Missing authorization header",
			})
			return
		}

		tokenString = tokenString[len("Bearer "):]
		//  veification of teken
		userID, err := jwttoken.VerifyToken(tokenString)
		if err != nil {
			response.WriteJson(w, http.StatusUnauthorized, response.Response{
				Status: response.Error,
				Error:  "Invalid token",
			})
			return
		}

		user, err := dbclient.PrismaClient.User.FindUnique(
			db.User.ID.Equals(userID),
		).Exec(r.Context())
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.Response{
				Status: response.Error,
				Error:  "Failed to fetch user",
			})
			return
		}
		// Return the user's data
		response.WriteJson(w, http.StatusOK, user)
	}
}
