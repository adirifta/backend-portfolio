package config

import (
	"os"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
	Port       string
}

func LoadConfig() *Config {
	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "portfolio_db"),
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

	// Jika dijalankan di Cloud Run, gunakan koneksi via Cloud SQL Proxy
	if os.Getenv("K_SERVICE") != "" {
		// Gunakan koneksi via Cloud SQL Unix socket
		dsn = fmt.Sprintf(
			"user=%s password=%s dbname=%s host=/cloudsql/project-adi-474909:us-central1:portfolio-db sslmode=disable",
			cfg.DBUser, cfg.DBPassword, cfg.DBName,
		)
	} else {
		// Jalankan secara lokal
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=Asia/Jakarta",
			cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
		)
	}

	log.Printf("Connecting to database with DSN: %s", dsn)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}


// Atau jika ingin lebih simple, gunakan ini saja:
func InitDatabase(cfg *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)
	
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}