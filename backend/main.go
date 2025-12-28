package main

import (
	"context"
	"log"
	"time"

	accountRoutes "github.com/bricksocoolxd/bengi-investment-system/module/account/routes"
	authRoutes "github.com/bricksocoolxd/bengi-investment-system/module/auth/routes"
	"github.com/bricksocoolxd/bengi-investment-system/module/instrument/repository"
	instrumentRoutes "github.com/bricksocoolxd/bengi-investment-system/module/instrument/routes"
	"github.com/bricksocoolxd/bengi-investment-system/module/instrument/service"
	orderRoutes "github.com/bricksocoolxd/bengi-investment-system/module/order/routes"
	portfolioRoutes "github.com/bricksocoolxd/bengi-investment-system/module/portfolio/routes"
	tradeRoutes "github.com/bricksocoolxd/bengi-investment-system/module/trade/routes"
	watchlistRoutes "github.com/bricksocoolxd/bengi-investment-system/module/watchlist/routes"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/cache"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/config"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/core/database"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/seeder"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/ws"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	config.LoadConfig()
	database.ConnextMongoDB()

	// Initialize Redis (optional - continues if Redis is unavailable)
	if err := cache.Initialize(); err != nil {
		log.Printf("⚠️ Redis not available: %v (caching disabled)", err)
	}

	// Run seeders (create default roles, etc.)
	seeder.RunSeeders()

	// Start Symbol Sync Service (fetches all stocks/ETFs/crypto from Finnhub)
	instrumentRepo := repository.NewInstrumentRepository()
	symbolSyncService := service.NewSymbolSyncService(instrumentRepo)
	ctx := context.Background()

	// Start periodic sync (every 24 hours)
	go symbolSyncService.StartPeriodicSync(ctx, 24*time.Hour)
	log.Println("📈 Symbol Sync Service started (syncs all instruments from Finnhub)")

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
	accountRoutes.RegisterDemoRoutes(app) // Demo trading routes
	instrumentRoutes.RegisterRoutes(app)
	portfolioRoutes.RegisterRoutes(app)
	orderRoutes.RegisterRoutes(app)
	tradeRoutes.RegisterRoutes(app)
	watchlistRoutes.RegisterRoutes(app)

	// WebSocket routes
	ws.RegisterRoutes(app)

	// Start server
	log.Printf("🚀 Server starting on port %s", config.AppConfig.Port)
	log.Fatal(app.Listen(":" + config.AppConfig.Port))
}
