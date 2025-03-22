package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/thutasann/ecommerce-cart/pkg/controllers"
	"github.com/thutasann/ecommerce-cart/pkg/database"
	"github.com/thutasann/ecommerce-cart/pkg/middleware"
	"github.com/thutasann/ecommerce-cart/pkg/routes"
)

// Ecommerce Cart Rest API
func main() {
	envErr := godotenv.Load("../../.env")
	if envErr != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	app := controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(middleware.Authentication())

	routes.UserRoutes(router)
	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.RemoveItem())
	router.GET("/cartcheckout", app.BuyFromCart())
	router.GET("/instantbuy", app.InstantBuy())

	log.Fatal(router.Run(":" + port))
}
