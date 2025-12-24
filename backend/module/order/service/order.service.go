package service

import (
	"context"
	"errors"

	"github.com/bricksocoolxd/bengi-investment-system/module/order/dto"
	"github.com/bricksocoolxd/bengi-investment-system/module/order/model"
	"github.com/bricksocoolxd/bengi-investment-system/module/order/repository"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/ws"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrOrderNotFound       = errors.New("order not found")
	ErrUnauthorized        = errors.New("unauthorized access")
	ErrCannotCancelOrder   = errors.New("order cannot be cancelled")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrInvalidOrderType    = errors.New("invalid order type")
)

type OrderService struct {
	repo *repository.OrderRepository
}

func NewOrderService(repo *repository.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
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

	order := &model.Order{
		UserID:       userObjectID,
		AccountID:    accountObjectID,
		PortfolioID:  portfolioObjectID,
		InstrumentID: primitive.NewObjectID(), // TODO: lookup from instrument service
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

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, err
	}

	ws.PublishOrderUpdate(userID, &ws.OrderPayload{
		OrderID:   order.ID.Hex(),
		Symbol:    order.Symbol,
		Side:      string(order.Side),
		Status:    string(order.Status),
		FilledQty: 0,
	})

	return s.toOrderResponse(order), nil
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
