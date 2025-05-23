package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/thutasann/restaurant-api/controllers"
)

// Food Routes
func FoodRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/foods", controllers.GetFoods())
	incomingRoutes.GET("/foods/:food_id", controllers.GetFood())
	incomingRoutes.POST("/foods", controllers.CreateFood())
	incomingRoutes.PATCH("/foods/:food_id", controllers.UpdateFood())
}
