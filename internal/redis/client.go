package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(ctx context.Context, url string) (*RedisClient, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		log.Fatalf("Failed to parse Redis URL: %v", err)
	}
	rdb := redis.NewClient(opt)

	// Wait until Redis is reachable (5 retries)
	for i := 0; i < 5; i++ {
		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()

		_, err := rdb.Ping(ctx).Result()
		if err == nil {
			return &RedisClient{client: rdb}, nil
		}

		time.Sleep(2 * time.Second)
	}
	return nil, fmt.Errorf("could not connect to Redis at %s", url)
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisClient) Set(ctx context.Context, key string, val string, ttlSeconds int) error {
	return r.client.Set(ctx, key, val, time.Duration(ttlSeconds)*time.Second).Err()
}

func (r *RedisClient) Eval(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error) {
	return r.client.Eval(ctx, script, keys, args...).Result()
}
