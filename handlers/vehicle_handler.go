package handlers

import (
	"errors"
	"myproject/models"
	"myproject/repositories"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type VehicleHandler struct {
	repo         *repositories.VehicleRepository
	locationRepo *repositories.LocationRepository
}

func NewVehicleHandler(repo *repositories.VehicleRepository, locationRepo *repositories.LocationRepository) *VehicleHandler {
	return &VehicleHandler{
		repo:         repo,
		locationRepo: locationRepo,
	}
}

func (h *VehicleHandler) GetAll(c *gin.Context) {
	vehicles, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, vehicles)
}

func (h *VehicleHandler) GetVehiclesByLocation(c *gin.Context) {
	// Get location name from path parameter
	locationPermalink := c.Param("name")

	// First find the location
	location, err := h.locationRepo.FindByPermalink(locationPermalink)
	if err != nil {
		// Check if it's a "record not found" error from GORM
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Location not found"})
			return
		}

		// Other database errors
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	vehiclesChan := make(chan []models.Vehicle)
	statsChan := make(chan map[string]int)
	errorChan := make(chan error, 2) // Buffer to collect potential errors

	// Get vehicles concurrently
	go func() {
		vehicles, err := h.repo.GetVehiclesByLocation(location.ID)
		if err != nil {
			errorChan <- err
			return
		}
		vehiclesChan <- vehicles
	}()

	// Get statistics concurrently (e.g., count by vehicle type)
	go func() {
		stats, err := h.repo.GetVehicleStatsByLocation(location.ID)
		if err != nil {
			errorChan <- err
			return
		}
		statsChan <- stats
	}()

	// Collect results with timeout
	var vehicles []models.Vehicle
	var stats map[string]int

	timeout := time.After(3 * time.Second)

	for i := 0; i < 2; i++ { // We expect 2 results
		select {
		case err := <-errorChan:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		case vehicles = <-vehiclesChan:
			// Got vehicles
		case stats = <-statsChan:
			// Got stats
		case <-timeout:
			c.JSON(http.StatusRequestTimeout, gin.H{"error": "Request timed out"})
			return
		}
	}

	// Transform to response format
	type VehicleResponse struct {
		Permalink string `json:"permalink"`
	}

	vehiclesResponse := make([]VehicleResponse, len(vehicles))
	for i, vehicle := range vehicles {
		vehiclesResponse[i] = VehicleResponse{
			Permalink: vehicle.Permalink,
		}
	}

	// Return combined response
	c.JSON(http.StatusOK, gin.H{
		"vehicles": vehiclesResponse,
		"stats":    stats,
	})
}
