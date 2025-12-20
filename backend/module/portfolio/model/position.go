package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const PositionCollection = "positions"

type Position struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PortfolioID  primitive.ObjectID `bson:"portfolioId" json:"portfolioId"`
	InstrumentID primitive.ObjectID `bson:"instrumentId" json:"instrumentId"`
	Symbol       string             `bson:"symbol" json:"symbol"`
	Quantity     float64            `bson:"quantity" json:"quantity"`
	AvgCost      float64            `bson:"avgCost" json:"avgCost"`
	TotalCost    float64            `bson:"totalCost" json:"totalCost"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt" json:"updatedAt"`
}
