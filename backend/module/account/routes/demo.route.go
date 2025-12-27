package routes

import (
	"github.com/bricksocoolxd/bengi-investment-system/module/account/controller"
	"github.com/bricksocoolxd/bengi-investment-system/module/account/repository"
	"github.com/bricksocoolxd/bengi-investment-system/module/account/service"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterDemoRoutes(app *fiber.App) {
	repo := repository.NewAccountRepository()
	demoService := service.NewDemoService(repo)
	ctrl := controller.NewDemoController(demoService)

	// Demo routes - all require authentication
	demo := app.Group("/api/v1/demo", middleware.AuthRequired())

	// Get or create demo account
	demo.Get("/", ctrl.GetOrCreateDemo)

	// Create new demo account
	demo.Post("/", ctrl.CreateDemo)

	// Deposit virtual funds
	demo.Post("/:accountId/deposit", ctrl.Deposit)

	// Reset demo account
	demo.Post("/:accountId/reset", ctrl.Reset)

	// Get demo stats
	demo.Get("/:accountId/stats", ctrl.GetStats)
}
