package main

import (
	"log"
	"myproject/config"
	"myproject/database"
)

func main() {
    // Initialize configuration
    dbConfig := &config.DBConfig{
        Host:     "localhost",
        Port:     "5432",
        User:     "postgres",
        Password: "password",
        DBName:   "myapp",
        SSLMode:  "disable",
    }

    // Create database connection
    db, err := database.NewDatabase(dbConfig)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    // Auto Migrate schemas
    err = db.AutoMigrate(&User{}, &Profile{})
    if err != nil {
        log.Fatalf("Failed to migrate database: %v", err)
    }
}