package service

import (
	"context"
	"errors"
	"time"

	accountModel "github.com/bricksocoolxd/bengi-investment-system/module/account/model"
	accountRepo "github.com/bricksocoolxd/bengi-investment-system/module/account/repository"
	orderModel "github.com/bricksocoolxd/bengi-investment-system/module/order/model"
	orderRepo "github.com/bricksocoolxd/bengi-investment-system/module/order/repository"
	portfolioModel "github.com/bricksocoolxd/bengi-investment-system/module/portfolio/model"
	portfolioRepo "github.com/bricksocoolxd/bengi-investment-system/module/portfolio/repository"
	"github.com/bricksocoolxd/bengi-investment-system/module/trade/dto"
	tradeModel "github.com/bricksocoolxd/bengi-investment-system/module/trade/model"
	tradeRepo "github.com/bricksocoolxd/bengi-investment-system/module/trade/repository"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/ws"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CommissionRate = 0.001 // 0.1% commission

var (
	ErrTradeNotFound       = errors.New("trade not found")
	ErrOrderNotFound       = errors.New("order not found")
	ErrOrderAlreadyFilled  = errors.New("order is already filled")
	ErrOrderNotExecutable  = errors.New("order cannot be executed")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrInsufficientShares  = errors.New("insufficient shares to sell")
	ErrUnauthorized        = errors.New("unauthorized access")
)

type TradeService struct {
	tradeRepository     *tradeRepo.TradeRepository
	orderRepository     *orderRepo.OrderRepository
	accountRepository   *accountRepo.AccountRepository
	portfolioRepository *portfolioRepo.PortfolioRepository
}

func NewTradeService(
	tradeRepository *tradeRepo.TradeRepository,
	orderRepository *orderRepo.OrderRepository,
	accountRepository *accountRepo.AccountRepository,
	portfolioRepository *portfolioRepo.PortfolioRepository,
) *TradeService {
	return &TradeService{
		tradeRepository:     tradeRepository,
		orderRepository:     orderRepository,
		accountRepository:   accountRepository,
		portfolioRepository: portfolioRepository,
	}
}

// ExecuteTrade executes a trade for an order
func (s *TradeService) ExecuteTrade(ctx context.Context, req *dto.ExecuteTradeRequest) (*dto.TradeResponse, error) {
	// 1. Get order
	order, err := s.orderRepository.FindByID(ctx, req.OrderID)
	if err != nil {
		return nil, ErrOrderNotFound
	}

	// 2. Validate order can be executed (not already filled or cancelled)
	if order.Status == orderModel.OrderStatusFilled ||
		order.Status == orderModel.OrderStatusCancelled ||
		order.Status == orderModel.OrderStatusRejected {
		return nil, ErrOrderNotExecutable
	}

	// 3. Calculate trade values
	total := req.Quantity * req.Price
	commission := total * CommissionRate
	var netAmount float64

	if order.Side == orderModel.OrderSideBuy {
		netAmount = total + commission // Pay total + commission
	} else {
		netAmount = total - commission // Receive total - commission
	}

	// 4. Validate balance for BUY orders
	if order.Side == orderModel.OrderSideBuy {
		account, err := s.accountRepository.FindByID(ctx, order.AccountID.Hex())
		if err != nil {
			return nil, err
		}
		if account.Balance < netAmount {
			return nil, ErrInsufficientBalance
		}
	}

	// 5. Validate shares for SELL orders
	if order.Side == orderModel.OrderSideSell {
		position, err := s.portfolioRepository.FindPositionByPortfolioAndSymbol(ctx, order.PortfolioID, order.Symbol)
		if err != nil || position.Quantity < req.Quantity {
			return nil, ErrInsufficientShares
		}
	}

	// 6. Create trade record
	trade := &tradeModel.Trade{
		OrderID:      order.ID,
		UserID:       order.UserID,
		AccountID:    order.AccountID,
		PortfolioID:  order.PortfolioID,
		InstrumentID: order.InstrumentID,
		Symbol:       order.Symbol,
		Side:         tradeModel.TradeSide(order.Side),
		Quantity:     req.Quantity,
		Price:        req.Price,
		Total:        total,
		Commission:   commission,
		NetAmount:    netAmount,
		ExecutedAt:   time.Now(),
	}

	if err := s.tradeRepository.Create(ctx, trade); err != nil {
		return nil, err
	}

	// 7. Update order status
	newFilledQty := order.FilledQty + req.Quantity
	newAvgPrice := s.calculateAvgPrice(order.AvgFillPrice, order.FilledQty, req.Price, req.Quantity)

	var newStatus orderModel.OrderStatus
	if newFilledQty >= order.Quantity {
		newStatus = orderModel.OrderStatusFilled
	} else {
		newStatus = orderModel.OrderStatusPartiallyFilled
	}

	if err := s.orderRepository.UpdateFill(ctx, order.ID, newFilledQty, newAvgPrice, newStatus); err != nil {
		return nil, err
	}

	// 8. Update account balance
	account, _ := s.accountRepository.FindByID(ctx, order.AccountID.Hex())
	var newBalance float64
	if order.Side == orderModel.OrderSideBuy {
		newBalance = account.Balance - netAmount
	} else {
		newBalance = account.Balance + netAmount
	}
	s.accountRepository.UpdateBalance(ctx, account.ID, newBalance)

	// 9. Create account transaction
	s.accountRepository.CreateTransaction(ctx, &accountModel.Transaction{
		AccountID:     account.ID,
		Type:          accountModel.TransactionTypeTrade,
		Amount:        netAmount,
		BalanceBefore: account.Balance,
		BalanceAfter:  newBalance,
		Status:        accountModel.TransactionStatusCompleted,
		Description:   string(trade.Side) + " " + trade.Symbol,
	})

	// 10. Update portfolio position
	s.updatePosition(ctx, trade, order.Side)

	s.publishTradeEvents(trade, order, newFilledQty, newAvgPrice, newStatus)

	return s.toTradeResponse(trade), nil
}

// GetTrades returns trades for a user with filtering
func (s *TradeService) GetTrades(ctx context.Context, userID string, filter *dto.TradeFilter) (*dto.TradeListResponse, error) {
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

	modelFilter := &tradeModel.TradeFilter{
		Symbol: filter.Symbol,
		Side:   filter.Side,
	}

	trades, total, err := s.tradeRepository.FindByUserID(ctx, userID, modelFilter, page, limit)
	if err != nil {
		return nil, err
	}

	var responses []dto.TradeResponse
	for _, trade := range trades {
		responses = append(responses, *s.toTradeResponse(&trade))
	}

	return &dto.TradeListResponse{
		Trades: responses,
		Total:  int(total),
		Page:   page,
		Limit:  limit,
	}, nil
}

// GetTradeByID returns a single trade
func (s *TradeService) GetTradeByID(ctx context.Context, tradeID, userID string) (*dto.TradeResponse, error) {
	trade, err := s.tradeRepository.FindByID(ctx, tradeID)
	if err != nil {
		return nil, ErrTradeNotFound
	}

	if trade.UserID.Hex() != userID {
		return nil, ErrUnauthorized
	}

	return s.toTradeResponse(trade), nil
}

// GetTradesByOrderID returns all trades for an order
func (s *TradeService) GetTradesByOrderID(ctx context.Context, orderID, userID string) ([]dto.TradeResponse, error) {
	order, err := s.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		return nil, ErrOrderNotFound
	}

	if order.UserID.Hex() != userID {
		return nil, ErrUnauthorized
	}

	trades, err := s.tradeRepository.FindByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	var responses []dto.TradeResponse
	for _, trade := range trades {
		responses = append(responses, *s.toTradeResponse(&trade))
	}

	return responses, nil
}

