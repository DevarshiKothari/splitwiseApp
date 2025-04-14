package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"splitwise-app/models"
	"strconv"

	"github.com/gorilla/mux"
)

func CreateGroupHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		var group models.Group
		err := json.NewDecoder(r.Body).Decode(&group)
		if err != nil {
			http.Error(w, "Invalid req payload", http.StatusBadRequest)
			return
		}

		createdGroup, err := models.CreateGroup(db, group.Name, group.CreatedBy)
		if err != nil {
			fmt.Println("Error creating user:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdGroup)
	}
}

func GetGroupByIdHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		params := mux.Vars(r)
		idStr := params["id"]

		id, err := strconv.Atoi(idStr)
		if err != nil {
			return
		}

		fetchedGroup, err := models.GetGroupByID(db, id)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(fetchedGroup)
	}
}
