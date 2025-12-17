package controller

import (
	"errors"

	"github.com/bricksocoolxd/bengi-investment-system/module/auth/dto"
	"github.com/bricksocoolxd/bengi-investment-system/module/auth/service"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/common"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/middleware"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (ctrl *AuthController) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return common.BadRequest(c, "Invalid request body")
	}

	result, err := ctrl.authService.Register(c.Context(), &req)
	if err != nil {
		if errors.Is(err, service.ErrEmailExists) {
			return common.BadRequest(c, "Email already exists")
		}
		return common.InternalError(c, err.Error())
	}

	return common.Success(c, result, "Register successfully")
}

func (ctrl *AuthController) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return common.BadRequest(c, "Invalid request body")
	}
	result, err := ctrl.authService.Login(c.Context(), &req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return common.Unauthorized(c, "Invalid email or password")
		}
		return common.InternalError(c, err.Error())
	}

	utils.SetAuthCookies(c, result.AccessToken, result.RefreshToken)
	return common.Success(c, fiber.Map{
		"user": result.User,
	}, "Login successful")
}

func (ctrl *AuthController) RefreshToken(c *fiber.Ctx) error {
	// Get refresh token from cookie
	refreshToken := c.Cookies(utils.RefreshTokenCookie)
	if refreshToken == "" {
		return common.Unauthorized(c, "Refresh token not found")
	}
	// Validate and generate new tokens
	result, err := ctrl.authService.RefreshToken(c.Context(), refreshToken)
	if err != nil {
		utils.ClearAuthCookies(c)
		return common.Unauthorized(c, "Invalid refresh token")
	}
	// Set new cookies
	utils.SetAuthCookies(c, result.AccessToken, result.RefreshToken)
	return common.Success(c, nil, "Token refreshed")
}

func (ctrl *AuthController) Logout(c *fiber.Ctx) error {

	utils.ClearAuthCookies(c)

	return common.Success(c, nil, "Logout successful")
}

func (ctrl *AuthController) GetProfile(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return common.Unauthorized(c, "User not authenticated")
	}

	result, err := ctrl.authService.GetProfile(c.Context(), userID)
	if err != nil {
		return common.NotFound(c, "User not found")
	}

	return common.Success(c, result, "Get profile successfully")
}
