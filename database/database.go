package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Client

func ConnectDB() {
	clientOptions := options.Client().ApplyURI("mongodb://")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	fmt.Println("Connected to MongoDB!")
	DB = client
}


func GetCollection(databaseName, collectionName string) *mongo.Collection {
	return DB.Database(databaseName).Collection(collectionName)
}

func DisconnectDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := DB.Disconnect(ctx); err != nil {
		log.Fatalf("Error disconnecting from MongoDB: %v", err)
	}
	fmt.Println("Disconnected from MongoDB")
}
