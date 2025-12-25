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

	if validationErrors := utils.ParseAndValidate(c, &req); validationErrors != nil {
		return common.ValidationError(c, validationErrors)
	}

	result, err := ctrl.authService.Register(c.Context(), &req)
	if err != nil {
		if errors.Is(err, service.ErrEmailExists) {
			return common.BadRequest(c, "Email already exists")
		}
		return common.InternalError(c, err.Error())
	}

	return common.Created(c, result, "Register successfully")
}

func (ctrl *AuthController) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest

	if validationErrors := utils.ParseAndValidate(c, &req); validationErrors != nil {
		return common.ValidationError(c, validationErrors)
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

	refreshToken := c.Cookies(utils.RefreshTokenCookie)
	if refreshToken == "" {
		return common.Unauthorized(c, "Refresh token not found")
	}

	result, err := ctrl.authService.RefreshToken(c.Context(), refreshToken)
	if err != nil {
		utils.ClearAuthCookies(c)
		return common.Unauthorized(c, "Invalid refresh token")
	}

	utils.SetAuthCookies(c, result.AccessToken, result.RefreshToken)
	return common.Success(c, nil, "Token refreshed")
}

func (ctrl *AuthController) Logout(c *fiber.Ctx) error {
	// Delete session from Redis
	sessionID := c.Get("X-Session-ID")
	if sessionID != "" {
		_ = ctrl.authService.Logout(c.Context(), sessionID)
	}

	utils.ClearAuthCookies(c)

	return common.Success(c, nil, "Logout successful")
}

// LogoutAll logs out from all devices
func (ctrl *AuthController) LogoutAll(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return common.Unauthorized(c, "User not authenticated")
	}

	if err := ctrl.authService.LogoutAll(c.Context(), userID); err != nil {
		return common.InternalError(c, err.Error())
	}

	utils.ClearAuthCookies(c)

	return common.Success(c, nil, "Logged out from all devices")
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

func (ctrl *AuthController) ChangePassword(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return common.Unauthorized(c, "User not authenticated")
	}

	var req dto.ChangePasswordRequest
	if validationErrors := utils.ParseAndValidate(c, &req); validationErrors != nil {
		return common.ValidationError(c, validationErrors)
	}

	err := ctrl.authService.ChangePassword(c.Context(), userID, &req)
	if err != nil {
		if errors.Is(err, service.ErrWrongPassword) {
			return common.BadRequest(c, "Current password is incorrect")
		}
		return common.InternalError(c, err.Error())
	}

	return common.Success(c, nil, "Password changed successfully")
}
