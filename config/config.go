package config

import (
	"os"
	"strings"
)

type Config struct {
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	DBInstanceName   string
	JWTSecret        string
	JWTRefreshSecret string
	Port             string

	// Cookie settings
	CookieDomain   string
	CookieSecure   bool
	CookieSameSite string
	AllowedOrigins []string

	// Token expiry (in minutes)
	AccessTokenExpiry  int
	RefreshTokenExpiry int
}

func LoadConfig() *Config {
	secure := getEnv("COOKIE_SECURE", "true") == "true"

	// Default allowed origins (production domains + local dev)
	defaultOrigins := []string{
		"https://adirdk.cloud",
		"https://adirdk.com",
		"https://dashboard.adirdk.com",
		"http://localhost:3000",
		"http://localhost:8080",
		"https://www.adirdk.cloud",
		"https://www.adirdk.com",
		"https://www.dashboard.adirdk.com",
	}

	// Override via ALLOWED_ORIGINS env var (comma-separated)
	allowedOrigins := defaultOrigins
	if envOrigins := os.Getenv("ALLOWED_ORIGINS"); envOrigins != "" {
		allowedOrigins = strings.Split(envOrigins, ",")
		for i := range allowedOrigins {
			allowedOrigins[i] = strings.TrimSpace(allowedOrigins[i])
		}
	}

	return &Config{
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "5432"),
		DBUser:             getEnv("DB_USER", "postgres"),
		DBPassword:         getEnv("DB_PASSWORD", ""),
		DBName:             getEnv("DB_NAME", "portfolio-db"),
		DBInstanceName:     getEnv("DB_INSTANCE_NAME", ""),
		JWTSecret:          getEnv("JWT_SECRET", "your-secret-key"),
		JWTRefreshSecret:   getEnv("JWT_REFRESH_SECRET", "your-refresh-secret-key"),
		Port:               getEnv("PORT", "8080"),
		CookieDomain:       getEnv("COOKIE_DOMAIN", ""),
		CookieSecure:       secure,
		CookieSameSite:     getEnv("COOKIE_SAMESITE", "Strict"),
		AllowedOrigins:     allowedOrigins,
		AccessTokenExpiry:  15,    // 15 minutes
		RefreshTokenExpiry: 10080, // 7 days
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
