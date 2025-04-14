package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"splitwise-app/models"

	"github.com/gorilla/mux"
)

// GetUserByIDHandler handles GET /users/{id}
func GetUserByIDHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r) //fetches path parameters like /users/{id}. // mux.Vars(r) gives you a map of URL path variables
		idStr := params["id"]

		id, err := strconv.Atoi(idStr) //Converts that id from string to integer. // Atoi stands for ASCII to Integer
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		user, err := models.GetUserByID(db, id)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

func CreateUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var userInput models.User
		err := json.NewDecoder(r.Body).Decode(&userInput) //Decodes the JSON request body into the user struct.
		if err != nil {
			http.Error(w, "Invalid req payload", http.StatusBadRequest)
			return
		}

		createdUser, err := models.CreateUser(db, userInput.Name, userInput.Email) //Inserts the new user into the database.
		if err != nil {
			fmt.Println("Error creating user:", err) // to print the error
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)      //Sets the HTTP status code to 201 Created.
		json.NewEncoder(w).Encode(createdUser) //Encodes the createdUser back to JSON and sends it in the response.
	}
}
