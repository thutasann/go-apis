package controllers

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/thutasann/ecommerce-cart/pkg/database"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var UserCollection *mongo.Collection = database.UserData(database.Client, "Users")
var ProductCollection *mongo.Collection = database.ProductData(database.Client, "Products")

// Hash Password
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

// Verify Password
func verifyPassword(userPassword string, givenPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(givenPassword), []byte(userPassword))
	valid := true
	msg := ""
	if err != nil {
		msg = "Login or Password is Incorrect"
		valid = false
	}
	return valid, msg
}

// SignUp Controller
func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

// Login Controller
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

// Product Viewer Admin Controller
func ProductViewerAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

// Search Product Controller
func SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

// Search Product By Query Controller
func SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
