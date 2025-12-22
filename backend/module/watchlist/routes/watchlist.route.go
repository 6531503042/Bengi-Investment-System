package routes

import (
	"github.com/bricksocoolxd/bengi-investment-system/module/watchlist/controller"
	"github.com/bricksocoolxd/bengi-investment-system/module/watchlist/repository"
	"github.com/bricksocoolxd/bengi-investment-system/module/watchlist/service"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	// Wire up dependencies
	repo := repository.NewWatchlistRepository()
	watchlistSvc := service.NewWatchlistService(repo)
	ctrl := controller.NewWatchlistController(watchlistSvc)

	// All routes are protected
	watchlists := app.Group("/api/v1/watchlists", middleware.AuthRequired())

	watchlists.Post("/", ctrl.CreateWatchlist)
	watchlists.Get("/", ctrl.GetWatchlists)
	watchlists.Get("/:id", ctrl.GetWatchlistByID)
	watchlists.Put("/:id", ctrl.UpdateWatchlist)
	watchlists.Delete("/:id", ctrl.DeleteWatchlist)
	watchlists.Post("/:id/symbols", ctrl.AddSymbol)
	watchlists.Delete("/:id/symbols/:symbol", ctrl.RemoveSymbol)
}
