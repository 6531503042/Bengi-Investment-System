package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const OrderCollection = "orders"

type OrderSide string
type OrderType string
type OrderStatus string
type TimeInForce string

const (
	OrderSideBuy  OrderSide = "BUY"
	OrderSideSell OrderSide = "SELL"
)

const (
	OrderTypeMarket OrderType = "MARKET"
	OrderTypeLimit  OrderType = "LIMIT"
	OrderTypeStop   OrderType = "STOP"
)

const (
	OrderStatusPending         OrderStatus = "PENDING"
	OrderStatusOpen            OrderStatus = "OPEN"
	OrderStatusPartiallyFilled OrderStatus = "PARTIALLY_FILLED"
	OrderStatusFilled          OrderStatus = "FILLED"
	OrderStatusCancelled       OrderStatus = "CANCELLED"
	OrderStatusRejected        OrderStatus = "REJECTED"
	OrderStatusExpired         OrderStatus = "EXPIRED"
)

const (
	TimeInForceGTC TimeInForce = "GTC" // Good Till Cancelled
	TimeInForceDay TimeInForce = "DAY" // Day Order
	TimeInForceIOC TimeInForce = "IOC" // Immediate or Cancel
	TimeInForceFOK TimeInForce = "FOK" // Fill or Kill
)

type Order struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID       primitive.ObjectID `bson:"userId" json:"userId"`
	AccountID    primitive.ObjectID `bson:"accountId" json:"accountId"`
	PortfolioID  primitive.ObjectID `bson:"portfolioId" json:"portfolioId"`
	InstrumentID primitive.ObjectID `bson:"instrumentId" json:"instrumentId"`
	Symbol       string             `bson:"symbol" json:"symbol"`
	Side         OrderSide          `bson:"side" json:"side"`
	Type         OrderType          `bson:"type" json:"type"`
	Status       OrderStatus        `bson:"status" json:"status"`
	TimeInForce  TimeInForce        `bson:"timeInForce" json:"timeInForce"`
	Quantity     float64            `bson:"quantity" json:"quantity"`
	FilledQty    float64            `bson:"filledQty" json:"filledQty"`
	Price        float64            `bson:"price,omitempty" json:"price,omitempty"`
	StopPrice    float64            `bson:"stopPrice,omitempty" json:"stopPrice,omitempty"`
	AvgFillPrice float64            `bson:"avgFillPrice,omitempty" json:"avgFillPrice,omitempty"`
	Commission   float64            `bson:"commission" json:"commission"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt" json:"updatedAt"`
	FilledAt     *time.Time         `bson:"filledAt,omitempty" json:"filledAt,omitempty"`
	CancelledAt  *time.Time         `bson:"cancelledAt,omitempty" json:"cancelledAt,omitempty"`
}
