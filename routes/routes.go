package routes

import (
	"myproject/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, vh *handlers.VehicleHandler) {
	// Group routes
	api := r.Group("/api")
	{
		api.GET("/vehicles", vh.GetAll)
	}
}
