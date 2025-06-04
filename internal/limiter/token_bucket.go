package limiter

import (
	"fmt"
	"time"

	"github.com/dheeraj-sn/distributed-rate-limiter/internal/redis"
)

type TokenBucketLimiter struct {
	redis *redis.RedisClient
}

func NewTokenBucketLimiter(redisClient *redis.RedisClient) *TokenBucketLimiter {
	return &TokenBucketLimiter{redis: redisClient}
}

const tokenBucketScript = `
local key = KEYS[1]
local rate = tonumber(ARGV[1])
local interval = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
local ttl = interval

local data = redis.call("HMGET", key, "tokens", "last")
local tokens = tonumber(data[1]) or rate
local last = tonumber(data[2]) or now

local delta = math.max(0, now - last)
local refill = math.floor(delta * (rate / interval))
tokens = math.min(rate, tokens + refill)

local allowed = tokens > 0
if allowed then
	tokens = tokens - 1
end

redis.call("HMSET", key, "tokens", tokens, "last", now)
redis.call("EXPIRE", key, ttl)

return allowed
`

func (t *TokenBucketLimiter) Allow(req RateLimitRequest) RateLimitResponse {
	now := time.Now().Unix()

	result, err := t.redis.Eval(tokenBucketScript,
		[]string{req.Key},
		req.Rate, req.Interval, now,
	)
	if err != nil {
		fmt.Println("Redis error:", err)
		return RateLimitResponse{Allowed: true} // fallback
	}

	allowed, _ := result.(int64)
	return RateLimitResponse{
		Allowed:       allowed == 1,
		RetryAfterSec: 1,
	}
}
