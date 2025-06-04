package limiter

import "context"

type RateLimitRequest struct {
	Key      string
	Rate     int // tokens per interval
	Interval int // in seconds
}

type RateLimitResponse struct {
	Allowed       bool
	RetryAfterSec int
}

type Limiter interface {
	Allow(ctx context.Context, req RateLimitRequest) RateLimitResponse
}
