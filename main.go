package main

import (
	"log"
	"net/http"

	"splitwise-app/config"
	"splitwise-app/migrations"
	"splitwise-app/routes"

	"github.com/gorilla/mux"
)

func main() {
	config.ConnectDB()
	defer config.DB.Close()

	err := config.DB.Ping() // Check if the database connection is alive
	if err != nil {
		log.Fatal("Database connection failed: ", err)
	}

	// Run table migrations
	migrations.RunMigrations(config.DB)

	// Initialize router // Creates a new router
	router := mux.NewRouter()

	// Register all user-related routes
	routes.RegisterUserRoutes(router, config.DB)

	// 	Starts the HTTP server on port 8080
	log.Println("Server is running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
