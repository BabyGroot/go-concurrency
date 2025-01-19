package main

import (
	"log"
	"myproject/config"
	"myproject/database"
	"myproject/repositories"
)

func main() {
	// Initialize configuration
	dbConfig := &config.DBConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "test_app",
		Password: "",
		DBName:   "playground",
		SSLMode:  "disable",
	}

	// Create database connection
	db, err := database.NewDatabase(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create repository
	vehicleRepo := repositories.NewVehicleRepository(db)

	// Get all vehicles
	vehicles, err := vehicleRepo.GetAll()
	if err != nil {
		log.Fatalf("Failed to fetch vehicles: %v", err)
	}

	// Print the vehicles
	for _, vehicle := range vehicles {
		log.Printf("Vehicle: ID=%d, Name=%s, Permalink=%s\n",
			vehicle.ID, vehicle.Name, vehicle.Permalink)
	}
}
