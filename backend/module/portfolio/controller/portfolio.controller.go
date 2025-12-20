package controller

import (
	"errors"

	"github.com/bricksocoolxd/bengi-investment-system/module/portfolio/dto"
	"github.com/bricksocoolxd/bengi-investment-system/module/portfolio/service"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/common"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/middleware"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type PortfolioController struct {
	portfolioService *service.PortfolioService
}

func NewPortfolioController(portfolioService *service.PortfolioService) *PortfolioController {
	return &PortfolioController{
		portfolioService: portfolioService,
	}
}

// CreatePortfolio creates a new portfolio
// POST /api/v1/portfolios
func (ctrl *PortfolioController) CreatePortfolio(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return common.Unauthorized(c, "User not authenticated")
	}

	var req dto.CreatePortfolioRequest
	if validationErrors := utils.ParseAndValidate(c, &req); validationErrors != nil {
		return common.ValidationError(c, validationErrors)
	}

	result, err := ctrl.portfolioService.CreatePortfolio(c.Context(), userID, &req)
	if err != nil {
		return common.InternalError(c, err.Error())
	}

	return common.Created(c, result, "Portfolio created successfully")
}

// GetPortfolios returns all portfolios for current user
// GET /api/v1/portfolios
func (ctrl *PortfolioController) GetPortfolios(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return common.Unauthorized(c, "User not authenticated")
	}

	result, err := ctrl.portfolioService.GetPortfolios(c.Context(), userID)
	if err != nil {
		return common.InternalError(c, err.Error())
	}

	return common.Success(c, result, "")
}

// GetPortfolioByID returns a single portfolio
// GET /api/v1/portfolios/:id
func (ctrl *PortfolioController) GetPortfolioByID(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return common.Unauthorized(c, "User not authenticated")
	}

	portfolioID := c.Params("id")
	result, err := ctrl.portfolioService.GetPortfolioByID(c.Context(), portfolioID, userID)
	if err != nil {
		if errors.Is(err, service.ErrPortfolioNotFound) {
			return common.NotFound(c, "Portfolio not found")
		}
		if errors.Is(err, service.ErrUnauthorized) {
			return common.Unauthorized(c, "Access denied")
		}
		return common.InternalError(c, err.Error())
	}

	return common.Success(c, result, "")
}

// GetPortfolioSummary returns portfolio with positions and P&L
// GET /api/v1/portfolios/:id/summary
func (ctrl *PortfolioController) GetPortfolioSummary(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return common.Unauthorized(c, "User not authenticated")
	}

	portfolioID := c.Params("id")
	result, err := ctrl.portfolioService.GetPortfolioSummary(c.Context(), portfolioID, userID)
	if err != nil {
		if errors.Is(err, service.ErrPortfolioNotFound) {
			return common.NotFound(c, "Portfolio not found")
		}
		if errors.Is(err, service.ErrUnauthorized) {
			return common.Unauthorized(c, "Access denied")
		}
		return common.InternalError(c, err.Error())
	}

	return common.Success(c, result, "")
}

// UpdatePortfolio updates a portfolio
// PUT /api/v1/portfolios/:id
func (ctrl *PortfolioController) UpdatePortfolio(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return common.Unauthorized(c, "User not authenticated")
	}

	portfolioID := c.Params("id")
	var req dto.UpdatePortfolioRequest
	if validationErrors := utils.ParseAndValidate(c, &req); validationErrors != nil {
		return common.ValidationError(c, validationErrors)
	}

	result, err := ctrl.portfolioService.UpdatePortfolio(c.Context(), portfolioID, userID, &req)
	if err != nil {
		if errors.Is(err, service.ErrPortfolioNotFound) {
			return common.NotFound(c, "Portfolio not found")
		}
		if errors.Is(err, service.ErrUnauthorized) {
			return common.Unauthorized(c, "Access denied")
		}
		return common.InternalError(c, err.Error())
	}

	return common.Success(c, result, "Portfolio updated successfully")
}

// DeletePortfolio deletes a portfolio
// DELETE /api/v1/portfolios/:id
func (ctrl *PortfolioController) DeletePortfolio(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return common.Unauthorized(c, "User not authenticated")
	}

	portfolioID := c.Params("id")
	err := ctrl.portfolioService.DeletePortfolio(c.Context(), portfolioID, userID)
	if err != nil {
		if errors.Is(err, service.ErrPortfolioNotFound) {
			return common.NotFound(c, "Portfolio not found")
		}
		if errors.Is(err, service.ErrUnauthorized) {
			return common.Unauthorized(c, "Access denied")
		}
		return common.InternalError(c, err.Error())
	}

	return common.Success(c, nil, "Portfolio deleted successfully")
}

// GetPositions returns all positions in a portfolio
// GET /api/v1/portfolios/:id/positions
func (ctrl *PortfolioController) GetPositions(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return common.Unauthorized(c, "User not authenticated")
	}

	portfolioID := c.Params("id")
	result, err := ctrl.portfolioService.GetPositions(c.Context(), portfolioID, userID)
	if err != nil {
		if errors.Is(err, service.ErrPortfolioNotFound) {
			return common.NotFound(c, "Portfolio not found")
		}
		if errors.Is(err, service.ErrUnauthorized) {
			return common.Unauthorized(c, "Access denied")
		}
		return common.InternalError(c, err.Error())
	}

	return common.Success(c, result, "")
}

// GetPositionDetail returns position with lots
// GET /api/v1/positions/:id
func (ctrl *PortfolioController) GetPositionDetail(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return common.Unauthorized(c, "User not authenticated")
	}

	positionID := c.Params("id")
	result, err := ctrl.portfolioService.GetPositionDetail(c.Context(), positionID, userID)
	if err != nil {
		if errors.Is(err, service.ErrPositionNotFound) {
			return common.NotFound(c, "Position not found")
		}
		if errors.Is(err, service.ErrUnauthorized) {
			return common.Unauthorized(c, "Access denied")
		}
		return common.InternalError(c, err.Error())
	}

	return common.Success(c, result, "")
}
