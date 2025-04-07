package routes

import (
	"myproject/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, vh *handlers.VehicleHandler, lh *handlers.LocationHandler) {
	// Group routes
	api := r.Group("/api")
	{
		api.GET("/vehicles", vh.GetAll)
		api.GET("locations", lh.GetAll)
		api.GET("/locations/:name/vehicles", vh.GetVehiclesByLocation)
	}
}
