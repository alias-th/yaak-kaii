package main

import (
	"encoding/json"
	"log"
	"os"
	"yaak-kaii/services/product-service/internal/models"

	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type SeedFile struct {
	Category   models.Category            `json:"category"`
	Attributes []models.CategoryAttribute `json:"attributes"`
}

func main() {
	db, err := gorm.Open(postgres.Open("postgres://root:password123@localhost:5432/yaak_kaii?sslmode=disable"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	file, err := os.ReadFile("./categories.json")
	if err != nil {
		panic(err)
	}
	var seed SeedFile
	if err := json.Unmarshal(file, &seed); err != nil {
		panic(err)
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	category := models.Category{
		ID:          seed.Category.ID,
		Name:        seed.Category.Name,
		Description: seed.Category.Description,
	}

	if err := tx.Create(&models.Category{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}).Error; err != nil {
		tx.Rollback()
		panic(err)
	}

	for _, attr := range seed.Attributes {
		log.Println(attr)
		optionsJSON, _ := json.Marshal(attr.Options)

		categoryAttribute := models.CategoryAttribute{
			CategoryID: category.ID,
			Key:        attr.Key,
			Label:      attr.Label,
			Type:       attr.Type,
			Required:   attr.Required,
			Options:    datatypes.JSON(optionsJSON),
		}
		if err := tx.Create(&categoryAttribute).Error; err != nil {
			tx.Rollback()
			panic(err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		panic(err)
	}

	log.Println("Seeding completed successfully")
}
