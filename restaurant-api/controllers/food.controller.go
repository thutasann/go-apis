package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/thutasann/restaurant-api/database"
	"github.com/thutasann/restaurant-api/helpers"
	"github.com/thutasann/restaurant-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// validator instance
var validate = validator.New()

// Food Mongo Collection
var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

// Get Foods
func GetFoods() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

// Get Food By ID
func GetFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

		foodId := c.Param("food_id")
		var food models.Food

		err := foodCollection.FindOne(ctx, bson.M{"food_id": foodId}).Decode(&food)
		defer cancel()
		if err != nil {
			helpers.Error(c, "Error occured while fetching food by Id", 0, err)
		}
		helpers.Success(c, "Food fetched successfully", food)
	}
}

// Create Food
func CreateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		var menu models.Menu
		var food models.Food

		if err := c.BindJSON(&food); err != nil {
			helpers.Error(c, "JSON Bind error", 0, err)
		}

		validationErr := validate.Struct(food)
		if validationErr != nil {
			helpers.Error(c, "Validation Error", 0, validationErr)
		}

		err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.Menu_id}).Decode(&menu)
		defer cancel()

		if err != nil {
			msg := fmt.Sprintf("menu was not found %s", *food.Menu_id)
			helpers.Error(c, msg, 0, err)
			return
		}

		food.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.ID = primitive.NewObjectID()
		food.Food_id = food.ID.Hex()
		var num = toFixed(*food.Price, 2)
		food.Price = &num

		result, insertErr := foodCollection.InsertOne(ctx, food)
		if insertErr != nil {
			msg := "Food item was not created"
			helpers.Error(c, msg, 0, insertErr)
			return
		}
		defer cancel()
		helpers.Success(c, "Food was created", result)
	}
}

// Update Food
func UpdateFood() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

// Private: round utility
func round(num float64) int {
	return 0
}

// Privae: toFixed utility
func toFixed(num float64, precision int) float64 {
	return 0
}
