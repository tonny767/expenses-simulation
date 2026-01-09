package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"

	"backend/db"
	"backend/routes"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("No .env file found, relying on system env")
	}

	db.Connect()

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{os.Getenv("NUXT_URL")}, // Your Nuxt dev URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Cookie"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // Allows cookies for auth
		MaxAge:           24 * time.Hour,
	}))

	apiGroup := router.Group("/api")
	routes.AuthRoutes(apiGroup)

	router.Use(static.Serve("/", static.LocalFile("./dist", true)))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	serverAddr := "localhost:" + port
	log.Println("Starting server at", serverAddr)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
