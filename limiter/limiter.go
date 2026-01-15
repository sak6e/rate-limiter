package limiter

import (
	"sync"
	"time"
)

const (
	tokenCount = 10
	period     = time.Minute
)

type LimiterHead struct {
	limiters map[string]*Limiter
	mu       sync.Mutex
}

type Limiter struct {
	tokens         int
	lastRefillTime time.Time
}

func NewRateLimiter() *LimiterHead {
	return &LimiterHead{
		limiters: make(map[string]*Limiter),
	}
}

func (lh *LimiterHead) Allow(ip string) bool {
	lh.mu.Lock()
	defer lh.mu.Unlock()
	rl, exists := lh.limiters[ip]
	if !exists {
		lh.limiters[ip] = &Limiter{
			tokens:         tokenCount - 1,
			lastRefillTime: time.Now(),
		}
		return true
	}
	rl.RefillTokens()
	if rl.tokens > 0 {
		return true
	}
	return false
}

func (rl *Limiter) RefillTokens() {
	duration := time.Now().Sub(rl.lastRefillTime)
	if duration >= period {
		rl.tokens = tokenCount
		rl.lastRefillTime = time.Now()
	}
}
