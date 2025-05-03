package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/thutasann/restaurant-api/controllers"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/users", controllers.GetUsers())
	incomingRoutes.GET("/users/:user_id", controllers.GetUser())
	incomingRoutes.POST("/users/signup", controllers.SignUp())
	incomingRoutes.POST("/users/signin", controllers.SignIn())
}
