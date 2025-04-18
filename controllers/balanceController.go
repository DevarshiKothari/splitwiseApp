package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"splitwise-app/models"
	"strconv"

	"github.com/gorilla/mux"
)

func GetGroupBalanceHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		params := mux.Vars(r)
		groupIDStr := params["groupID"]
		groupID, err := strconv.Atoi(groupIDStr)
		if err != nil {
			http.Error(w, "Invalid group ID", http.StatusBadRequest)
			return
		}

		balances, err := models.CalculateGroupBalances(db, groupID)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Group ID not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(balances)
	}
}
