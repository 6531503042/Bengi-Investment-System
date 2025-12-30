package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TradeCollection is the MongoDB collection name for executed trades.
const TradeCollection = "trades"

// TradeSide represents the direction of a trade.
type TradeSide string

const (
	TradeSideBuy  TradeSide = "BUY"
	TradeSideSell TradeSide = "SELL"
)

// String returns the string representation of TradeSide.
func (s TradeSide) String() string {
	return string(s)
}

// IsBuy returns true if this is a buy trade.
func (s TradeSide) IsBuy() bool {
	return s == TradeSideBuy
}

// Trade represents a single executed trade.
// Created when an order is filled (partially or fully).
type Trade struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OrderID      primitive.ObjectID `bson:"orderId" json:"orderId"`
	UserID       primitive.ObjectID `bson:"userId" json:"userId"`
	AccountID    primitive.ObjectID `bson:"accountId" json:"accountId"`
	PortfolioID  primitive.ObjectID `bson:"portfolioId" json:"portfolioId"`
	InstrumentID primitive.ObjectID `bson:"instrumentId" json:"instrumentId"`
	Symbol       string             `bson:"symbol" json:"symbol"`
	Side         TradeSide          `bson:"side" json:"side"`
	Quantity     float64            `bson:"quantity" json:"quantity"`
	Price        float64            `bson:"price" json:"price"`
	Total        float64            `bson:"total" json:"total"`           // Quantity * Price
	Commission   float64            `bson:"commission" json:"commission"` // Trading fee
	NetAmount    float64            `bson:"netAmount" json:"netAmount"`   // Total +/- Commission based on side
	ExecutedAt   time.Time          `bson:"executedAt" json:"executedAt"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
}

// NewTrade creates a trade with calculated totals.
// Commission is subtracted for buys, added for sells.
func NewTrade(
	orderID, userID, accountID, portfolioID, instrumentID primitive.ObjectID,
	symbol string,
	side TradeSide,
	quantity, price, commissionRate float64,
) *Trade {
	total := quantity * price
	commission := total * commissionRate

	netAmount := total
	if side.IsBuy() {
		netAmount = total + commission // Pay more when buying
	} else {
		netAmount = total - commission // Receive less when selling
	}

	now := time.Now()
	return &Trade{
		ID:           primitive.NewObjectID(),
		OrderID:      orderID,
		UserID:       userID,
		AccountID:    accountID,
		PortfolioID:  portfolioID,
		InstrumentID: instrumentID,
		Symbol:       symbol,
		Side:         side,
		Quantity:     quantity,
		Price:        price,
		Total:        total,
		Commission:   commission,
		NetAmount:    netAmount,
		ExecutedAt:   now,
		CreatedAt:    now,
	}
}
