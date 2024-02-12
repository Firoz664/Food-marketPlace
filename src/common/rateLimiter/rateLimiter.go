package ratelimiter

import (
	"time"

	"github.com/gin-contrib/limiter"
	"github.com/gin-gonic/gin"
)

// Initialize the rate limiter
var rateLimiter *limiter.Limiter

func init() {
	rateLimiter = NewRateLimiter()
}

func NewRateLimiter() *limiter.Limiter {
	// Create a new rate limiter with a rate limit of 10 requests per minute
	rate := limiter.NewRateLimiter(limiter.IPKey, &limiter.Rate{
		Max:    10,
		Window: time.Minute,
	})

	return rate
}

func RateLimiterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Apply the rate limiter middleware to the current request
		if err := rateLimiter.Limit(c.Writer, c.Request); err != nil {
			// Return a 429 Too Many Requests error if the rate limit is exceeded
			c.AbortWithStatus(429)
			return
		}

		// Continue processing the request if the rate limit is not exceeded
		c.Next()
	}
}
