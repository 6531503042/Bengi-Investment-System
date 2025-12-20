package routes

import (
	"github.com/bricksocoolxd/bengi-investment-system/module/portfolio/controller"
	"github.com/bricksocoolxd/bengi-investment-system/module/portfolio/repository"
	"github.com/bricksocoolxd/bengi-investment-system/module/portfolio/service"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	// Wire up dependencies
	repo := repository.NewPortfolioRepository()
	portfolioSvc := service.NewPortfolioService(repo)
	ctrl := controller.NewPortfolioController(portfolioSvc)

	// All routes are protected
	portfolios := app.Group("/api/v1/portfolios", middleware.AuthRequired())

	portfolios.Post("/", ctrl.CreatePortfolio)
	portfolios.Get("/", ctrl.GetPortfolios)
	portfolios.Get("/:id", ctrl.GetPortfolioByID)
	portfolios.Get("/:id/summary", ctrl.GetPortfolioSummary)
	portfolios.Put("/:id", ctrl.UpdatePortfolio)
	portfolios.Delete("/:id", ctrl.DeletePortfolio)
	portfolios.Get("/:id/positions", ctrl.GetPositions)

	// Position routes
	positions := app.Group("/api/v1/positions", middleware.AuthRequired())
	positions.Get("/:id", ctrl.GetPositionDetail)
}
