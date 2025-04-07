package handlers

import (
	"myproject/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LocationHandler struct {
	repo *repositories.LocationRepository
}

func NewLocationHandler(repo *repositories.LocationRepository) *LocationHandler {
	return &LocationHandler{
		repo: repo,
	}
}

func (h *LocationHandler) GetAll(c *gin.Context) {
	locations, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	permalinks := make([]string, len(locations))
	for i, location := range locations {
		permalinks[i] = location.Permalink
	}

	c.JSON(http.StatusOK, permalinks)
}
