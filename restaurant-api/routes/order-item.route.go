package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/thutasann/restaurant-api/controllers"
)

func OrderItemRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/order-items", controllers.GetOrderItems())
	incomingRoutes.GET("/order-item/:order_item_id", controllers.GetOrderItem())
	incomingRoutes.GET("/order-item-order/:order_id", controllers.GetOrderItemsByOrder())
	incomingRoutes.POST("/order-items", controllers.CreateOrderItem())
	incomingRoutes.PATCH("/order-item/:order_item_id", controllers.UpdateOrderItem())
}
