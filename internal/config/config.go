package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost            string
	DBPort            int
	DBUser            string
	DBPassword        string
	DBName            string
	JWTSecret         string
	Port              string
	ChirpStackHost    string
	ChirpStackPort    string
	ChirpStackToken   string
	ChirpStackEnabled bool
}

func Load() (*Config, error) {
	// Load .env file if it exists
	godotenv.Load()

	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		dbPort = 5432
	}

	chirpStackEnabled := getEnv("CHIRPSTACK_ENABLED", "true") == "true"

	return &Config{
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            dbPort,
		DBUser:            getEnv("DB_USER", "postgres"),
		DBPassword:        getEnv("DB_PASSWORD", "password123"),
		DBName:            getEnv("DB_NAME", "auth_db"),
		JWTSecret:         getEnv("JWT_SECRET", "your-secret-key"),
		Port:              getEnv("PORT", "8080"),
		ChirpStackHost:    getEnv("CHIRPSTACK_HOST", "192.168.0.21"),
		ChirpStackPort:    getEnv("CHIRPSTACK_PORT", "8090"),
		ChirpStackToken:   getEnv("CHIRPSTACK_TOKEN", ""),
		ChirpStackEnabled: chirpStackEnabled,
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
