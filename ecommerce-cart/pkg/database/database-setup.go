package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database setup
func DBSet() *mongo.Client {
	ctx, channel := context.WithTimeout(context.Background(), 20*time.Second)
	defer channel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017/ecommerce_cart"))

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Println("failed to connect to mongodb")
		return nil
	}
	fmt.Println(":::: successfully connected to mongodb ::::")
	return client
}

// Database Client
var Client *mongo.Client = DBSet()

// User Data Collection
func UserData(client *mongo.Client, CollectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("ecommerce_cart").Collection(CollectionName)
	return collection
}

// Product Data Collection
func ProductData(client *mongo.Client, CollectionName string) *mongo.Collection {
	var productCollection *mongo.Collection = client.Database("ecommerce_cart").Collection(CollectionName)
	return productCollection
}
