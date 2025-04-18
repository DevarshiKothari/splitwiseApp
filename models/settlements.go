package models

import (
	"database/sql"
	"fmt"
	"splitwise-app/utils"
)

func AddSettlementsToDB(db *sql.DB, groupID int, settlements []utils.Balance) error {
	query := `
        INSERT INTO settlements (group_id, from_user_id, to_user_id, amount)
        VALUES ($1, $2, $3, $4)
    `

	for _, s := range settlements {
		_, err := db.Exec(query, groupID, s.FromUserID, s.ToUserID, s.Amount)
		if err != nil {
			return err
		}
	}

	msg := fmt.Sprintf("Settlements added for group %d", groupID)
	err := CreateActivity(db, groupID, settlements[0].FromUserID, "settlement_added", msg)
	if err != nil {
		fmt.Println("Activity logging failed for settlement addition:", err)
	}

	return nil
}
