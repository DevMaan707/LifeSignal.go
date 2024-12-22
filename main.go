package main

import (
	"life-signal/database"
	"life-signal/routes"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	client, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := database.DisconnectDB(client); err != nil {
			log.Printf("Error disconnecting database: %v", err)
		}
	}()

	router := gin.Default()

	routes.Routes(router, client)

	router.Run()
}
