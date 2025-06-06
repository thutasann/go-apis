package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/thutasann/restaurant-api/controllers"
)

func OrderRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/orders", controllers.GetOrders())
	incomingRoutes.GET("/order/:order_id", controllers.GetOrder())
	incomingRoutes.POST("/orders", controllers.CreateOrder())
	incomingRoutes.PATCH("/order/:order_id", controllers.UpdateOrder())
}
