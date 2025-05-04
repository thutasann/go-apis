package controllers

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thutasann/restaurant-api/database"
	"github.com/thutasann/restaurant-api/helpers"
	"github.com/thutasann/restaurant-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

		menuId := c.Param("menu_id")
		var menu models.Menu

		err := menuCollection.FindOne(ctx, bson.M{"menu_id": menuId}).Decode(&menu)
		defer cancel()
		if err != nil {
			helpers.Error(c, "Error occured while fetching menu by ID", 0, err)
			return
		}
		helpers.Success(c, "Food fetched successfully", menu)
	}
}

// Create Menu
func CreateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		var menu models.Menu
		var ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)

		if err := c.BindJSON(&menu); err != nil {
			helpers.Error(c, "JSON bind err", 0, err)
		}

		validationErr := validate.Struct(menu)
		if validationErr != nil {
			helpers.Error(c, "Validation error", 0, validationErr)
		}

		menu.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.ID = primitive.NewObjectID()
		menu.Menu_id = menu.ID.Hex()

		result, insertErr := menuCollection.InsertOne(ctx, menu)
		if insertErr != nil {
			helpers.Error(c, "Insert err", 0, insertErr)
		}
		defer cancel()
		helpers.Success(c, "Food was created", result)
	}
}

// Update Menu
func UpdateMenu() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}
