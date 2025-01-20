package handlers

import (
	"myproject/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)

type VehicleHandler struct {
	repo *repositories.VehicleRepository
}

func NewVehicleHandler(repo *repositories.VehicleRepository) *VehicleHandler {
	return &VehicleHandler{repo: repo}
}

func (h *VehicleHandler) GetAll(c *gin.Context) {
	vehicles, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, vehicles)
}
