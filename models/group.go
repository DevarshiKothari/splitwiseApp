package models

import (
	"database/sql"
	"time"
)

type Group struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedBy int       `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

func CreateGroup(db *sql.DB, name string, createdBy int) (Group, error) {
	var group Group
	query := `
		INSERT INTO groups (name, createdBy)
		VALUES ($1, $2)
		RETURNING id, name, createdBy, created_at;
	`
	err := db.QueryRow(query, name, createdBy).Scan(
		&group.ID,
		&group.Name,
		&group.CreatedBy,
		&group.CreatedAt,
	)
	if err != nil {
		return Group{}, err
	}

	return group, nil
}

func GetGroupByID(db *sql.DB, groupID int) (Group, error) {
	var group Group
	query := `
		SELECT (groupID) from groups
		VALUES ($1)
	`
	err := db.QueryRow(query, groupID).Scan(
		&group.ID,
		&group.Name,
		&group.CreatedBy,
		&group.CreatedAt,
	)
	if err != nil {
		return Group{}, err
	}

	return group, nil
}
