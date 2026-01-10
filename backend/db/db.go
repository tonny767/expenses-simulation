package db

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Build URL from parts if DATABASE_URL not set
		dbUser := os.Getenv("DB_USER")
		dbPassword := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		if dbUser == "" || dbName == "" {
			log.Fatal("DATABASE_URL or DB_USER/DB_NAME must be set")
		}
		if dbHost == "" {
			dbHost = "localhost"
		}
		if dbPort == "" {
			dbPort = "5432"
		}
		dbURL = "postgres://" + dbUser + ":" + dbPassword + "@" + dbHost + ":" + dbPort + "/" + dbName + "?sslmode=disable"
	}
	log.Println("Connecting to DB with URL:", dbURL)

	var err error
	DB, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC() // <-- store UTC in DB
		},
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Connected to PostgreSQL successfully")
}
