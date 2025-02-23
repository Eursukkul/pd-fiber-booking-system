package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	JWTSecret string
	APIKey    string
}

func LoadConfig() (*Config, error) {
	// load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using default configurations")
	}

	return &Config{
		Port:      getEnv("PORT", ":3001"),
		JWTSecret: getEnv("JWT_SECRET", "your_default_jwt_secret"),
		APIKey:    getEnv("API_KEY", "your_default_api_key"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}