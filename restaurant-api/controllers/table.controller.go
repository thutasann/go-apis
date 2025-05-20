package controllers

import (
	"context"
	"fmt"
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

// Table Mongo Collection
var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")

// Get Tables
func GetTables() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		result, err := tableCollection.Find(context.TODO(), bson.M{})
		if err != nil {
			helpers.Error(c, "get tables failed", 0, err)
			return
		}
		var allTables []bson.M
		if err = result.All(ctx, &allTables); err != nil {
			log.Fatal(err)
		}
		helpers.Success(c, "get tables success", allTables)
	}
}

// Get Table By Id
func GetTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		tableId := c.Param("table_id")
		var table models.Table

		err := tableCollection.FindOne(ctx, bson.M{"table_id": tableId}).Decode(&table)
		if err != nil {
			helpers.Error(c, "get table by id failed", 0, err)
			return
		}
		helpers.Success(c, "get table by id success", table)
	}
}

// Create Table
func CreateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var table models.Table

		if err := c.BindJSON(&table); err != nil {
			helpers.Error(c, "JSON Bind error", 400, err)
			return
		}

		validationErr := validate.Struct(table)

		if validationErr != nil {
			helpers.Error(c, "validation error", 400, validationErr)
			return
		}

		table.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		table.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		table.ID = primitive.NewObjectID()
		table.Table_id = table.ID.Hex()

		result, insertErr := tableCollection.InsertOne(ctx, table)

		if insertErr != nil {
			msg := fmt.Sprintf("table insert error: %s", insertErr)
			helpers.Error(c, msg, 0, insertErr)
			return
		}

		helpers.Success(c, "table create success", result)
	}
}

// Update Table Controller
func UpdateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		// var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		// defer cancel()

		// var table models.Table

		// if err := c.BindJSON(&table); err != nil {

		// }
	}
}
