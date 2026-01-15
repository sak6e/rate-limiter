package old

import (
	"fmt"
	"sync"
	"time"
)

type RateLimiter struct {
	mu      sync.Mutex
	buckets map[string]int
	limit   int
}

func NewRateLimiter() *RateLimiter {
	rl := &RateLimiter{
		buckets: make(map[string]int),
		limit:   5,
	}
	go rl.refillLoop()
	return rl
}

func (rl *RateLimiter) refillLoop() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		for ip := range rl.buckets {
			rl.buckets[ip] = rl.limit
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) IsRateLimited(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if _, ok := rl.buckets[ip]; !ok {
		rl.buckets[ip] = rl.limit
	}

	if rl.buckets[ip] <= 0 {
		return true
	}

	rl.buckets[ip]--
	return false
}

func main() {
	rl := NewRateLimiter()

	// receive new req
	ip := "1.1.1.1"
	for i := 0; i < 10; i++ {
		if rl.IsRateLimited(ip) {
			fmt.Print("request is rate limited\n")
		} else {
			fmt.Print("request succeeded\n")
		}
	}
}
