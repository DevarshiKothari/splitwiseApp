package main

import (
	"log"
	"splitwise-app/config"
	"splitwise-app/migrations"
)

func main() {
	config.ConnectDB()
	defer config.DB.Close()

	err := config.DB.Ping()
	if err != nil {
		log.Fatal("Database connection failed: ", err)
	}

	migrations.RunMigrations(config.DB)
}
