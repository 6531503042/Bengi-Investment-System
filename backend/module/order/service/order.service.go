package service

import (
	"context"
	"errors"
	"time"

	accountRepo "github.com/bricksocoolxd/bengi-investment-system/module/account/repository"
	"github.com/bricksocoolxd/bengi-investment-system/module/order/dto"
	"github.com/bricksocoolxd/bengi-investment-system/module/order/model"
	"github.com/bricksocoolxd/bengi-investment-system/module/order/repository"
	portfolioModel "github.com/bricksocoolxd/bengi-investment-system/module/portfolio/model"
	portfolioRepo "github.com/bricksocoolxd/bengi-investment-system/module/portfolio/repository"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/ws"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrOrderNotFound       = errors.New("order not found")
	ErrUnauthorized        = errors.New("unauthorized access")
	ErrCannotCancelOrder   = errors.New("order cannot be cancelled")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrInvalidOrderType    = errors.New("invalid order type")
)

type OrderService struct {
	repo          *repository.OrderRepository
	portfolioRepo *portfolioRepo.PortfolioRepository
	accountRepo   *accountRepo.AccountRepository
}

func NewOrderService(repo *repository.OrderRepository) *OrderService {
	return &OrderService{
		repo:          repo,
		portfolioRepo: portfolioRepo.NewPortfolioRepository(),
		accountRepo:   accountRepo.NewAccountRepository(),
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, userID string, req *dto.CreateOrderRequest) (*dto.OrderResponse, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	accountObjectID, err := primitive.ObjectIDFromHex(req.AccountID)
	if err != nil {
		return nil, err
	}

	portfolioObjectID, err := primitive.ObjectIDFromHex(req.PortfolioID)
	if err != nil {
		return nil, err
	}

	// Set default TimeInForce
	timeInForce := model.TimeInForceGTC
	if req.TimeInForce != "" {
		timeInForce = model.TimeInForce(req.TimeInForce)
	}

	// Validate order type requirements
	if req.Type == "LIMIT" && req.Price <= 0 {
		return nil, ErrInvalidOrderType
	}
	if req.Type == "STOP" && req.StopPrice <= 0 {
		return nil, ErrInvalidOrderType
	}

	// For market orders, use current price (passed from frontend or mock)
	fillPrice := req.Price
	if fillPrice <= 0 {
		fillPrice = 100.0 // Default mock price if not provided
	}
	totalCost := req.Quantity * fillPrice

	// Check balance for BUY orders
	if req.Side == "BUY" {
		account, err := s.accountRepo.FindByID(ctx, req.AccountID)
		if err != nil {
			return nil, err
		}
		if account.Balance < totalCost {
			return nil, ErrInsufficientBalance
		}
	}

	order := &model.Order{
		UserID:       userObjectID,
		AccountID:    accountObjectID,
		PortfolioID:  portfolioObjectID,
		InstrumentID: primitive.NewObjectID(),
		Symbol:       req.Symbol,
		Side:         model.OrderSide(req.Side),
		Type:         model.OrderType(req.Type),
		TimeInForce:  timeInForce,
		Quantity:     req.Quantity,
		FilledQty:    0,
		Price:        req.Price,
		StopPrice:    req.StopPrice,
		Commission:   0,
	}

	// Create order first
	if err := s.repo.Create(ctx, order); err != nil {
		return nil, err
	}

	// For MARKET orders, execute immediately
	if req.Type == "MARKET" {
		if err := s.executeMarketOrder(ctx, order, fillPrice); err != nil {
			// Update order status to REJECTED
			s.repo.UpdateStatus(ctx, order.ID, model.OrderStatusRejected)
			return nil, err
		}
	}

	ws.PublishOrderUpdate(userID, &ws.OrderPayload{
		OrderID:   order.ID.Hex(),
		Symbol:    order.Symbol,
		Side:      string(order.Side),
		Status:    string(order.Status),
		FilledQty: order.FilledQty,
		AvgPrice:  order.AvgFillPrice,
	})

	return s.toOrderResponse(order), nil
}

// executeMarketOrder executes a market order immediately
func (s *OrderService) executeMarketOrder(ctx context.Context, order *model.Order, fillPrice float64) error {
	now := time.Now()
	totalCost := order.Quantity * fillPrice
	commission := totalCost * 0.001 // 0.1% commission

	// Update order as filled
	order.Status = model.OrderStatusFilled
	order.FilledQty = order.Quantity
	order.AvgFillPrice = fillPrice
	order.Commission = commission
	order.FilledAt = &now

	if err := s.repo.UpdateFill(ctx, order.ID, order.Quantity, fillPrice, model.OrderStatusFilled); err != nil {
		return err
	}

	// Update/create position in portfolio
	if order.Side == model.OrderSideBuy {
		if err := s.addPosition(ctx, order, fillPrice); err != nil {
			return err
		}
		// Deduct balance
		if err := s.accountRepo.UpdateBalanceDelta(ctx, order.AccountID, -totalCost-commission); err != nil {
			return err
		}
	} else {
		// SELL - reduce position and credit balance
		if err := s.reducePosition(ctx, order, fillPrice); err != nil {
			return err
		}
		if err := s.accountRepo.UpdateBalanceDelta(ctx, order.AccountID, totalCost-commission); err != nil {
			return err
		}
	}

	return nil
}

// addPosition creates or updates a position for a BUY order
func (s *OrderService) addPosition(ctx context.Context, order *model.Order, price float64) error {
	existingPos, err := s.portfolioRepo.FindPositionByPortfolioAndSymbol(ctx, order.PortfolioID, order.Symbol)

	if err == mongo.ErrNoDocuments || existingPos == nil {
		// Create new position
		newPosition := &portfolioModel.Position{
			PortfolioID:  order.PortfolioID,
			InstrumentID: order.InstrumentID,
			Symbol:       order.Symbol,
			Quantity:     order.Quantity,
			AvgCost:      price,
			TotalCost:    order.Quantity * price,
		}
		return s.portfolioRepo.CreatePosition(ctx, newPosition)
	}

	// Update existing position with weighted average
	newTotalQty := existingPos.Quantity + order.Quantity
	newTotalCost := existingPos.TotalCost + (order.Quantity * price)
	newAvgCost := newTotalCost / newTotalQty

	return s.portfolioRepo.UpdatePosition(ctx, existingPos.ID, bson.M{
		"quantity":  newTotalQty,
		"avgCost":   newAvgCost,
		"totalCost": newTotalCost,
	})
}

// reducePosition reduces a position for a SELL order
func (s *OrderService) reducePosition(ctx context.Context, order *model.Order, price float64) error {
	existingPos, err := s.portfolioRepo.FindPositionByPortfolioAndSymbol(ctx, order.PortfolioID, order.Symbol)
	if err != nil {
		return errors.New("no position to sell")
	}

	if existingPos.Quantity < order.Quantity {
		return errors.New("insufficient shares")
	}

	newQty := existingPos.Quantity - order.Quantity
	if newQty <= 0 {
		// Delete position if fully sold
		return s.portfolioRepo.DeletePosition(ctx, existingPos.ID)
	}

	newTotalCost := newQty * existingPos.AvgCost
	return s.portfolioRepo.UpdatePosition(ctx, existingPos.ID, bson.M{
		"quantity":  newQty,
		"totalCost": newTotalCost,
	})
}

func (s *OrderService) GetOrders(ctx context.Context, userID string, filter *dto.OrderFilter) (*dto.OrderListResponse, error) {
	page := filter.Page
	limit := filter.Limit

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	modelFilter := &model.OrderFilter{
		Status: filter.Status,
		Side:   filter.Side,
		Symbol: filter.Symbol,
	}

	orders, total, err := s.repo.FindByUserID(ctx, userID, modelFilter, page, limit)
	if err != nil {
		return nil, err
	}

	var responses []dto.OrderResponse
	for _, order := range orders {
		responses = append(responses, *s.toOrderResponse(&order))
	}

	return &dto.OrderListResponse{
		Orders: responses,
		Total:  int(total),
		Page:   page,
		Limit:  limit,
	}, nil
}

func (s *OrderService) GetOrderByID(ctx context.Context, orderID, userID string) (*dto.OrderResponse, error) {
	order, err := s.repo.FindByID(ctx, orderID)
	if err != nil {
		return nil, ErrOrderNotFound
	}

	if order.UserID.Hex() != userID {
		return nil, ErrUnauthorized
	}

	return s.toOrderResponse(order), nil
}

func (s *OrderService) CancelOrder(ctx context.Context, orderID, userID string) (*dto.OrderResponse, error) {
	order, err := s.repo.FindByID(ctx, orderID)
	if err != nil {
		return nil, ErrOrderNotFound
	}

	if order.UserID.Hex() != userID {
		return nil, ErrUnauthorized
	}

	// Check if order can be cancelled
	if order.Status != model.OrderStatusPending &&
		order.Status != model.OrderStatusOpen &&
		order.Status != model.OrderStatusPartiallyFilled {
		return nil, ErrCannotCancelOrder
	}

	if err := s.repo.UpdateStatus(ctx, order.ID, model.OrderStatusCancelled); err != nil {
		return nil, err
	}

	ws.PublishOrderUpdate(userID, &ws.OrderPayload{
		OrderID:   order.ID.Hex(),
		Symbol:    order.Symbol,
		Side:      string(order.Side),
		Status:    string(model.OrderStatusCancelled),
		FilledQty: order.FilledQty,
		AvgPrice:  order.AvgFillPrice,
	})

	// Fetch updated order
	updated, err := s.repo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	return s.toOrderResponse(updated), nil
}

func (s *OrderService) toOrderResponse(order *model.Order) *dto.OrderResponse {
	resp := &dto.OrderResponse{
		ID:           order.ID.Hex(),
		UserID:       order.UserID.Hex(),
		AccountID:    order.AccountID.Hex(),
		PortfolioID:  order.PortfolioID.Hex(),
		InstrumentID: order.InstrumentID.Hex(),
		Symbol:       order.Symbol,
		Side:         string(order.Side),
		Type:         string(order.Type),
		Status:       string(order.Status),
		TimeInForce:  string(order.TimeInForce),
		Quantity:     order.Quantity,
		FilledQty:    order.FilledQty,
		Price:        order.Price,
		StopPrice:    order.StopPrice,
		AvgFillPrice: order.AvgFillPrice,
		Commission:   order.Commission,
		CreatedAt:    order.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if order.FilledAt != nil {
		t := order.FilledAt.Format("2006-01-02T15:04:05Z07:00")
		resp.FilledAt = &t
	}
	if order.CancelledAt != nil {
		t := order.CancelledAt.Format("2006-01-02T15:04:05Z07:00")
		resp.CancelledAt = &t
	}

	return resp
}
