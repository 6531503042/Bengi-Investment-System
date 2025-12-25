package middleware

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bricksocoolxd/bengi-investment-system/pkg/cache"
	"github.com/gofiber/fiber/v2"
)

// RateLimitConfig defines rate limiting configuration
type RateLimitConfig struct {
	// Max requests per window
	Max int
	// Window duration (e.g., 1 minute)
	Window time.Duration
	// Key generator function (default: IP-based)
	KeyGenerator func(*fiber.Ctx) string
	// Skip function to bypass rate limiting
	Skip func(*fiber.Ctx) bool
}

// DefaultRateLimitConfig returns default config (100 req/min)
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		Max:    100,
		Window: time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		Skip: nil,
	}
}

// AuthRateLimitConfig returns stricter config for auth endpoints (10 req/min)
func AuthRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		Max:    10,
		Window: time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return "auth:" + c.IP()
		},
		Skip: nil,
	}
}

// RateLimit creates a rate limiting middleware using Redis
func RateLimit(config ...RateLimitConfig) fiber.Handler {
	cfg := DefaultRateLimitConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	return func(c *fiber.Ctx) error {
		// Check if should skip
		if cfg.Skip != nil && cfg.Skip(c) {
			return c.Next()
		}

		// Check if Redis is available
		if !cache.IsConnected() {
			// If Redis is not available, allow request (fail open)
			return c.Next()
		}

		// Generate key
		key := "ratelimit:" + cfg.KeyGenerator(c)

		// Get current count
		count, err := cache.Incr(key)
		if err != nil {
			// If error, allow request (fail open)
			return c.Next()
		}

		// Set expiry on first request
		if count == 1 {
			_ = cache.Expire(key, cfg.Window)
		}

		// Get remaining TTL
		ttl, _ := cache.TTL(key)
		resetTime := time.Now().Add(ttl).Unix()

		// Set rate limit headers
		c.Set("X-RateLimit-Limit", strconv.Itoa(cfg.Max))
		c.Set("X-RateLimit-Remaining", strconv.Itoa(max(0, cfg.Max-int(count))))
		c.Set("X-RateLimit-Reset", strconv.FormatInt(resetTime, 10))

		// Check if over limit
		if int(count) > cfg.Max {
			c.Set("Retry-After", strconv.FormatInt(int64(ttl.Seconds()), 10))
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"message": fmt.Sprintf("Too many requests. Try again in %d seconds", int(ttl.Seconds())),
				"error":   "RATE_LIMIT_EXCEEDED",
			})
		}

		return c.Next()
	}
}

// RateLimitByUser creates a rate limiter that uses user ID instead of IP
func RateLimitByUser(max int, window time.Duration) fiber.Handler {
	return RateLimit(RateLimitConfig{
		Max:    max,
		Window: window,
		KeyGenerator: func(c *fiber.Ctx) string {
			userID := c.Locals("userId")
			if userID != nil {
				return "user:" + userID.(string)
			}
			return "ip:" + c.IP()
		},
	})
}

// max returns the larger of x or y
func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
