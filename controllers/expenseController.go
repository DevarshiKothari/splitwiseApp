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

func CreateExpenseHandler(db *sql.DB) http.HandlerFunc {
	// fmt.Println("Handler hit")
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		params := mux.Vars(r)
		groupIDStr := params["groupID"]
		groupID, err := strconv.Atoi(groupIDStr)
		if err != nil {
			http.Error(w, "Invalid Group ID", http.StatusBadRequest)
			return
		}

		var payload struct {
			PaidBy      int     `json:"paid_by"`
			Description string  `json:"description"`
			TotalAmount float64 `json:"total_amount"`
		}
		err = json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		expense, err := models.CreateExpense(db, groupID, payload.PaidBy, payload.Description, payload.TotalAmount)
		if err != nil {
			fmt.Println("Error creating expense:", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(expense)
	}
}

func GetExpenseByIDHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// id := r.URL.Query().Get("expenseID") // This method parses the query parameters and not from the URL path of the incoming request. It returns a url.Values object, which is essentially a map of query parameter keys to their values.
		params := mux.Vars(r)
		idStr := params["expenseID"]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid expense ID", http.StatusBadRequest)
			return
		}

		expense, err := models.GetExpenseByID(db, id)
		if err != nil {
			fmt.Println("Error fetching expense:", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(expense)
	}
}

func GetExpensesByGroupIDHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		params := mux.Vars(r)
		groupIDStr := params["groupID"]
		groupID, err := strconv.Atoi(groupIDStr)
		if err != nil {
			http.Error(w, "Invalid Group ID", http.StatusBadRequest)
			return
		}

		expenses, err := models.GetExpensesByGroupID(db, groupID)
		if err != nil {
			fmt.Println("Error fetching expenses:", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(expenses)
	}
}
