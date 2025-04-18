package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Expense struct {
	ID          int       `json:"id"`
	GroupID     int       `json:"group_id"`
	PaidBy      int       `json:"paid_by"`
	Description string    `json:"description"`
	TotalAmount float64   `json:"total_amount"`
	CreatedAt   time.Time `json:"created_at"`
}

func CreateExpense(db *sql.DB, groupId int, paidBy int, description string, totalAmount float64) (Expense, error) {
	var expense Expense

	query := `INSERT INTO expenses (group_id, paid_by, description, total_amount) 
	VALUES ($1, $2, $3, $4)
	RETURNING id, group_id, paid_by, description, total_amount, created_at;`

	err := db.QueryRow(query, groupId, paidBy, description, totalAmount).Scan(
		&expense.ID,
		&expense.GroupID,
		&expense.PaidBy,
		&expense.Description,
		&expense.TotalAmount,
		&expense.CreatedAt,
	)
	if err != nil {
		return Expense{}, err
	}

	msg := fmt.Sprintf("%s Expense added by %d in %d", description, paidBy, groupId)
	err = CreateActivity(db, groupId, paidBy, "expense_created", msg)
	if err != nil {
		fmt.Println("Activity logging failed for expense creation:", err)
	}

	return expense, nil
}

func GetExpenseByID(db *sql.DB, id int) (Expense, error) {
	var expense Expense

	query := `SELECT id, group_id, paid_by, description, total_amount, created_at 
		FROM expenses
		WHERE id = $1;`

	err := db.QueryRow(query, id).Scan(
		&expense.ID,
		&expense.GroupID,
		&expense.PaidBy,
		&expense.Description,
		&expense.TotalAmount,
		&expense.CreatedAt,
	)
	if err != nil {
		return Expense{}, err
	}

	return expense, nil
}

func GetExpensesByGroupID(db *sql.DB, groupID int) ([]Expense, error) {
	var expenses []Expense

	query := `SELECT id, group_id, paid_by, description, total_amount, created_at 
		FROM expenses
		WHERE group_id = $1;`

	rows, err := db.Query(query, groupID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var expense Expense
		err := rows.Scan(
			&expense.ID,
			&expense.GroupID,
			&expense.PaidBy,
			&expense.Description,
			&expense.TotalAmount,
			&expense.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		expenses = append(expenses, expense)
	}

	return expenses, nil
}
