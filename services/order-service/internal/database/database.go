package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB() *gorm.DB {
	// "postgres://root:password123@postgres-service:5432/yaak_kaii?sslmode=disable"
	pgHost := os.Getenv("PG_HOST")
	pgUser := os.Getenv("PG_USERNAME")
	pgPassword := os.Getenv("PG_PASSWORD")
	pgDatabase := os.Getenv("PG_DATABASE")
	pgPort := os.Getenv("PG_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		pgHost, pgUser, pgPassword, pgDatabase, pgPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get generic database object:", err)
	}

	sqlDB.SetMaxIdleConns(10)           // connection ที่เปิดค้างไว้
	sqlDB.SetMaxOpenConns(100)          // connection สูงสุดที่อนุญาต
	sqlDB.SetConnMaxLifetime(time.Hour) // ระยะเวลาสูงสุดของแต่ละ connection

	log.Println("Database connection successfully")
	return db
}
