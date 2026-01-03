package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Redis RedisConfig
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

// LoadConfig loads the configuration from environment variables with fallback to default values
func LoadConfig() *Config {
	// Try to load dev.env file if it exists
	if err := loadEnvFile(); err != nil {
		log.Printf("Note: dev.env file not loaded: %v", err)
	}
	return &Config{
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
		},
	}
}

// Helper function to get environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Helper function to get environment variable as integer
func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// Helper function to get environment variable as boolean
func getEnvAsBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

// loadEnvFile attempts to load environment variables from dev.env file
func loadEnvFile() error {
	// List of possible env file locations, in order of preference
	envFiles := []string{
		"dev.env",           // Current directory
		"config/dev.env",    // Config subdirectory
		"../config/dev.env", // One level up (for tests)
		filepath.Join(os.Getenv("HOME"), "web-scrapper/config/dev.env"), // Home directory
	}

	// Try each possible location
	for _, file := range envFiles {
		if err := godotenv.Load(file); err == nil {
			log.Printf("Loaded environment from %s", file)
			return nil
		}
	}

	return godotenv.Load("dev.env") // Return error from default location
}
