package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Environment represents the application environment.
type Environment string

const (
	EnvDevelopment Environment = "development"
	EnvProduction  Environment = "production"
	EnvTesting     Environment = "testing"
)

// Config holds all application configuration.
// Loaded from environment variables on startup.
type Config struct {
	// Server settings
	Port string
	Env  Environment

	// MongoDB connection
	MongoURI      string
	MongoDatabase string

	// Market data APIs
	TwelveDataAPIKey string
	FinnhubAPIKey    string

	// JWT authentication
	JWTSecret         string
	JWTExpireDuration time.Duration

	// Redis cache (optional)
	RedisURI string

	// Kafka messaging (optional)
	KafkaBrokers string
	KafkaGroupID string
}

// AppConfig is the global configuration instance.
var AppConfig *Config

// LoadConfig reads configuration from environment variables.
// Falls back to defaults if .env file is not found.
func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	AppConfig = &Config{
		Port: getEnv("PORT", "8080"),
		Env:  Environment(getEnv("ENV", string(EnvDevelopment))),

		MongoURI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDatabase: getEnv("MONGO_DATABASE", "bengi-investment"),

		TwelveDataAPIKey: getEnv("TWELVEDATA_API_KEY", ""),
		FinnhubAPIKey:    getEnv("FINNHUB_API_KEY", ""),

		JWTSecret:         getEnv("JWT_SECRET", "change-this-in-production"),
		JWTExpireDuration: parseDuration(getEnv("JWT_EXPIRE", "24h")),

		RedisURI: getEnv("REDIS_URI", "redis://localhost:6379"),

		KafkaBrokers: getEnv("KAFKA_BROKERS", "localhost:9092"),
		KafkaGroupID: getEnv("KAFKA_GROUP_ID", "bengi-investment"),
	}
}

// IsDevelopment returns true if running in development mode.
func (c *Config) IsDevelopment() bool {
	return c.Env == EnvDevelopment
}

// IsProduction returns true if running in production mode.
func (c *Config) IsProduction() bool {
	return c.Env == EnvProduction
}

// getEnv returns the environment variable value or a default.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// parseDuration parses a duration string, defaults to 24h on error.
func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 24 * time.Hour
	}
	return d
}
