package service

import (
	"context"
	"errors"

	"github.com/bricksocoolxd/bengi-investment-system/module/instrument/dto"
	"github.com/bricksocoolxd/bengi-investment-system/module/instrument/model"
	"github.com/bricksocoolxd/bengi-investment-system/module/instrument/repository"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	ErrInstrumentNotFound = errors.New("instrument not found")
	ErrSymbolExists       = errors.New("symbol already exists")
)

type InstrumentService struct {
	repository        *repository.InstrumentRepository
	MarketDataService *MarketDataService
}

func NewInstrumentService(repo *repository.InstrumentRepository, marketSvc *MarketDataService) *InstrumentService {
	return &InstrumentService{
		repository:        repo,
		MarketDataService: marketSvc,
	}
}

func (s *InstrumentService) CreateInstrument(ctx context.Context, req *dto.CreateInstrumentRequest) (*dto.InstrumentResponse, error) {

	exists, err := s.repository.SymbolExists(ctx, req.Symbol)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrSymbolExists
	}

	instrument := &model.Instrument{
		Symbol:      req.Symbol,
		Name:        req.Name,
		Type:        model.InstrumentType(req.Type),
		Exchange:    req.Exchange,
		Currency:    req.Currency,
		Description: req.Description,
		LogoURL:     req.LogoURL,
	}

	if err := s.repository.CreateInstrument(ctx, instrument); err != nil {
		return nil, err
	}

	return s.toInstrumentResponse(instrument), nil
}

// GetInstruments returns paginated list of instruments
func (s *InstrumentService) GetInstruments(ctx context.Context, page, limit int) (*dto.InstrumentListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	instruments, total, err := s.repository.FindAll(ctx, page, limit)
	if err != nil {
		return nil, err
	}

	var responses []dto.InstrumentResponse
	for _, inst := range instruments {
		responses = append(responses, *s.toInstrumentResponse(&inst))
	}

	return &dto.InstrumentListResponse{
		Instruments: responses,
		Total:       int(total),
		Page:        page,
		Limit:       limit,
	}, nil
}

// GetInstrumentBySymbol returns a single instrument by symbol
func (s *InstrumentService) GetInstrumentBySymbol(ctx context.Context, symbol string) (*dto.InstrumentResponse, error) {
	instrument, err := s.repository.FindBySymbol(ctx, symbol)
	if err != nil {
		return nil, ErrInstrumentNotFound
	}

	return s.toInstrumentResponse(instrument), nil
}

// SearchInstruments searches instruments
func (s *InstrumentService) SearchInstruments(ctx context.Context, query *dto.SearchQuery) (*dto.InstrumentListResponse, error) {
	page := query.Page
	limit := query.Limit

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	instruments, total, err := s.repository.Search(ctx, query.Query, query.Type, query.Exchange, page, limit)
	if err != nil {
		return nil, err
	}

	var responses []dto.InstrumentResponse
	for _, inst := range instruments {
		responses = append(responses, *s.toInstrumentResponse(&inst))
	}

	return &dto.InstrumentListResponse{
		Instruments: responses,
		Total:       int(total),
		Page:        page,
		Limit:       limit,
	}, nil
}

// GetQuote returns real-time quote for a symbol
func (s *InstrumentService) GetQuote(ctx context.Context, symbol string) (*dto.QuoteResponse, error) {
	// First verify instrument exists in our database
	_, err := s.repository.FindBySymbol(ctx, symbol)
	if err != nil {
		return nil, ErrInstrumentNotFound
	}

	// Fetch live quote from market data service
	quote, err := s.MarketDataService.GetQuote(symbol)
	if err != nil {
		return nil, err
	}

	return s.toQuoteResponse(quote), nil
}

// UpdateInstrument updates an instrument (admin only)
func (s *InstrumentService) UpdateInstrument(ctx context.Context, symbol string, req *dto.UpdateInstrumentRequest) (*dto.InstrumentResponse, error) {
	instrument, err := s.repository.FindBySymbol(ctx, symbol)
	if err != nil {
		return nil, ErrInstrumentNotFound
	}

	// Build update
	update := bson.M{}
	if req.Name != "" {
		update["name"] = req.Name
	}
	if req.Description != "" {
		update["description"] = req.Description
	}
	if req.LogoURL != "" {
		update["logoUrl"] = req.LogoURL
	}
	if req.Status != "" {
		update["status"] = req.Status
	}

	if len(update) > 0 {
		if err := s.repository.Update(ctx, instrument.ID, update); err != nil {
			return nil, err
		}
	}

	// Fetch updated instrument
	updated, err := s.repository.FindBySymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}

	return s.toInstrumentResponse(updated), nil
}

// Helper: Convert Instrument to InstrumentResponse
func (s *InstrumentService) toInstrumentResponse(inst *model.Instrument) *dto.InstrumentResponse {
	return &dto.InstrumentResponse{
		ID:          inst.ID.Hex(),
		Symbol:      inst.Symbol,
		Name:        inst.Name,
		Type:        string(inst.Type),
		Exchange:    inst.Exchange,
		Currency:    inst.Currency,
		Status:      string(inst.Status),
		Description: inst.Description,
		LogoURL:     inst.LogoURL,
	}
}

// Helper: Convert Quote to QuoteResponse
func (s *InstrumentService) toQuoteResponse(quote *model.Quote) *dto.QuoteResponse {
	return &dto.QuoteResponse{
		Symbol:        quote.Symbol,
		Price:         quote.Price,
		Open:          quote.Open,
		High:          quote.High,
		Low:           quote.Low,
		Close:         quote.Close,
		PreviousClose: quote.PreviousClose,
		Volume:        quote.Volume,
		Change:        quote.Change,
		ChangePercent: quote.ChangePercent,
		Timestamp:     quote.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
	}
}
