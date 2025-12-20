package service

import (
	"context"
	"errors"

	"github.com/bricksocoolxd/bengi-investment-system/module/portfolio/dto"
	"github.com/bricksocoolxd/bengi-investment-system/module/portfolio/model"
	"github.com/bricksocoolxd/bengi-investment-system/module/portfolio/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrPortfolioNotFound = errors.New("portfolio not found")
	ErrPositionNotFound  = errors.New("position not found")
	ErrUnauthorized      = errors.New("unauthorized access")
)

type PortfolioService struct {
	repo *repository.PortfolioRepository
}

func NewPortfolioService(repo *repository.PortfolioRepository) *PortfolioService {
	return &PortfolioService{repo: repo}
}

// ==================== Portfolio Methods ====================

func (s *PortfolioService) CreatePortfolio(ctx context.Context, userID string, req *dto.CreatePortfolioRequest) (*dto.PortfolioResponse, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	accountObjectID, err := primitive.ObjectIDFromHex(req.AccountID)
	if err != nil {
		return nil, err
	}

	portfolio := &model.Portfolio{
		UserID:    userObjectID,
		AccountID: accountObjectID,
		Name:      req.Name,
		IsDefault: false,
	}

	if err := s.repo.CreatePortfolio(ctx, portfolio); err != nil {
		return nil, err
	}

	return s.toPortfolioResponse(portfolio), nil
}

func (s *PortfolioService) GetPortfolios(ctx context.Context, userID string) ([]dto.PortfolioResponse, error) {
	portfolios, err := s.repo.FindPortfoliosByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var responses []dto.PortfolioResponse
	for _, p := range portfolios {
		responses = append(responses, *s.toPortfolioResponse(&p))
	}
	return responses, nil
}

func (s *PortfolioService) GetPortfolioByID(ctx context.Context, portfolioID, userID string) (*dto.PortfolioResponse, error) {
	portfolio, err := s.repo.FindPortfolioByID(ctx, portfolioID)
	if err != nil {
		return nil, ErrPortfolioNotFound
	}

	if portfolio.UserID.Hex() != userID {
		return nil, ErrUnauthorized
	}

	return s.toPortfolioResponse(portfolio), nil
}

func (s *PortfolioService) GetPortfolioSummary(ctx context.Context, portfolioID, userID string) (*dto.PortfolioSummary, error) {
	portfolio, err := s.repo.FindPortfolioByID(ctx, portfolioID)
	if err != nil {
		return nil, ErrPortfolioNotFound
	}

	if portfolio.UserID.Hex() != userID {
		return nil, ErrUnauthorized
	}

	positions, err := s.repo.FindPositionsByPortfolioID(ctx, portfolioID)
	if err != nil {
		return nil, err
	}

	var positionResponses []dto.PositionResponse
	var totalCost float64
	var totalValue float64

	for _, pos := range positions {
		resp := s.toPositionResponse(&pos)
		positionResponses = append(positionResponses, *resp)
		totalCost += pos.TotalCost
		totalValue += pos.TotalCost // Will be replaced with market value when integrated with price service
	}

	totalPnL := totalValue - totalCost
	var totalPnLPct float64
	if totalCost > 0 {
		totalPnLPct = (totalPnL / totalCost) * 100
	}

	return &dto.PortfolioSummary{
		Portfolio:   *s.toPortfolioResponse(portfolio),
		Positions:   positionResponses,
		TotalValue:  totalValue,
		TotalCost:   totalCost,
		TotalPnL:    totalPnL,
		TotalPnLPct: totalPnLPct,
	}, nil
}

