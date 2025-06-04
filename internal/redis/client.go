package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(url string) *RedisClient {
	opt, _ := redis.ParseURL(url)
	rdb := redis.NewClient(opt)

	return &RedisClient{client: rdb}
}

func (r *RedisClient) Get(key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisClient) Set(key string, val string, ttlSeconds int) error {
	return r.client.Set(ctx, key, val, time.Duration(ttlSeconds)*time.Second).Err()
}

func (r *RedisClient) Eval(script string, keys []string, args ...interface{}) (interface{}, error) {
	return r.client.Eval(ctx, script, keys, args...).Result()
}
