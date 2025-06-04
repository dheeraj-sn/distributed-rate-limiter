package limiter

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
	Allow(req RateLimitRequest) RateLimitResponse
}
