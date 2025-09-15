package database

import (
	"fmt"
	"log"
	"backend-portfolio/config"
	"backend-portfolio/models"
	"backend-portfolio/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate models
	err = db.AutoMigrate(
		&models.User{},
		&models.About{},
		&models.Portfolio{},
		&models.Skill{},
		&models.Qualification{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	DB = db
	log.Println("Database connected successfully")
	
	// Initialize JWT with secret from config
	utils.InitJWT(cfg)
}

func GetDB() *gorm.DB {
	return DB
}