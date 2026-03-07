package database

import (
	"backend-portfolio/config"
	"backend-portfolio/models"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// NewConnection creates a new database connection based on environment.
func NewConnection(cfg *config.Config) (*gorm.DB, error) {
	var dsn string

	if os.Getenv("K_SERVICE") != "" {
		log.Println("🚀 Running in Cloud Run environment")

		if cfg.DBInstanceName != "" {
			log.Printf("🔗 Connecting to Cloud SQL: %s", cfg.DBInstanceName)
			dsn = fmt.Sprintf(
				"user=%s password=%s database=%s host=/cloudsql/%s",
				cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBInstanceName,
			)
		} else {
			log.Printf("🔗 Connecting via TCP: %s:%s", cfg.DBHost, cfg.DBPort)
			dsn = fmt.Sprintf(
				"host=%s user=%s password=%s dbname=%s port=%s sslmode=require",
				cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
			)
		}
	} else {
		log.Println("💻 Running in local environment")
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
		)
	}

	log.Printf("📝 DSN: %s", maskPassword(dsn))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("✅ Database connection established successfully")
	return db, nil
}

func maskPassword(dsn string) string {
	return "host=*** user=*** password=*** dbname=***"
}

// InitDB initializes the database with retry mechanism, connection pooling, and auto-migration.
func InitDB(cfg *config.Config) {
	var db *gorm.DB
	var err error

	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		db, err = NewConnection(cfg)
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

	// Connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("✅ Database connection tested successfully")

	// Auto migrate
	err = db.AutoMigrate(
		&models.User{},
		&models.About{},
		&models.Portfolio{},
		&models.PortfolioMedia{},
		&models.Skill{},
		&models.Qualification{},
		&models.Visitor{},
	)
	if err != nil {
		log.Fatalf("❌ Failed to migrate database: %v", err)
	}

	DB = db
	log.Println("✅ Database migrated successfully")
}

// GetDB returns the global database instance.
func GetDB() *gorm.DB {
	return DB
}
