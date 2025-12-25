package routes

import (
	"github.com/bricksocoolxd/bengi-investment-system/module/auth/controller"
	"github.com/bricksocoolxd/bengi-investment-system/module/auth/repository"
	"github.com/bricksocoolxd/bengi-investment-system/module/auth/service"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {

	userRepo := repository.NewUserRepository()
	authService := service.NewAuthService(userRepo)
	ctrl := controller.NewAuthController(authService)

	auth := app.Group("/api/v1/auth")

	// Public routes
	auth.Post("/register", ctrl.Register)
	auth.Post("/login", ctrl.Login)
	auth.Post("/refresh", ctrl.RefreshToken)
	auth.Post("/logout", ctrl.Logout)

	// Protected routes
	auth.Get("/profile", middleware.AuthRequired(), ctrl.GetProfile)
	auth.Put("/password", middleware.AuthRequired(), ctrl.ChangePassword)
	auth.Post("/logout-all", middleware.AuthRequired(), ctrl.LogoutAll)
}
