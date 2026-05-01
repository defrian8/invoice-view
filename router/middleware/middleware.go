package middleware

import (
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func RateLimiter() fiber.Handler {
	max, err := strconv.Atoi(os.Getenv("LIMITER_INVOICE_MAX"))
	if err != nil {
		max = 60
	}

	duration, err := strconv.Atoi(os.Getenv("LIMITER_INVOICE_DURATION"))
	if err != nil {
		duration = 1
	}

	apiLimiter := limiter.New(limiter.Config{
		Max:        max,
		Expiration: time.Duration(duration) * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(429).JSON(fiber.Map{
				"error": "Too many requests",
			})
		},
		LimiterMiddleware: limiter.SlidingWindow{},
	})
	return apiLimiter
}

func APIKeyMiddleware() fiber.Handler {
	secretKey := os.Getenv("API_SECRET")

	return func(c *fiber.Ctx) error {
		apiKey := c.Get("X-API-Key")
		if apiKey == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "API key is missing",
			})
		}

		if apiKey != secretKey {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid API key",
			})
		}

		return c.Next()
	}
}
