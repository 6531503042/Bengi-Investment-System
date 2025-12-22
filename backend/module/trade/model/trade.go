package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const TradeCollection = "trades"

type TradeSide string

const (
	TradeSideBuy  TradeSide = "BUY"
	TradeSideSell TradeSide = "SELL"
)

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
	Total        float64            `bson:"total" json:"total"` // Quantity * Price
	Commission   float64            `bson:"commission" json:"commission"`
	NetAmount    float64            `bson:"netAmount" json:"netAmount"` // Total +/- Commission
	ExecutedAt   time.Time          `bson:"executedAt" json:"executedAt"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
}
