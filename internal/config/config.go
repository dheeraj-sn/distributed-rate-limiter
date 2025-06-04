package config

import (
	"os"
)

type Config struct {
	RedisURL string
	HTTPPort string
}

func Load() *Config {
	return &Config{
		RedisURL: getEnv("REDIS_URL", "redis://localhost:6379"),
		HTTPPort: getEnv("HTTP_PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