// GetTradeSummary returns aggregate stats
func (s *TradeService) GetTradeSummary(ctx context.Context, userID string) (*dto.TradeSummary, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	summary, err := s.tradeRepository.GetTradeSummary(ctx, userObjectID)
	if err != nil {
		return nil, err
	}

	return &dto.TradeSummary{
		TotalTrades:     summary.TotalTrades,
		TotalBuyValue:   summary.TotalBuyValue,
		TotalSellValue:  summary.TotalSellValue,
		TotalCommission: summary.TotalCommission,
		NetProfit:       summary.NetProfit,
	}, nil
}

// calculateAvgPrice calculates weighted average price
func (s *TradeService) calculateAvgPrice(oldAvg, oldQty, newPrice, newQty float64) float64 {
	totalQty := oldQty + newQty
	if totalQty == 0 {
		return newPrice
	}
	return ((oldAvg * oldQty) + (newPrice * newQty)) / totalQty
}

// updatePosition updates position after trade
func (s *TradeService) updatePosition(ctx context.Context, trade *tradeModel.Trade, side orderModel.OrderSide) {
	position, err := s.portfolioRepository.FindPositionByPortfolioAndSymbol(ctx, trade.PortfolioID, trade.Symbol)

	if side == orderModel.OrderSideBuy {
		if err != nil {
			// Create new position
			newPosition := &portfolioModel.Position{
				PortfolioID:  trade.PortfolioID,
				InstrumentID: trade.InstrumentID,
				Symbol:       trade.Symbol,
				Quantity:     trade.Quantity,
				AvgCost:      trade.Price,
				TotalCost:    trade.Total,
			}
			s.portfolioRepository.CreatePosition(ctx, newPosition)

			// Create position lot
			s.portfolioRepository.CreatePositionLot(ctx, &portfolioModel.PositionLot{
				PortfolioID:  trade.PortfolioID,
				PositionID:   newPosition.ID,
				InstrumentID: trade.InstrumentID,
				TradeID:      trade.ID,
				Quantity:     trade.Quantity,
				RemainingQty: trade.Quantity,
				CostPerUnit:  trade.Price,
				PurchasedAt:  trade.ExecutedAt,
			})
		} else {
			// Update existing position
			newQty := position.Quantity + trade.Quantity
			newTotalCost := position.TotalCost + trade.Total
			newAvgCost := newTotalCost / newQty

			s.portfolioRepository.UpdatePosition(ctx, position.ID, bson.M{
				"quantity":  newQty,
				"totalCost": newTotalCost,
				"avgCost":   newAvgCost,
			})

			// Add position lot
			s.portfolioRepository.CreatePositionLot(ctx, &portfolioModel.PositionLot{
				PortfolioID:  trade.PortfolioID,
				PositionID:   position.ID,
				InstrumentID: trade.InstrumentID,
				TradeID:      trade.ID,
				Quantity:     trade.Quantity,
				RemainingQty: trade.Quantity,
				CostPerUnit:  trade.Price,
				PurchasedAt:  trade.ExecutedAt,
			})
		}
	} else {
		// SELL - reduce position using FIFO
		if position != nil {
			remainingToSell := trade.Quantity
			lots, _ := s.portfolioRepository.FindLotsByPositionID(ctx, position.ID)

			for _, lot := range lots {
				if remainingToSell <= 0 {
					break
				}

				if lot.RemainingQty <= remainingToSell {
					remainingToSell -= lot.RemainingQty
					s.portfolioRepository.UpdatePositionLot(ctx, lot.ID, 0)
				} else {
					newRemaining := lot.RemainingQty - remainingToSell
					s.portfolioRepository.UpdatePositionLot(ctx, lot.ID, newRemaining)
					remainingToSell = 0
				}
			}

			// Update position quantity
			newQty := position.Quantity - trade.Quantity
			if newQty <= 0 {
				s.portfolioRepository.DeletePosition(ctx, position.ID)
			} else {
				newTotalCost := position.AvgCost * newQty
				s.portfolioRepository.UpdatePosition(ctx, position.ID, bson.M{
					"quantity":  newQty,
					"totalCost": newTotalCost,
				})
			}
		}
	}
}

