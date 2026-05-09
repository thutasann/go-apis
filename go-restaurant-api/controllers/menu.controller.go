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
	"go.mongodb.org/mongo-driver/mongo/options"
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
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var menu models.Menu
		if err := c.BindJSON(&menu); err != nil {
			helpers.Error(c, "JSON bind error", 0, err)
			return
		}

		menuId := c.Param("menu_id")
		filter := bson.M{"menu_id": menuId}

		var updateObj primitive.D

		if menu.Start_Date != nil && menu.End_Date != nil {
			if !inTimeSpan(*menu.Start_Date, *menu.End_Date, time.Now()) {
				msg := "kindly retype the time"
				helpers.Error(c, msg, 0, nil)
				return
			}

			updateObj = append(updateObj, bson.E{Key: "start_date", Value: menu.Start_Date})
			updateObj = append(updateObj, bson.E{Key: "end_date", Value: menu.End_Date})
		}

		if menu.Name != "" {
			updateObj = append(updateObj, bson.E{Key: "name", Value: menu.Name})
		}

		if menu.Category != "" {
			updateObj = append(updateObj, bson.E{Key: "category", Value: menu.Category})
		}

		menu.Updated_at = time.Now().UTC()
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: menu.Updated_at})

		upsert := true
		opt := options.UpdateOptions{Upsert: &upsert}

		result, err := menuCollection.UpdateOne(ctx, filter, bson.D{
			{Key: "$set", Value: updateObj},
		}, &opt)

		if err != nil {
			msg := "Menu update failed"
			helpers.Error(c, msg, 0, err)
			return
		}

		helpers.Success(c, "Update Menu Success", result)
	}
}

// Private: Function to Check the Start Date and End Date is in the time span
func inTimeSpan(start time.Time, end time.Time, now time.Time) bool {
	return start.After(now) && end.After(start)
}
