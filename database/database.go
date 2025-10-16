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
	"time"

	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) {
	var db *gorm.DB
	var err error
	
	// Coba konek dengan retry mechanism
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		db, err = config.NewDatabaseConnection(cfg)
		if err == nil {
			break
		}
		
		log.Printf("⚠️ Attempt %d/%d: Failed to connect to database: %v", i+1, maxRetries, err)
		
		if i < maxRetries-1 {
			waitTime := time.Duration(i+1) * 2 * time.Second
			log.Printf("⏳ Retrying in %v...", waitTime)
			time.Sleep(waitTime)
		}
	}
	
	if err != nil {
		log.Fatalf("❌ Failed to connect to database after %d attempts: %v", maxRetries, err)
	}

	// Test connection
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("❌ Failed to get database instance: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("❌ Database ping failed: %v", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("✅ Database connection tested successfully")

	// Auto migrate
	err = db.AutoMigrate(
		&models.User{},
		&models.About{},
		&models.Portfolio{},
		&models.Skill{},
		&models.Qualification{},
	)
	if err != nil {
		log.Fatalf("❌ Failed to migrate database: %v", err)
	}

	DB = db
	log.Println("✅ Database migrated successfully")
}

func GetDB() *gorm.DB {
	return DB
}