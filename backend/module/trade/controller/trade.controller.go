package controller

import (
	"errors"

	"github.com/bricksocoolxd/bengi-investment-system/module/trade/dto"
	"github.com/bricksocoolxd/bengi-investment-system/module/trade/service"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/common"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/middleware"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type TradeController struct {
	tradeService *service.TradeService
}

func NewTradeController(tradeService *service.TradeService) *TradeController {
	return &TradeController{
		tradeService: tradeService,
	}
}

// ExecuteTrade executes a trade (admin only - for manual execution)
// POST /api/v1/trades/execute
func (ctrl *TradeController) ExecuteTrade(c *fiber.Ctx) error {
	var req dto.ExecuteTradeRequest
	if validationErrors := utils.ParseAndValidate(c, &req); validationErrors != nil {
		return common.ValidationError(c, validationErrors)
	}

	result, err := ctrl.tradeService.ExecuteTrade(c.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrOrderNotFound):
			return common.NotFound(c, "Order not found")
		case errors.Is(err, service.ErrOrderAlreadyFilled):
			return common.BadRequest(c, "Order is already filled")
		case errors.Is(err, service.ErrOrderNotExecutable):
			return common.BadRequest(c, "Order cannot be executed")
		case errors.Is(err, service.ErrInsufficientBalance):
			return common.BadRequest(c, "Insufficient balance")
		case errors.Is(err, service.ErrInsufficientShares):
			return common.BadRequest(c, "Insufficient shares")
		default:
			return common.InternalError(c, err.Error())
		}
	}

	return common.Created(c, result, "Trade executed successfully")
}

// GetTrades returns trades for current user
// GET /api/v1/trades
func (ctrl *TradeController) GetTrades(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return common.Unauthorized(c, "User not authenticated")
	}

	filter := &dto.TradeFilter{
		Symbol: c.Query("symbol"),
		Side:   c.Query("side"),
		Page:   c.QueryInt("page", 1),
		Limit:  c.QueryInt("limit", 20),
	}

	result, err := ctrl.tradeService.GetTrades(c.Context(), userID, filter)
	if err != nil {
		return common.InternalError(c, err.Error())
	}

	return common.Success(c, result, "")
}

// GetTradeByID returns a single trade
// GET /api/v1/trades/:id
func (ctrl *TradeController) GetTradeByID(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return common.Unauthorized(c, "User not authenticated")
	}

	tradeID := c.Params("id")
	result, err := ctrl.tradeService.GetTradeByID(c.Context(), tradeID, userID)
	if err != nil {
		if errors.Is(err, service.ErrTradeNotFound) {
			return common.NotFound(c, "Trade not found")
		}
		if errors.Is(err, service.ErrUnauthorized) {
			return common.Unauthorized(c, "Access denied")
		}
		return common.InternalError(c, err.Error())
	}

	return common.Success(c, result, "")
}

// GetTradesByOrderID returns all trades for an order
// GET /api/v1/orders/:id/trades
func (ctrl *TradeController) GetTradesByOrderID(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return common.Unauthorized(c, "User not authenticated")
	}

	orderID := c.Params("id")
	result, err := ctrl.tradeService.GetTradesByOrderID(c.Context(), orderID, userID)
	if err != nil {
		if errors.Is(err, service.ErrOrderNotFound) {
			return common.NotFound(c, "Order not found")
		}
		if errors.Is(err, service.ErrUnauthorized) {
			return common.Unauthorized(c, "Access denied")
		}
		return common.InternalError(c, err.Error())
	}

	return common.Success(c, result, "")
}

// GetTradeSummary returns aggregate stats
// GET /api/v1/trades/summary
func (ctrl *TradeController) GetTradeSummary(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return common.Unauthorized(c, "User not authenticated")
	}

	result, err := ctrl.tradeService.GetTradeSummary(c.Context(), userID)
	if err != nil {
		return common.InternalError(c, err.Error())
	}

	return common.Success(c, result, "")
}
