package controller

import (
	"errors"

	"github.com/bricksocoolxd/bengi-investment-system/module/watchlist/dto"
	"github.com/bricksocoolxd/bengi-investment-system/module/watchlist/service"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/common"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/middleware"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type WatchlistController struct {
	watchlistService *service.WatchlistService
}

func NewWatchlistController(watchlistService *service.WatchlistService) *WatchlistController {
	return &WatchlistController{
		watchlistService: watchlistService,
	}
}

// CreateWatchlist creates a new watchlist
// POST /api/v1/watchlists
func (ctrl *WatchlistController) CreateWatchlist(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return common.Unauthorized(c, "User not authenticated")
	}

	var req dto.CreateWatchlistRequest
	if validationErrors := utils.ParseAndValidate(c, &req); validationErrors != nil {
		return common.ValidationError(c, validationErrors)
	}

	result, err := ctrl.watchlistService.CreateWatchlist(c.Context(), userID, &req)
	if err != nil {
		if errors.Is(err, service.ErrMaxWatchlistsReached) {
			return common.BadRequest(c, "Maximum number of watchlists reached")
		}
		return common.InternalError(c, err.Error())
	}

	return common.Created(c, result, "Watchlist created successfully")
}

// GetWatchlists returns all watchlists for current user
// GET /api/v1/watchlists
func (ctrl *WatchlistController) GetWatchlists(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return common.Unauthorized(c, "User not authenticated")
	}

	result, err := ctrl.watchlistService.GetWatchlists(c.Context(), userID)
	if err != nil {
		return common.InternalError(c, err.Error())
	}

	return common.Success(c, result, "")
}

// GetWatchlistByID returns a single watchlist
// GET /api/v1/watchlists/:id
func (ctrl *WatchlistController) GetWatchlistByID(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return common.Unauthorized(c, "User not authenticated")
	}

	watchlistID := c.Params("id")
	result, err := ctrl.watchlistService.GetWatchlistByID(c.Context(), watchlistID, userID)
	if err != nil {
		if errors.Is(err, service.ErrWatchlistNotFound) {
			return common.NotFound(c, "Watchlist not found")
		}
		if errors.Is(err, service.ErrUnauthorized) {
			return common.Unauthorized(c, "Access denied")
		}
		return common.InternalError(c, err.Error())
	}

	return common.Success(c, result, "")
}

// UpdateWatchlist updates a watchlist
// PUT /api/v1/watchlists/:id
func (ctrl *WatchlistController) UpdateWatchlist(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return common.Unauthorized(c, "User not authenticated")
	}

	watchlistID := c.Params("id")
	var req dto.UpdateWatchlistRequest
	if validationErrors := utils.ParseAndValidate(c, &req); validationErrors != nil {
		return common.ValidationError(c, validationErrors)
	}

	result, err := ctrl.watchlistService.UpdateWatchlist(c.Context(), watchlistID, userID, &req)
	if err != nil {
		if errors.Is(err, service.ErrWatchlistNotFound) {
			return common.NotFound(c, "Watchlist not found")
		}
		if errors.Is(err, service.ErrUnauthorized) {
			return common.Unauthorized(c, "Access denied")
		}
		return common.InternalError(c, err.Error())
	}

	return common.Success(c, result, "Watchlist updated successfully")
}

// DeleteWatchlist deletes a watchlist
// DELETE /api/v1/watchlists/:id
func (ctrl *WatchlistController) DeleteWatchlist(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return common.Unauthorized(c, "User not authenticated")
	}

	watchlistID := c.Params("id")
	err := ctrl.watchlistService.DeleteWatchlist(c.Context(), watchlistID, userID)
	if err != nil {
		if errors.Is(err, service.ErrWatchlistNotFound) {
			return common.NotFound(c, "Watchlist not found")
		}
		if errors.Is(err, service.ErrUnauthorized) {
			return common.Unauthorized(c, "Access denied")
		}
		return common.InternalError(c, err.Error())
	}

	return common.Success(c, nil, "Watchlist deleted successfully")
}

// AddSymbol adds a symbol to watchlist
// POST /api/v1/watchlists/:id/symbols
func (ctrl *WatchlistController) AddSymbol(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return common.Unauthorized(c, "User not authenticated")
	}

	watchlistID := c.Params("id")
	var req dto.AddSymbolRequest
	if validationErrors := utils.ParseAndValidate(c, &req); validationErrors != nil {
		return common.ValidationError(c, validationErrors)
	}

	result, err := ctrl.watchlistService.AddSymbol(c.Context(), watchlistID, userID, &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrWatchlistNotFound):
			return common.NotFound(c, "Watchlist not found")
		case errors.Is(err, service.ErrUnauthorized):
			return common.Unauthorized(c, "Access denied")
		case errors.Is(err, service.ErrMaxSymbolsReached):
			return common.BadRequest(c, "Maximum symbols reached")
		case errors.Is(err, service.ErrSymbolAlreadyExists):
			return common.BadRequest(c, "Symbol already in watchlist")
		default:
			return common.InternalError(c, err.Error())
		}
	}

	return common.Success(c, result, "Symbol added successfully")
}

// RemoveSymbol removes a symbol from watchlist
// DELETE /api/v1/watchlists/:id/symbols/:symbol
func (ctrl *WatchlistController) RemoveSymbol(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return common.Unauthorized(c, "User not authenticated")
	}

	watchlistID := c.Params("id")
	symbol := c.Params("symbol")

	req := &dto.RemoveSymbolRequest{Symbol: symbol}
	result, err := ctrl.watchlistService.RemoveSymbol(c.Context(), watchlistID, userID, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrWatchlistNotFound):
			return common.NotFound(c, "Watchlist not found")
		case errors.Is(err, service.ErrUnauthorized):
			return common.Unauthorized(c, "Access denied")
		case errors.Is(err, service.ErrSymbolNotFound):
			return common.NotFound(c, "Symbol not in watchlist")
		default:
			return common.InternalError(c, err.Error())
		}
	}

	return common.Success(c, result, "Symbol removed successfully")
}
