package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"splitwise-app/models"

	"github.com/gorilla/mux"
)

// GetUserByIDHandler handles GET /users/{id}
func GetUserByIDHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r) //fetches path parameters like /users/{id}.
		idStr := params["id"]

		id, err := strconv.Atoi(idStr) //Converts that id from string to integer.
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
