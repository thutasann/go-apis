package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB Instance
func DBInstance() *mongo.Client {
	mongoUrl := "mongodb://localhost:27017/go_restaurant"
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUrl))

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

var Client *mongo.Client = DBInstance()

// OpenCollection returns the collection from the database
func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("go_restaurant").Collection(collectionName)
	return collection
}
