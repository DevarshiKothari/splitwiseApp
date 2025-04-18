package models

import (
	"database/sql"
	"fmt"
	"splitwise-app/utils"
)

func CalculateGroupBalances(db *sql.DB, groupID int) ([]utils.Balance, error) {
	expenseRows, err := db.Query("SELECT id, paid_by, total_amount FROM expenses WHERE group_id = $1", groupID)
	if err != nil {
		return nil, err
	}
	defer expenseRows.Close()

	type Expense struct {
		ID          int
		PaidBy      int
		TotalAmount float64
	}

	var expenses []Expense
	var expenseIDs []int

	for expenseRows.Next() {
		var e Expense
		if err := expenseRows.Scan(&e.ID, &e.PaidBy, &e.TotalAmount); err != nil {
			return nil, err
		}
		expenses = append(expenses, e)
		expenseIDs = append(expenseIDs, e.ID)
	}

	if len(expenseIDs) == 0 {
		return []utils.Balance{}, nil // [] is not a valid way to return empty slice in Go so []Balance{}
	}

	// Fetch all splits for these expenses
	placeholders := make([]string, len(expenseIDs))
	args := make([]interface{}, len(expenseIDs))
	for i, id := range expenseIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}
	// Slice of type int containing expenses ids for a group, query so that it will accept this slice
	query := fmt.Sprintf("SELECT expense_id, user_id, amount FROM expense_splits WHERE expense_id IN (%s)",
		utils.StringJoin(placeholders, ",")) // this prepares a query like SELECT ... FROM expense_splits WHERE expense_id IN ($1, $2, $3...)

	splitRows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer splitRows.Close()

	type Split struct {
		ExpenseID int
		UserID    int
		Amount    float64
	}

	splitsMap := make(map[int][]Split)
	for splitRows.Next() {
		var s Split
		err := splitRows.Scan(
			&s.ExpenseID,
			&s.UserID,
			&s.Amount,
		)
		if err != nil {
			return nil, err
		}
		splitsMap[s.ExpenseID] = append(splitsMap[s.ExpenseID], s)
	}

	// Compute balances
	// balance[user_id] = amountOwed
	balances := make(map[int]float64)

	for _, expense := range expenses {
		balances[expense.PaidBy] += expense.TotalAmount // the payer gets credit

		for _, split := range splitsMap[expense.ID] {
			balances[split.UserID] -= split.Amount // debtor owes money
		}
	}

	// Fetch and subtract settlements
	settlementRows, err := db.Query(
		"SELECT from_user_id, to_user_id, amount FROM settlements WHERE group_id = $1",
		groupID,
	)
	if err != nil {
		fmt.Println("Error fetching settlements:", err)
		return nil, err
	}

	defer settlementRows.Close()

	for settlementRows.Next() {
		var from_user_id, to_user_id int
		var amount float64

		err := settlementRows.Scan(&from_user_id, &to_user_id, &amount)
		if err != nil {
			return nil, err
		}

		balances[from_user_id] += amount
		balances[to_user_id] -= amount
	}

	return utils.SimplifyBalances(balances), nil
}
