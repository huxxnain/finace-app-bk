package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoDBURI     string
	DatabaseName   string
	JWTSecret      string
	JWTExpiryHours int
	Port           string
}

func LoadConfig() *Config {
	// Load .env file if it exists
	_ = godotenv.Load()

	jwtExpiryHours, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))

	return &Config{
		MongoDBURI:     getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		DatabaseName:   getEnv("DATABASE_NAME", "finance_app"),
		JWTSecret:      getEnv("JWT_SECRET", "your-super-secret-key"),
		JWTExpiryHours: jwtExpiryHours,
		Port:           getEnv("PORT", "3000"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
