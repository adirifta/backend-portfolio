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
		DBName:     getEnv("DB_NAME", "portfolio_db"),
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
	if os.Getenv("K_SERVICE") != "" && cfg.DBInstanceName != "" {
		// Jalankan di Cloud Run (Gunakan Unix Socket)
		log.Println("Connecting to Cloud SQL via Unix socket...")
		dsn = fmt.Sprintf(
			"user=%s password=%s dbname=%s host=/cloudsql/%s sslmode=disable",
			cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBInstanceName, // <-- Gunakan variabel
		)
	} else {
		// Jalankan secara lokal (misalnya ke Supabase atau DB lokal)
		log.Println("Connecting to remote/local database (e.g., Supabase)...")
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=Asia/Jakarta",
			cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
		)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// Beri pesan error yang lebih spesifik
		return nil, fmt.Errorf("failed to connect database using DSN: %w", err)
	}

	log.Println("✅ Database connection established successfully")
	return db, nil
}

// Atau jika ingin lebih simple, gunakan ini saja:
// func InitDatabase(cfg *Config) (*gorm.DB, error) {
// 	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require",
// 		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

// 	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
// }
