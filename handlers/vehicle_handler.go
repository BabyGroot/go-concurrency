package handlers

import (
	"errors"
	"myproject/repositories"
	"net/http"

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

	// Get vehicles by location ID
	vehicles, err := h.repo.GetVehiclesByLocation(location.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	type VehicleResponse struct {
		Permalink string `json:"permalink"`
	}

	vehicles_response := make([]VehicleResponse, len(vehicles))
	for i, vehicle := range vehicles {
		vehicles_response[i] = VehicleResponse{
			Permalink: vehicle.Permalink,
		}
	}

	c.JSON(http.StatusOK, vehicles_response)
}
