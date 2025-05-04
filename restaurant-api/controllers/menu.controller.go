package controllers

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thutasann/restaurant-api/database"
	"github.com/thutasann/restaurant-api/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Menu Mongo Collection
var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")

// Get all Menus
func GetMenus() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		result, err := menuCollection.Find(context.TODO(), bson.M{})
		defer cancel()

		if err != nil {
			helpers.Error(c, "error occured while listing the menu list", 0, err)
		}
		allMenus := make([]bson.M, 0)
		if err = result.All(ctx, &allMenus); err != nil {
			log.Fatal(err)
		}
		helpers.Success(c, "fetched menus success", allMenus)
	}
}

// Get Menu By ID
func GetMenu() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

// Create Menu
func CreateMenu() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

// Update Menu
func UpdateMenu() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}
