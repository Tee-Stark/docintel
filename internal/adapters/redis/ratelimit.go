package redis

import "context"

type RateLimiter struct{}

func NewRateLimiter() *RateLimiter { return &RateLimiter{} }

func (r *RateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	panic("not implemented")
}
