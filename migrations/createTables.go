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
			amount NUMERIC NOT NULL,
			description TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,

		`CREATE TABLE IF NOT EXISTS expense_splits (
			id SERIAL PRIMARY KEY,
			expense_id INTEGER REFERENCES expenses(id) ON DELETE CASCADE,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			amount NUMERIC NOT NULL
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
