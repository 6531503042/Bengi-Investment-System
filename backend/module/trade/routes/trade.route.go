package routes

import (
	accountRepo "github.com/bricksocoolxd/bengi-investment-system/module/account/repository"
	authModel "github.com/bricksocoolxd/bengi-investment-system/module/auth/model"
	orderRepo "github.com/bricksocoolxd/bengi-investment-system/module/order/repository"
	portfolioRepo "github.com/bricksocoolxd/bengi-investment-system/module/portfolio/repository"
	"github.com/bricksocoolxd/bengi-investment-system/module/trade/controller"
	"github.com/bricksocoolxd/bengi-investment-system/module/trade/repository"
	"github.com/bricksocoolxd/bengi-investment-system/module/trade/service"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	// Wire up dependencies
	tradeRepo := repository.NewTradeRepository()
	orderRepository := orderRepo.NewOrderRepository()
	accountRepository := accountRepo.NewAccountRepository()
	portfolioRepository := portfolioRepo.NewPortfolioRepository()

	tradeSvc := service.NewTradeService(
		tradeRepo,
		orderRepository,
		accountRepository,
		portfolioRepository,
	)
	ctrl := controller.NewTradeController(tradeSvc)

	// Protected routes
	trades := app.Group("/api/v1/trades", middleware.AuthRequired())

	trades.Get("/", ctrl.GetTrades)
	trades.Get("/summary", ctrl.GetTradeSummary)
	trades.Get("/:id", ctrl.GetTradeByID)

	// Admin only - manual trade execution
	trades.Post("/execute",
		middleware.RoleRequired(authModel.RoleAdmin),
		ctrl.ExecuteTrade,
	)

	// Order trades route
	orders := app.Group("/api/v1/orders", middleware.AuthRequired())
	orders.Get("/:id/trades", ctrl.GetTradesByOrderID)
}
