package main

import (
	"fmt"
	"time"
)

var buckets = map[string]int{}

func main() {
	go func() {
		// refill buckets
		for {
			_ = time.AfterFunc(10*time.Second, refillBuckets)
		}
	}()

	// receive new req
	ip := "1.1.1.1"
	for i := 0; i < 10; i++ {
		req := isRateLimited(ip)
		if req {
			fmt.Print("request is rate limited\n")
			continue
		}
		fmt.Print("request succeeded\n")
	}
}

func refillBuckets() {
	for ip, _ := range buckets {
		buckets[ip] = 5
	}
}

func isRateLimited(ip string) bool {
	if _, ok := buckets[ip]; !ok {
		buckets[ip] = 5
	}
	if buckets[ip] <= 0 {
		return true
	}
	buckets[ip] = buckets[ip] - 1
	return false
}
