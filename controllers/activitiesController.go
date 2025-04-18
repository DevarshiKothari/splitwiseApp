package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"splitwise-app/models"
	"strconv"

	"github.com/gorilla/mux"
)

func GetGroupActivitiesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)
		groupIDStr := vars["groupID"]
		groupID, err := strconv.Atoi(groupIDStr)
		if err != nil {
			http.Error(w, "Invalid group ID", http.StatusBadRequest)
			return
		}

		activities, err := models.GetActivitiesByGroupID(db, groupID)
		if err != nil {
			http.Error(w, "Failed to fetch activities", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(activities)
	}
}
