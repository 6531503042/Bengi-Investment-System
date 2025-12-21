package controller

import (
	"errors"

	"github.com/bricksocoolxd/bengi-investment-system/module/order/dto"
	"github.com/bricksocoolxd/bengi-investment-system/module/order/service"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/common"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/middleware"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type OrderController struct {
	orderService *service.OrderService
}

func NewOrderController(orderService *service.OrderService) *OrderController {
	return &OrderController{
		orderService: orderService,
	}
}

// CreateOrder creates a new order
// POST /api/v1/orders
func (ctrl *OrderController) CreateOrder(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return common.Unauthorized(c, "User not authenticated")
	}

	var req dto.CreateOrderRequest
	if validationErrors := utils.ParseAndValidate(c, &req); validationErrors != nil {
		return common.ValidationError(c, validationErrors)
	}

	result, err := ctrl.orderService.CreateOrder(c.Context(), userID, &req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidOrderType) {
			return common.BadRequest(c, "Invalid order: LIMIT orders require price, STOP orders require stopPrice")
		}
		if errors.Is(err, service.ErrInsufficientBalance) {
			return common.BadRequest(c, "Insufficient balance")
		}
		return common.InternalError(c, err.Error())
	}

	return common.Created(c, result, "Order created successfully")
}

// GetOrders returns all orders for current user
// GET /api/v1/orders
func (ctrl *OrderController) GetOrders(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return common.Unauthorized(c, "User not authenticated")
	}

	filter := &dto.OrderFilter{
		Status: c.Query("status"),
		Side:   c.Query("side"),
		Symbol: c.Query("symbol"),
		Page:   c.QueryInt("page", 1),
		Limit:  c.QueryInt("limit", 20),
	}

	result, err := ctrl.orderService.GetOrders(c.Context(), userID, filter)
	if err != nil {
		return common.InternalError(c, err.Error())
	}

	return common.Success(c, result, "")
}

// GetOrderByID returns a single order
// GET /api/v1/orders/:id
func (ctrl *OrderController) GetOrderByID(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return common.Unauthorized(c, "User not authenticated")
	}

	orderID := c.Params("id")
	result, err := ctrl.orderService.GetOrderByID(c.Context(), orderID, userID)
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

// CancelOrder cancels an order
// POST /api/v1/orders/:id/cancel
func (ctrl *OrderController) CancelOrder(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return common.Unauthorized(c, "User not authenticated")
	}

	orderID := c.Params("id")
	result, err := ctrl.orderService.CancelOrder(c.Context(), orderID, userID)
	if err != nil {
		if errors.Is(err, service.ErrOrderNotFound) {
			return common.NotFound(c, "Order not found")
		}
		if errors.Is(err, service.ErrUnauthorized) {
			return common.Unauthorized(c, "Access denied")
		}
		if errors.Is(err, service.ErrCannotCancelOrder) {
			return common.BadRequest(c, "Order cannot be cancelled")
		}
		return common.InternalError(c, err.Error())
	}

	return common.Success(c, result, "Order cancelled successfully")
}