func (s *TradeService) toTradeResponse(trade *tradeModel.Trade) *dto.TradeResponse {
	return &dto.TradeResponse{
		ID:           trade.ID.Hex(),
		OrderID:      trade.OrderID.Hex(),
		UserID:       trade.UserID.Hex(),
		AccountID:    trade.AccountID.Hex(),
		PortfolioID:  trade.PortfolioID.Hex(),
		InstrumentID: trade.InstrumentID.Hex(),
		Symbol:       trade.Symbol,
		Side:         string(trade.Side),
		Quantity:     trade.Quantity,
		Price:        trade.Price,
		Total:        trade.Total,
		Commission:   trade.Commission,
		NetAmount:    trade.NetAmount,
		ExecutedAt:   trade.ExecutedAt.Format(time.RFC3339),
	}
}

// publishTradeEvents publishes WebSocket events for trade execution
func (s *TradeService) publishTradeEvents(
	trade *tradeModel.Trade,
	order *orderModel.Order,
	filledQty float64,
	avgPrice float64,
	status orderModel.OrderStatus,
) {
	userID := order.UserID.Hex()
	// Publish trade executed event
	ws.PublishTradeUpdate(userID, &ws.TradePayload{
		TradeID:    trade.ID.Hex(),
		OrderID:    trade.OrderID.Hex(),
		Symbol:     trade.Symbol,
		Side:       string(trade.Side),
		Quantity:   trade.Quantity,
		Price:      trade.Price,
		Commission: trade.Commission,
	})
	// Publish order update event
	ws.PublishOrderUpdate(userID, &ws.OrderPayload{
		OrderID:   order.ID.Hex(),
		Symbol:    order.Symbol,
		Side:      string(order.Side),
		Status:    string(status),
		FilledQty: filledQty,
		AvgPrice:  avgPrice,
	})
}
