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

		var groupMember models.GroupMember
		err := json.NewDecoder(r.Body).Decode(&groupMember)
		if err != nil {
			http.Error(w, "Invalid req payload", http.StatusBadRequest)
			return
		}

		addedGroupMember, err := models.AddGroupMember(db, groupMember.GroupID, groupMember.UserID)
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
		idStr := params["id"]

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
