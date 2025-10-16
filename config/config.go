package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	DBInstanceName string
	JWTSecret      string
	Port           string
}

func LoadConfig() *Config {
	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "portfolio-db"),
		DBInstanceName: getEnv("DB_INSTANCE_NAME", ""),
		JWTSecret:  getEnv("JWT_SECRET", "your-secret-key"),
		Port:       getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// NewDatabaseConnection untuk Supabase dengan SSL (CARA YANG BENAR)
// func NewDatabaseConnection(cfg *Config) (*gorm.DB, error) {
// 	// Gunakan sslmode=require atau verify-full untuk Supabase
// 	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=Asia/Jakarta",
// 		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

// 	log.Printf("Connecting to database: %s@%s:%s", cfg.DBUser, cfg.DBHost, cfg.DBPort)

// 	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
// }

func NewDatabaseConnection(cfg *Config) (*gorm.DB, error) {
	var dsn string

	// Deteksi lingkungan Cloud Run
	if os.Getenv("K_SERVICE") != "" {
		log.Println("🚀 Running in Cloud Run environment")
		
		if cfg.DBInstanceName != "" {
			// Gunakan Unix socket untuk Cloud SQL
			log.Printf("🔗 Connecting to Cloud SQL: %s", cfg.DBInstanceName)
			dsn = fmt.Sprintf(
				"user=%s password=%s database=%s host=/cloudsql/%s",
				cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBInstanceName,
			)
		} else {
			// Fallback ke koneksi TCP biasa
			log.Printf("🔗 Connecting via TCP: %s:%s", cfg.DBHost, cfg.DBPort)
			dsn = fmt.Sprintf(
				"host=%s user=%s password=%s dbname=%s port=%s sslmode=require",
				cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
			)
		}
	} else {
		// Local development
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
	// Sederhana saja, jangan log password sebenarnya
	return "host=*** user=*** password=*** dbname=***"
}

// Atau jika ingin lebih simple, gunakan ini saja:
// func InitDatabase(cfg *Config) (*gorm.DB, error) {
// 	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require",
// 		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

// 	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
// }
