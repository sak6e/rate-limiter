package main

import (
	"fmt"
	"net/http"

	"saksham.com/rate-limiter/limiter"
	"saksham.com/rate-limiter/middleware"
)

func main() {
	rl := limiter.NewRateLimiter()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Request allowed")
	})

	http.Handle("/", middleware.LimiterMiddleware(rl, handler))

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}
