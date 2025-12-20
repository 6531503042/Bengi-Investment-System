package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const PositionLotCollection = "positionLots"

type PositionLot struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PortfolioID  primitive.ObjectID `bson:"portfolioId" json:"portfolioId"`
	PositionID   primitive.ObjectID `bson:"positionId" json:"positionId"`
	InstrumentID primitive.ObjectID `bson:"instrumentId" json:"instrumentId"`
	TradeID      primitive.ObjectID `bson:"tradeId" json:"tradeId"`
	Quantity     float64            `bson:"quantity" json:"quantity"`
	RemainingQty float64            `bson:"remainingQty" json:"remainingQty"`
	CostPerUnit  float64            `bson:"costPerUnit" json:"costPerUnit"`
	PurchasedAt  time.Time          `bson:"purchasedAt" json:"purchasedAt"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
}
