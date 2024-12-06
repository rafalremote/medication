package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	DBDriver   string
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatalf("Invalid DB_PORT value: %v", err)
	}

	config := &Config{
		Port:       GetEnv("PORT", "8080"),
		DBDriver:   GetEnv("DB_DRIVER", "postgres"),
		DBHost:     GetEnv("DB_HOST", "localhost"),
		DBPort:     dbPort,
		DBUser:     GetEnv("DB_USER", "postgres"),
		DBPassword: GetEnv("DB_PASSWORD", "postgres"),
		DBName:     GetEnv("DB_NAME", "medication"),
		JWTSecret:  GetEnv("JWT_SECRET", "defaultsecret"),
	}

	return config, nil
}

func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
