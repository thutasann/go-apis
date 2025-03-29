package controllers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	prodCollection *mongo.Collection
	userCollection *mongo.Collection
}

// Cart Application
func NewApplication(prodCollection, userCollection *mongo.Collection) *Application {
	return &Application{
		prodCollection: prodCollection,
		userCollection: userCollection,
	}
}

// Add to cart
func (app *Application) AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

// Remove Item
func (app *Application) RemoveItem() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

// Get Item From Cart
func (app *Application) GetItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

// Buy Cart
func (app *Application) BuyFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

// Instant Buy
func (app *Application) InstantBuy() gin.HandlerFunc {
	return func(c *gin.Context) {}
}
