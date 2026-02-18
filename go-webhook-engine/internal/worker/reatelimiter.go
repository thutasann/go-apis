package worker

import (
	"context"
	"time"
)

// RateLimiter implements a simple token bucket
type RateLimiter struct {
	tokens chan struct{}
	ticker *time.Ticker
}

// NewRateLimiter creates a limiter with `rate` tokens per second
func NewRateLimiter(rate int) *RateLimiter {
	rl := &RateLimiter{
		tokens: make(chan struct{}, rate), // bucker size = rate
		ticker: time.NewTicker(time.Second / time.Duration(rate)),
	}

	// Fill tokens continuously
	go func() {
		for t := range rl.ticker.C {
			select {
			case rl.tokens <- struct{}{}:
			default: // bucket full
			}
			_ = t // avoid unused var
		}
	}()

	return rl
}

// Wait blocks until a token is available or context is canceled
func (rl *RateLimiter) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-rl.tokens:
		return nil
	}
}

// Stop the ticker
func (rl *RateLimiter) Stop() {
	rl.ticker.Stop()
	close(rl.tokens)
}
