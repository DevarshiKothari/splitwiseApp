package migrations

import (
	"database/sql" // Go's standard package to interact with SQL databases (like PostgreSQL, MySQL).
	"fmt"
)

func RunMigrations(db *sql.DB) { //*sql.DB is a a reference to a database connection pool provided by Go's standard library (database/sql). It represents a pool of zero or more underlying connections to our actual SQL database (like PostgreSQL). It allows our app to send SQL queries to the actual database efficiently
	queries := []string{ // Added ON DELETE CASCADE to foreign keys to allow deletion of parent records and automatically delete child rows
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,

		`CREATE TABLE IF NOT EXISTS groups (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			created_by INTEGER REFERENCES users(id) ON DELETE CASCADE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,

		`CREATE TABLE IF NOT EXISTS group_members (
			id SERIAL PRIMARY KEY,
			group_id INTEGER REFERENCES groups(id) ON DELETE CASCADE,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE
		);`,

		`CREATE TABLE IF NOT EXISTS expenses (
			id SERIAL PRIMARY KEY,
			group_id INTEGER REFERENCES groups(id) ON DELETE CASCADE,
			paid_by INTEGER REFERENCES users(id) ON DELETE CASCADE,
			total_amount NUMERIC NOT NULL,
			description TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,

		`CREATE TABLE IF NOT EXISTS expense_splits (
			id SERIAL PRIMARY KEY,
			expense_id INTEGER REFERENCES expenses(id) ON DELETE CASCADE,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			amount NUMERIC NOT NULL
		);`,

		`CREATE TABLE IF NOT EXISTS settlements (
			id SERIAL PRIMARY KEY,
			group_id INTEGER NOT NULL REFERENCES groups(id),
			from_user_id INTEGER NOT NULL REFERENCES users(id),
			to_user_id INTEGER NOT NULL REFERENCES users(id),
			amount NUMERIC(10,2) NOT NULL CHECK (amount > 0),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS activities (
			id SERIAL PRIMARY KEY,
			group_id INT NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
			user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE, -- who performed the action
			activity_type TEXT NOT NULL, -- "expense_added", "split_added", "settlement", "member_added"
			description TEXT NOT NULL,   -- human-readable description like "Riya was added to the group"
			created_at TIMESTAMP DEFAULT NOW()
		);`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			fmt.Printf("Error running query:\n%s\nError: %v\n", query, err)
		}
	}

	fmt.Println("All migrations applied.")
}
