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

func AddGroupMemberHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		// Extract groupID from URL path
		params := mux.Vars(r)
		groupIDStr := params["groupID"]
		groupID, err := strconv.Atoi(groupIDStr)
		if err != nil {
			http.Error(w, "Invalid Group ID", http.StatusBadRequest)
			return
		}

		var payload struct {
			UserID int `json:"user_id"`
		}
		err = json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		addedGroupMember, err := models.AddGroupMember(db, groupID, payload.UserID)
		if err != nil {
			fmt.Println("Error creating user:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(addedGroupMember)
	}
}

func GetGroupMembersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		params := mux.Vars(r)
		idStr := params["groupID"]

		id, err := strconv.Atoi(idStr)
		if err != nil {
			return
		}

		listOfGroupMembers, err := models.GetGroupMembers(db, id)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Group Members not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(listOfGroupMembers)
	}
}
