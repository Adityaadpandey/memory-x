package users

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/adityaadpandey/memory-x/go-api/internal/dbclient"
	"github.com/adityaadpandey/memory-x/go-api/internal/types"
	"github.com/adityaadpandey/memory-x/go-api/internal/utils/jwttoken"
	response "github.com/adityaadpandey/memory-x/go-api/internal/utils/response"
	"github.com/adityaadpandey/memory-x/go-api/prisma/db"
	"github.com/golang-jwt/jwt/v5"
)

func Post() http.HandlerFunc {
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

		// Create the user in the DB and capture the result
		createdUser, err := dbclient.PrismaClient.User.CreateOne(
			db.User.Name.Set(user.Name),
			db.User.Email.Set(user.Email),
			db.User.Password.Set(user.Password),
		).Exec(r.Context())
		if err != nil {
			log.Printf("Error creating user: %v", err)
			response.WriteJson(w, http.StatusInternalServerError, response.Response{
				Status: response.Error,
				Error:  "Failed to create user",
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

func Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Step 1: Get the token from the Authorization header
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			response.WriteJson(w, http.StatusUnauthorized, response.Response{
				Status: response.Error,
				Error:  "Missing authorization header",
			})
			return
		}

		// Remove "Bearer " prefix
		tokenString = tokenString[len("Bearer "):]

		// Step 2: Parse the token and verify it
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwttoken.SecretKey, nil
		})
		if err != nil || !token.Valid {
			response.WriteJson(w, http.StatusUnauthorized, response.Response{
				Status: response.Error,
				Error:  "Invalid token",
			})
			return
		}

		// Step 3: Get claims (user ID, username)
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["id"] == nil {
			response.WriteJson(w, http.StatusUnauthorized, response.Response{
				Status: response.Error,
				Error:  "Invalid token claims",
			})
			return
		}
		fmt.Println(claims)

		// Step 4: Fetch the user by ID from the DB
		userID := claims["id"].(string)
		user, err := dbclient.PrismaClient.User.FindUnique(
			db.User.ID.Equals(userID),
		).Exec(r.Context())
		fmt.Printf("User: %v\n", user)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.Response{
				Status: response.Error,
				Error:  "Failed to fetch user",
			})
			return
		}

		// Step 5: Return the user's data
		response.WriteJson(w, http.StatusOK, user)
	}
}
