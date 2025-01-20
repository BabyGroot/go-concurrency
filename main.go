package main

import (
	"log"
	"myproject/config"
	"myproject/database"
	"myproject/handlers"
	"myproject/repositories"
	"myproject/routes"

	"github.com/gin-gonic/gin"
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

	// Create handler
    vehicleHandler := handlers.NewVehicleHandler(vehicleRepo)

	// Setup Gin
    r := gin.Default()

    // Routes
    routes.SetupRoutes(r, vehicleHandler)

    // Start server
    log.Println("Server starting on :8080")
    r.Run(":8080") // listen and serve on 0.0.0.0:8080

	// Get all vehicles
	// vehicles, err := vehicleRepo.GetAll()
	// if err != nil {
	// 	log.Fatalf("Failed to fetch vehicles: %v", err)
	// }

	// Print the vehicles
	// for _, vehicle := range vehicles {
	// 	log.Printf("Vehicle: ID=%d, Name=%s, Permalink=%s\n",
	// 		vehicle.ID, vehicle.Name, vehicle.Permalink)
	// }
}