func (s *PortfolioService) UpdatePortfolio(ctx context.Context, portfolioID, userID string, req *dto.UpdatePortfolioRequest) (*dto.PortfolioResponse, error) {
	portfolio, err := s.repo.FindPortfolioByID(ctx, portfolioID)
	if err != nil {
		return nil, ErrPortfolioNotFound
	}

	if portfolio.UserID.Hex() != userID {
		return nil, ErrUnauthorized
	}

	update := bson.M{}
	if req.Name != "" {
		update["name"] = req.Name
	}

	if req.IsDefault {
		if err := s.repo.ClearDefaultPortfolios(ctx, portfolio.UserID); err != nil {
			return nil, err
		}
		update["isDefault"] = true
	}

	if len(update) > 0 {
		if err := s.repo.UpdatePortfolio(ctx, portfolio.ID, update); err != nil {
			return nil, err
		}
	}

	updated, err := s.repo.FindPortfolioByID(ctx, portfolioID)
	if err != nil {
		return nil, err
	}

	return s.toPortfolioResponse(updated), nil
}

func (s *PortfolioService) DeletePortfolio(ctx context.Context, portfolioID, userID string) error {
	portfolio, err := s.repo.FindPortfolioByID(ctx, portfolioID)
	if err != nil {
		return ErrPortfolioNotFound
	}

	if portfolio.UserID.Hex() != userID {
		return ErrUnauthorized
	}

	return s.repo.DeletePortfolio(ctx, portfolio.ID)
}

// ==================== Position Methods ====================

func (s *PortfolioService) GetPositions(ctx context.Context, portfolioID, userID string) ([]dto.PositionResponse, error) {
	portfolio, err := s.repo.FindPortfolioByID(ctx, portfolioID)
	if err != nil {
		return nil, ErrPortfolioNotFound
	}

	if portfolio.UserID.Hex() != userID {
		return nil, ErrUnauthorized
	}

	positions, err := s.repo.FindPositionsByPortfolioID(ctx, portfolioID)
	if err != nil {
		return nil, err
	}

	var responses []dto.PositionResponse
	for _, pos := range positions {
		responses = append(responses, *s.toPositionResponse(&pos))
	}
	return responses, nil
}

func (s *PortfolioService) GetPositionDetail(ctx context.Context, positionID, userID string) (*dto.PositionDetailResponse, error) {
	position, err := s.repo.FindPositionByID(ctx, positionID)
	if err != nil {
		return nil, ErrPositionNotFound
	}

	portfolio, err := s.repo.FindPortfolioByID(ctx, position.PortfolioID.Hex())
	if err != nil {
		return nil, ErrPortfolioNotFound
	}

	if portfolio.UserID.Hex() != userID {
		return nil, ErrUnauthorized
	}

	lots, err := s.repo.FindLotsByPositionID(ctx, position.ID)
	if err != nil {
		return nil, err
	}

	var lotResponses []dto.PositionLotResponse
	for _, lot := range lots {
		lotResponses = append(lotResponses, dto.PositionLotResponse{
			ID:           lot.ID.Hex(),
			Quantity:     lot.Quantity,
			RemainingQty: lot.RemainingQty,
			CostPerUnit:  lot.CostPerUnit,
			PurchasedAt:  lot.PurchasedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return &dto.PositionDetailResponse{
		Position: *s.toPositionResponse(position),
		Lots:     lotResponses,
	}, nil
}

// ==================== Helpers ====================

func (s *PortfolioService) toPortfolioResponse(p *model.Portfolio) *dto.PortfolioResponse {
	return &dto.PortfolioResponse{
		ID:        p.ID.Hex(),
		UserID:    p.UserID.Hex(),
		AccountID: p.AccountID.Hex(),
		Name:      p.Name,
		IsDefault: p.IsDefault,
	}
}

func (s *PortfolioService) toPositionResponse(pos *model.Position) *dto.PositionResponse {
	return &dto.PositionResponse{
		ID:           pos.ID.Hex(),
		PortfolioID:  pos.PortfolioID.Hex(),
		InstrumentID: pos.InstrumentID.Hex(),
		Symbol:       pos.Symbol,
		Quantity:     pos.Quantity,
		AvgCost:      pos.AvgCost,
		TotalCost:    pos.TotalCost,
	}
}
