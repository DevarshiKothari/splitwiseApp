package models

import (
	"database/sql"
	"splitwise-app/utils"
)

func SaveSettlementsToDB(db *sql.DB, groupID int, settlements []utils.Balance) error {
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

	return nil
}