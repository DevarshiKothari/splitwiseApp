package models

import (
	"database/sql"
	"fmt"
)

type ExpenseSplit struct {
	ID        int     `json:"id"`
	ExpenseID int     `json:"expense_id"`
	UserID    int     `json:"user_id"`
	Amount    float64 `json:"amount"`
}

func CreateExpenseSplit(db *sql.DB, expenseID int, userID int, amount float64) (ExpenseSplit, error) {
	var expenseSplit ExpenseSplit

	query := `INSERT INTO expense_splits (expense_id, user_id, amount) 
	VALUES ($1, $2, $3)
	RETURNING id, expense_id, user_id, amount;`

	err := db.QueryRow(query, expenseID, userID, amount).Scan(
		&expenseSplit.ID,
		&expenseSplit.ExpenseID,
		&expenseSplit.UserID,
		&expenseSplit.Amount,
	)
	if err != nil {
		return ExpenseSplit{}, err
	}

	// Fetch group_id from expenses table
	var groupID int
	err = db.QueryRow("SELECT group_id FROM expenses WHERE id = $1", expenseID).Scan(&groupID)
	if err != nil {
		fmt.Println("Failed to fetch group ID for activity log:", err)
	}

	description := fmt.Sprintf("Expense split of %f added by user %d for expense %d in group %d", amount, userID, expenseID, groupID)
	err = CreateActivity(db, groupID, userID, "expense_split_created", description)
	if err != nil {
		fmt.Println("Activity logging failed for expense split creation:", err)
	}

	return expenseSplit, nil
}

func GetSplitsByExpenseID(db *sql.DB, expenseID int) ([]ExpenseSplit, error) { // Fetch all splits for a given expense
	var splits []ExpenseSplit

	query := `SELECT id, expense_id, user_id, amount 
		FROM expense_splits
		WHERE expense_id = $1;`

	rows, err := db.Query(query, expenseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var split ExpenseSplit
		err := rows.Scan(&split.ID, &split.ExpenseID, &split.UserID, &split.Amount)
		if err != nil {
			return nil, err
		}

		splits = append(splits, split)
	}

	return splits, nil
}
