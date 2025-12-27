package routes

import (
	"github.com/bricksocoolxd/bengi-investment-system/module/auth/model"
	"github.com/bricksocoolxd/bengi-investment-system/module/instrument/controller"
	"github.com/bricksocoolxd/bengi-investment-system/module/instrument/repository"
	"github.com/bricksocoolxd/bengi-investment-system/module/instrument/service"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	// Wire up dependencies
	repo := repository.NewInstrumentRepository()
	marketSvc := service.NewMarketDataService()
	instrumentSvc := service.NewInstrumentService(repo, marketSvc)
	ctrl := controller.NewInstrumentController(instrumentSvc)

	instruments := app.Group("/api/v1/instruments")

	// Public routes (no auth required)
	instruments.Get("/", ctrl.GetInstruments)
	instruments.Get("/search", ctrl.SearchInstruments)
	instruments.Get("/:symbol", ctrl.GetInstrumentBySymbol)
	instruments.Get("/:symbol/quote", ctrl.GetQuote)
	instruments.Get("/:symbol/candles", ctrl.GetCandles)

	// Admin routes (auth + admin role required)
	admin := instruments.Group("", middleware.AuthRequired(), middleware.RoleRequired(model.RoleAdmin))
	admin.Post("/", ctrl.CreateInstrument)
	admin.Put("/:symbol", ctrl.UpdateInstrument)
}
