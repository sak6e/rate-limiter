package middleware

import (
	"net"
	"net/http"

	"saksham.com/rate-limiter/limiter"
)

func LimiterMiddleware(lh *limiter.LimiterHead, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "Invalid IP", http.StatusBadRequest)
			return
		}

		if !lh.Allow(ip) {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
