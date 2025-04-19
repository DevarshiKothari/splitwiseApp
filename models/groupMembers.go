package models

import (
	"database/sql"
	"fmt"
)

type GroupMember struct {
	ID      int `json:"id"`
	GroupID int `json:"group_id"`
	UserID  int `json:"user_id"`
}

func AddGroupMember(db *sql.DB, groupID int, userID int) (GroupMember, error) {
	// Check if user is already a member of the group
	var exists bool
	checkQuery := `
		SELECT EXISTS (
			SELECT 1 FROM group_members 
			WHERE group_id = $1 AND user_id = $2
		);
	`

	err := db.QueryRow(checkQuery, groupID, userID).Scan(&exists)
	if err != nil {
		return GroupMember{}, fmt.Errorf("failed to check membership: %w", err)
	}

	if exists {
		return GroupMember{}, fmt.Errorf("user %d is already a member of group %d", userID, groupID)
	}

	var groupMember GroupMember
	query := `
		INSERT INTO group_members (group_id, user_id)
		VALUES ($1, $2)
		RETURNING id, group_id, user_id;
		`

	err = db.QueryRow(query, groupID, userID).Scan(
		&groupMember.ID,
		&groupMember.GroupID,
		&groupMember.UserID,
	)
	if err != nil {
		// fmt.Printf()
		return GroupMember{}, err
	}

	msg := fmt.Sprintf("User %d added to group %d", userID, groupID)
	err = CreateActivity(db, groupID, userID, "member_added", msg)
	if err != nil {
		fmt.Println("Activity logging failed for member addition:", err)
	}

	return groupMember, nil
}

func GetGroupMembers(db *sql.DB, groupID int) ([]User, error) {
	var members []User

	query := `
		SELECT u.id, u.name, u.email, u.created_at
		FROM group_members gm
		JOIN users u ON gm.user_id = u.id
		WHERE gm.group_id = $1;
	`
	rows, err := db.Query(query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() { // cannot range over rows, because it is of type *sql.Rows
		var user User

		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		members = append(members, user)
	}

	return members, nil
}
