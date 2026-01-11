package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"

	"backend/db"
	"backend/routes"

	"github.com/joho/godotenv"
)

// @title Expense Management System API
// @version 1.0
// @description Expense submission, approval, and payment processing system
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email tonny.wijayajuly2001@gmail.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey CookieAuth
// @in cookie
// @name access_token
func main() {
	if _, ok := os.LookupEnv("DATABASE_URL"); !ok {
		envPath := filepath.Join("..", ".env") // adjust relative path
		if err := godotenv.Load(envPath); err != nil {
			log.Println("No local .env file found at", envPath)
		} else {
			log.Println("Loaded local .env from", envPath)
		}
	}

	db.Connect()

	router := gin.Default()

	// CORS configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{os.Getenv("NUXT_URL")},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Cookie"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	}))

	// API routes
	apiGroup := router.Group("/api")
	routes.AuthRoutes(apiGroup)

	// Serve frontend static files
	router.Use(static.Serve("/", static.LocalFile("./dist", true)))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Starting server at", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
