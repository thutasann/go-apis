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
	"go.mongodb.org/mongo-driver/mongo"
)

// Order Mongo Collection
var orderCollection *mongo.Collection = database.OpenCollection(database.Client, "order")

// Get all Orders
func GetOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		result, err := orderCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			helpers.Error(c, "error occured while listing the orders list", 0, err)
		}

		var allOrders []bson.M
		if err = result.All(ctx, &allOrders); err != nil {
			log.Fatal(err)
		}
		helpers.Success(c, "get orders success", allOrders)
	}
}

// Get Order by ID
func GetOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		orderId := c.Param("order_id")
		var order models.Order

		err := orderCollection.FindOne(ctx, bson.M{"order_id": orderId}).Decode(&order)
		defer cancel()

		if err != nil {
			helpers.Error(c, "Error occured while fetching order by Id", 0, err)
		}
		helpers.Success(c, "Order fetched successfully", order)
	}
}

// Create Order
func CreateOrder() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

func UpdateOrder() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}
