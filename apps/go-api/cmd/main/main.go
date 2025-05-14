package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/adityaadpandey/memory-x/go-api/internal/config"
	"github.com/adityaadpandey/memory-x/go-api/internal/dbclient"
	"github.com/adityaadpandey/memory-x/go-api/prisma/db"
)

func main() {
	cfg := config.MustLoad()

	// Initialize Prisma client
	if err := dbclient.InitClient(); err != nil {
		log.Fatal("Error initializing Prisma client:", err)
	}
	defer dbclient.Disconnect()

	// Create a new HTTP router
	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		// Access the Prisma client via dbclient.PrismaClient
		users, err := dbclient.PrismaClient.User.FindMany().Exec(context.Background())
		if err != nil {
			log.Fatal("Error fetching users:", err)
		}

		// Print the fetched users (or handle them as needed)
		fmt.Println(users)

		// Send a simple response
		w.Write([]byte("Hello, World!"))
	})

	router.HandleFunc("GET /test", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		// create a new user
		user, err := dbclient.PrismaClient.User.CreateOne(
			db.User.Name.Set("Test User"),
			db.User.Email.Set("x@gmail.com"),
			db.User.Password.Set("password"),
		).Exec(context.Background())

		if err != nil {
			return
		}
		// Print the created user (or handle it as needed)
		fmt.Println(user)
		// Send a simple response
		w.Write([]byte("User created successfully!"))
	})

	// Set up the server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	// Start the server
	fmt.Println("Server starting on", cfg.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
