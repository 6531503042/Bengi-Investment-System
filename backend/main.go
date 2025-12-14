package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Bengi Investment System",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Bengi Investment System is running",
		})
	})

	// TODO: Register modules here
	// auth.RegisterRoutes(app)
	// account.RegisterRoutes(app)
	// portfolio.RegisterRoutes(app)
	// order.RegisterRoutes(app)
	// trade.RegisterRoutes(app)
	// instrument.RegisterRoutes(app)

	// Start server
	log.Fatal(app.Listen(":3000"))
}
