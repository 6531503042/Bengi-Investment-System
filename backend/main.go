package main

import (
	"log"

	accountRoutes "github.com/bricksocoolxd/bengi-investment-system/module/account/routes"
	authRoutes "github.com/bricksocoolxd/bengi-investment-system/module/auth/routes"
	instrumentRoutes "github.com/bricksocoolxd/bengi-investment-system/module/instrument/routes"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/config"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/core/database"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/seeder"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	config.LoadConfig()
	database.ConnextMongoDB()

	// Run seeders (create default roles, etc.)
	seeder.RunSeeders()

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

	// Register modules
	authRoutes.RegisterRoutes(app)
	accountRoutes.RegisterRoutes(app)
	instrumentRoutes.RegisterRoutes(app)

	// Start server
	log.Printf("🚀 Server starting on port %s", config.AppConfig.Port)
	log.Fatal(app.Listen(":" + config.AppConfig.Port))
}
