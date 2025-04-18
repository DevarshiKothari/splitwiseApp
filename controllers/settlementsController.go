package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"splitwise-app/models"
	"splitwise-app/utils"

	"github.com/gorilla/mux"
)

func AddSettlement(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)
		groupIDStr := vars["groupID"]
		groupID, err := strconv.Atoi(groupIDStr)
		if err != nil {
			http.Error(w, "Invalid group ID", http.StatusBadRequest)
			return
		}

		var req []utils.Balance
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil { // r.Body is a stream — read it only once
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err = models.SaveSettlementsToDB(db, groupID, req)
		if err != nil {
			http.Error(w, "Failed to record settlement", http.StatusInternalServerError)
			return
		}

		// Return updated balances
		balances, err := models.CalculateGroupBalances(db, groupID)
		if err != nil {
			http.Error(w, "Failed to calculate updated balances", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(balances)
	}
}
