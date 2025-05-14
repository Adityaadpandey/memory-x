package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/adityaadpandey/memory-x/go-api/internal/config"
	"github.com/adityaadpandey/memory-x/go-api/internal/dbclient"
	users "github.com/adityaadpandey/memory-x/go-api/internal/handlers"
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

	router.HandleFunc("POST /test", users.Post())
	router.HandleFunc("GET /test", users.Get())

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
