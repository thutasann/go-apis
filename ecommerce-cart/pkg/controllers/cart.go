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

}

// Remove Item
func (app *Application) RemoveItem() gin.HandlerFunc {

}

// Get Item From Cart
func (app *Application) GetItemFromCart() gin.HandlerFunc {

}

// Buy Cart
func (app *Application) BuyFromCart() gin.HandlerFunc {

}

// Instant Buy
func (app *Application) InstantBuy() gin.HandlerFunc {

}
