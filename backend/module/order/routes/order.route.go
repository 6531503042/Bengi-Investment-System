package routes

import (
	"github.com/bricksocoolxd/bengi-investment-system/module/order/controller"
	"github.com/bricksocoolxd/bengi-investment-system/module/order/repository"
	"github.com/bricksocoolxd/bengi-investment-system/module/order/service"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	// Wire up dependencies
	repo := repository.NewOrderRepository()
	orderSvc := service.NewOrderService(repo)
	ctrl := controller.NewOrderController(orderSvc)

	// All routes are protected
	orders := app.Group("/api/v1/orders", middleware.AuthRequired())

	orders.Post("/", ctrl.CreateOrder)
	orders.Get("/", ctrl.GetOrders)
	orders.Get("/:id", ctrl.GetOrderByID)
	orders.Post("/:id/cancel", ctrl.CancelOrder)
}
