package app

import (
	"log"
	"net/http"

	"github.com/recursionexcursion/dd-go-api/internal/api"
	"golang.org/x/time/rate"
)

func LoggerMW(next api.HandlerFn) api.HandlerFn {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Logger MW")
		next(w, r)
	}
}

func AuthMW(next api.HandlerFn) api.HandlerFn {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get(("Authorization"))
		log.Printf("Auth MW token is '%v'", token)
		next(w, r)
	}
}

// refil rate 5/sec, total bucket size is 10
var limiter = rate.NewLimiter(5, 10)

func RateLimitMW(next api.HandlerFn) api.HandlerFn {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("RL MW")
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next(w, r)
	}
}
