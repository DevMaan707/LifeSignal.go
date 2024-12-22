package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Client

func ConnectDB() error {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		return fmt.Errorf("MONGO_URI environment variable is not set")
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	clientOptions := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	err = client.Database("admin").RunCommand(ctx, bson.D{{"ping", 1}}).Err()
	if err != nil {
		return fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	DB = client
	fmt.Println("Connected to MongoDB!")

	return nil
}

func GetCollection(databaseName, collectionName string) *mongo.Collection {
	return DB.Database(databaseName).Collection(collectionName)
}

func DisconnectDB() error {
	if DB == nil {
		return fmt.Errorf("no active database connection")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := DB.Disconnect(ctx)
	if err != nil {
		return fmt.Errorf("failed to disconnect from MongoDB: %v", err)
	}

	fmt.Println("Disconnected from MongoDB")
	return nil
}
