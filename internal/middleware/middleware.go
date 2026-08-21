package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

func SetupMiddleware(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))
	app.Use(limiter.New(limiter.Config{
		Max:               100,
		Expiration:        1 * time.Minute,
		LimiterMiddleware: limiter.SlidingWindow{},
	}))
	app.Use(func(c fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)

		log.Printf("%s %s %d %s",
			c.Method(),
			c.Path(),
			c.Response().StatusCode(),
			duration,
		)
		return err
	})
}
