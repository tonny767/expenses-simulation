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
		log.Fatal("DATABASE_URL not set")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{
		NowFunc: func() time.Time {
			loc, _ := time.LoadLocation("Asia/Jakarta") // set time localize to ID/Jakarta
			return time.Now().In(loc)
		},
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
}
