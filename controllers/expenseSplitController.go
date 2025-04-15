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

func CreateExpenseSplitHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Extract expenseID from URL path
		params := mux.Vars(r)
		expenseIDStr := params["expenseID"]
		expenseID, err := strconv.Atoi(expenseIDStr)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		var payload struct {
			UserID int     `json:"user_id"`
			Amount float64 `json:"amount"`
		}
		err = json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		newSplit, err := models.CreateExpenseSplit(db, expenseID, payload.UserID, payload.Amount)
		if err != nil {
			fmt.Println("Error creating expense user:", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newSplit)
	}
}

func GetSplitsByExpenseIDHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		params := mux.Vars(r)
		expenseIDStr := params["expenseID"]
		expenseID, err := strconv.Atoi(expenseIDStr)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		splits, err := models.GetSplitsByExpenseID(db, expenseID)
		if err != nil {
			fmt.Println("Error fetching expense splits:", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(splits)
	}
}
