package main

import (
	"log"

	"github.com/dheeraj-sn/distributed-rate-limiter/internal/config"
	"github.com/dheeraj-sn/distributed-rate-limiter/internal/http"
	"github.com/dheeraj-sn/distributed-rate-limiter/internal/limiter"
	"github.com/dheeraj-sn/distributed-rate-limiter/internal/redis"
)

func main() {
	cfg := config.Load()

	redisClient := redis.NewRedisClient(cfg.RedisURL)
	rateLimiter := limiter.NewTokenBucketLimiter(redisClient)

	server := http.NewServer(cfg.HTTPPort, rateLimiter)
	log.Printf("Starting server on port %s...", cfg.HTTPPort)
	if err := server.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
