package controller

import (
	"errors"

	"github.com/bricksocoolxd/bengi-investment-system/module/account/dto"
	"github.com/bricksocoolxd/bengi-investment-system/module/account/service"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/common"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

type DemoController struct {
	service *service.DemoService
}

func NewDemoController(service *service.DemoService) *DemoController {
	return &DemoController{
		service: service,
	}
}

// GetOrCreateDemo gets or creates a demo account
func (c *DemoController) GetOrCreateDemo(ctx *fiber.Ctx) error {
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return common.Unauthorized(ctx, "User not authenticated")
	}

	stats, err := c.service.GetOrCreateDemoAccount(ctx.Context(), userID)
	if err != nil {
		return common.InternalError(ctx, err.Error())
	}

	return common.Success(ctx, stats, "Demo account retrieved")
}

// CreateDemo creates a new demo account
func (c *DemoController) CreateDemo(ctx *fiber.Ctx) error {
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return common.Unauthorized(ctx, "User not authenticated")
	}

	var req dto.CreateDemoAccountRequest
	if err := ctx.BodyParser(&req); err != nil {
		return common.BadRequest(ctx, "Invalid request body")
	}

	result, err := c.service.CreateDemoAccount(ctx.Context(), userID, &req)
	if err != nil {
		return common.InternalError(ctx, err.Error())
	}

	return common.Created(ctx, result, result.Message)
}

// Deposit adds virtual funds to demo account
func (c *DemoController) Deposit(ctx *fiber.Ctx) error {
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return common.Unauthorized(ctx, "User not authenticated")
	}
	accountID := ctx.Params("accountId")

	var req dto.DemoDepositRequest
	if err := ctx.BodyParser(&req); err != nil {
		return common.BadRequest(ctx, "Invalid request body")
	}

	result, err := c.service.DemoDeposit(ctx.Context(), accountID, userID, &req)
	if err != nil {
		if err.Error() == "account not found" {
			return common.NotFound(ctx, err.Error())
		}
		if errors.Is(err, service.ErrNotDemoAccount) {
			return common.BadRequest(ctx, err.Error())
		}
		return common.InternalError(ctx, err.Error())
	}

	return common.Success(ctx, result, result.Message)
}

// Reset resets demo account to initial state
func (c *DemoController) Reset(ctx *fiber.Ctx) error {
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return common.Unauthorized(ctx, "User not authenticated")
	}
	accountID := ctx.Params("accountId")

	var req dto.DemoResetRequest
	if err := ctx.BodyParser(&req); err != nil {
		req = dto.DemoResetRequest{} // Use defaults if no body
	}

	result, err := c.service.DemoReset(ctx.Context(), accountID, userID, &req)
	if err != nil {
		if err.Error() == "account not found" {
			return common.NotFound(ctx, err.Error())
		}
		if errors.Is(err, service.ErrNotDemoAccount) {
			return common.BadRequest(ctx, err.Error())
		}
		return common.InternalError(ctx, err.Error())
	}

	return common.Success(ctx, result, result.Message)
}

// GetStats returns demo account statistics
func (c *DemoController) GetStats(ctx *fiber.Ctx) error {
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return common.Unauthorized(ctx, "User not authenticated")
	}
	accountID := ctx.Params("accountId")

	stats, err := c.service.GetDemoStats(ctx.Context(), accountID, userID)
	if err != nil {
		if err.Error() == "account not found" {
			return common.NotFound(ctx, err.Error())
		}
		if errors.Is(err, service.ErrNotDemoAccount) {
			return common.BadRequest(ctx, err.Error())
		}
		return common.InternalError(ctx, err.Error())
	}

	return common.Success(ctx, stats, "Demo account stats retrieved")
}
