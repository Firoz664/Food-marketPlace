package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// var MONGODB_URL = "mongodb+srv://shams:Test123@cluster0.yu3acw9.mongodb.net/golang-ecomm"
var MONGODB_URL = ("mongodb+srv://shams:Test123@cluster0.yu3acw9.mongodb.net/golang-ecomm")

// DBInstance creates and returns a new MongoDB client instance.
func DBInstance() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(MONGODB_URL))
	if err != nil {
		log.Fatalf("Failed to create new MongoDB client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // Reduced timeout to a more practical value
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil) // Changed from context.TODO() to use the same context as the connect method
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	fmt.Println("Connected to MongoDB successfully!")
	return client
}

var Client *mongo.Client = DBInstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("food-marketPlace").Collection(collectionName)

	return collection
}
