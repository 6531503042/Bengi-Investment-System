package routes

import (
	"github.com/bricksocoolxd/bengi-investment-system/module/account/controller"
	"github.com/bricksocoolxd/bengi-investment-system/module/account/repository"
	"github.com/bricksocoolxd/bengi-investment-system/module/account/service"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	// Wire up dependencies
	repo := repository.NewAccountRepository()
	accountService := service.NewAccountService(repo)
	ctrl := controller.NewAccountController(accountService)

	// All routes are protected
	accounts := app.Group("/api/v1/accounts", middleware.AuthRequired())

	accounts.Post("/", ctrl.CreateAccount)
	accounts.Get("/", ctrl.GetAccounts)
	accounts.Get("/:id", ctrl.GetAccountByID)
	accounts.Post("/:id/deposit", ctrl.Deposit)
	accounts.Post("/:id/withdraw", ctrl.Withdraw)
	accounts.Get("/:id/transactions", ctrl.GetTransactions)
}
