package models

import (
	"database/sql"
	"time"
)

// User struct maps to the 'users' table
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateUser inserts a new user into the DB
// The name and email values are plugged into $1 and $2. The database driver will handle escaping and quoting these values to prevent SQL injection attacks.
// The Scan method is used to read the returned values into the user struct. The order of the arguments in Scan must match the order of the columns in the SELECT statement.
func CreateUser(db *sql.DB, name string, email string) (User, error) {
	var user User
	query := `
		INSERT INTO users (name, email)
		VALUES ($1, $2)
		RETURNING id, name, email, created_at;
	`

	err := db.QueryRow(query, name, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
	)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// GetUserByID fetches a user by ID
func GetUserByID(db *sql.DB, id int) (User, error) {
	var user User
	query := `SELECT id, name, email, created_at FROM users WHERE id = $1`

	err := db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
	)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
