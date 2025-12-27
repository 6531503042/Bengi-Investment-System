package controller

import (
	"errors"
	"net/url"
	"time"

	"github.com/bricksocoolxd/bengi-investment-system/module/instrument/dto"
	"github.com/bricksocoolxd/bengi-investment-system/module/instrument/service"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/common"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type InstrumentController struct {
	instrumentService *service.InstrumentService
}

func NewInstrumentController(instrumentService *service.InstrumentService) *InstrumentController {
	return &InstrumentController{
		instrumentService: instrumentService,
	}
}

// GetInstruments returns paginated list of instruments
// GET /api/v1/instruments
func (ctrl *InstrumentController) GetInstruments(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)

	result, err := ctrl.instrumentService.GetInstruments(c.Context(), page, limit)
	if err != nil {
		return common.InternalError(c, err.Error())
	}

	return common.Success(c, result, "")
}

// SearchInstruments searches instruments
// GET /api/v1/instruments/search?q=AAPL&type=STOCK
func (ctrl *InstrumentController) SearchInstruments(c *fiber.Ctx) error {
	query := &dto.SearchQuery{
		Query:    c.Query("q"),
		Type:     c.Query("type"),
		Exchange: c.Query("exchange"),
		Page:     c.QueryInt("page", 1),
		Limit:    c.QueryInt("limit", 20),
	}

	result, err := ctrl.instrumentService.SearchInstruments(c.Context(), query)
	if err != nil {
		return common.InternalError(c, err.Error())
	}

	return common.Success(c, result, "")
}

// GetInstrumentBySymbol returns a single instrument
// GET /api/v1/instruments/:symbol
func (ctrl *InstrumentController) GetInstrumentBySymbol(c *fiber.Ctx) error {
	symbol, _ := url.PathUnescape(c.Params("symbol"))

	result, err := ctrl.instrumentService.GetInstrumentBySymbol(c.Context(), symbol)
	if err != nil {
		if errors.Is(err, service.ErrInstrumentNotFound) {
			return common.NotFound(c, "Instrument not found")
		}
		return common.InternalError(c, err.Error())
	}

	return common.Success(c, result, "")
}

// GetQuote returns real-time quote for a symbol
// GET /api/v1/instruments/:symbol/quote
func (ctrl *InstrumentController) GetQuote(c *fiber.Ctx) error {
	symbol, _ := url.PathUnescape(c.Params("symbol"))

	result, err := ctrl.instrumentService.GetQuote(c.Context(), symbol)
	if err != nil {
		if errors.Is(err, service.ErrInstrumentNotFound) {
			return common.NotFound(c, "Instrument not found")
		}
		if errors.Is(err, service.ErrQuoteNotFound) || errors.Is(err, service.ErrAPIError) {
			return common.InternalError(c, "Unable to fetch quote")
		}
		return common.InternalError(c, err.Error())
	}

	return common.Success(c, result, "")
}

// GetCandles returns historical candlestick data for a symbol
// GET /api/v1/instruments/:symbol/candles?resolution=D&from=1234567890&to=1234567890
func (ctrl *InstrumentController) GetCandles(c *fiber.Ctx) error {
	symbol, _ := url.PathUnescape(c.Params("symbol"))
	resolution := c.Query("resolution", "D") // Default to daily
	from := c.QueryInt("from", 0)
	to := c.QueryInt("to", 0)

	// Default to last 30 days if no time range specified
	now := time.Now().Unix()
	if from == 0 {
		from = int(now) - (30 * 24 * 60 * 60) // 30 days ago
	}
	if to == 0 {
		to = int(now)
	}

	result, err := ctrl.instrumentService.GetCandles(c.Context(), symbol, resolution, int64(from), int64(to))
	if err != nil {
		if errors.Is(err, service.ErrInstrumentNotFound) {
			return common.NotFound(c, "Instrument not found")
		}
		if errors.Is(err, service.ErrQuoteNotFound) {
			return common.NotFound(c, "No candle data available")
		}
		return common.InternalError(c, err.Error())
	}

	return common.Success(c, result, "")
}

// CreateInstrument creates a new instrument (admin only)
// POST /api/v1/instruments
func (ctrl *InstrumentController) CreateInstrument(c *fiber.Ctx) error {
	var req dto.CreateInstrumentRequest
	if validationErrors := utils.ParseAndValidate(c, &req); validationErrors != nil {
		return common.ValidationError(c, validationErrors)
	}

	result, err := ctrl.instrumentService.CreateInstrument(c.Context(), &req)
	if err != nil {
		if errors.Is(err, service.ErrSymbolExists) {
			return common.BadRequest(c, "Symbol already exists")
		}
		return common.InternalError(c, err.Error())
	}

	return common.Created(c, result, "Instrument created successfully")
}

// UpdateInstrument updates an instrument (admin only)
// PUT /api/v1/instruments/:symbol
func (ctrl *InstrumentController) UpdateInstrument(c *fiber.Ctx) error {
	symbol := c.Params("symbol")

	var req dto.UpdateInstrumentRequest
	if validationErrors := utils.ParseAndValidate(c, &req); validationErrors != nil {
		return common.ValidationError(c, validationErrors)
	}

	result, err := ctrl.instrumentService.UpdateInstrument(c.Context(), symbol, &req)
	if err != nil {
		if errors.Is(err, service.ErrInstrumentNotFound) {
			return common.NotFound(c, "Instrument not found")
		}
		return common.InternalError(c, err.Error())
	}

	return common.Success(c, result, "Instrument updated successfully")
}
