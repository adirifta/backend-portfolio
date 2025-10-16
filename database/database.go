// package database

// import (
// 	"backend-portfolio/config"
// 	"backend-portfolio/models"
// 	"fmt"
// 	"log"

// 	"gorm.io/driver/postgres"
// 	"gorm.io/gorm"
// )

// var DB *gorm.DB

// func InitDB(cfg *config.Config) {
// 	// Gunakan sslmode=require untuk Supabase
// 	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=Asia/Jakarta",
// 		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

// 	log.Printf("Connecting to Supabase: %s", cfg.DBHost)

// 	db, err := gorm.Open(postgres.New(postgres.Config{
// 		DSN: dsn,
// 		PreferSimpleProtocol: true,
// 	}), &gorm.Config{})
// 	if err != nil {
// 		log.Fatal("Failed to connect to Supabase:", err)
// 	}

// 	// Auto migrate models
// 	err = db.AutoMigrate(
// 		&models.User{},
// 		&models.About{},
// 		&models.Portfolio{},
// 		&models.Skill{},
// 		&models.Qualification{},
// 	)
// 	if err != nil {
// 		log.Fatal("Failed to migrate database:", err)
// 	}

// 	DB = db
// 	log.Println("✅ Connected to Supabase successfully")
// }

// func GetDB() *gorm.DB {
// 	return DB
// }

package database

import (
	"backend-portfolio/config"
	"backend-portfolio/models"
	"log"

	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) {
	db, err := config.NewDatabaseConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.About{},
		&models.Portfolio{},
		&models.Skill{},
		&models.Qualification{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	DB = db
	log.Println("✅ Database migrated successfully")
}

func GetDB() *gorm.DB {
	return DB
}