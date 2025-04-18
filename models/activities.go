package models

import "database/sql"

type Activity struct {
	ID           int    `json:"id"`
	GroupID      int    `json:"group_id"`
	UserID       int    `json:"user_id"`
	ActivityType string `json:"activity_type"` // e.g., "expense_added", "split_added", "settlement", "member_added"
	Description  string `json:"description"`
	CreatedAt    string `json:"created_at"`
}

func CreateActivity(db *sql.DB, groupID int, userID int, activityType string, description string) error {
	var activity Activity
	query := `
		INSERT INTO activities (group_id, user_id, activity_type, description)
		VALUES ($1, $2, $3, $4)
		RETURNING id, group_id, user_id, activity_type, description, created_at;
	`

	err := db.QueryRow(query, groupID, userID, activityType, description).Scan(
		&activity.ID,
		&activity.GroupID,
		&activity.UserID,
		&activity.ActivityType,
		&activity.Description,
		&activity.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func GetActivitiesByGroupID(db *sql.DB, groupID int) ([]Activity, error) {
	var activities []Activity

	query := `
		SELECT id, group_id, user_id, activity_type, description, created_at
		FROM activities
		WHERE group_id = $1
		ORDER BY created_at DESC;
	`

	rows, err := db.Query(query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var activity Activity
		err := rows.Scan(
			&activity.ID,
			&activity.GroupID,
			&activity.UserID,
			&activity.ActivityType,
			&activity.Description,
			&activity.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}

	return activities, nil
}
