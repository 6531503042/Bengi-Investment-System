package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PositionCollection is the MongoDB collection name for portfolio positions.
const PositionCollection = "positions"

// Position represents a holding of a specific instrument in a portfolio.
// Tracks quantity, average cost, and total invested amount.
type Position struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PortfolioID  primitive.ObjectID `bson:"portfolioId" json:"portfolioId"`
	InstrumentID primitive.ObjectID `bson:"instrumentId" json:"instrumentId"`
	Symbol       string             `bson:"symbol" json:"symbol"`
	Quantity     float64            `bson:"quantity" json:"quantity"`   // Number of shares/units held
	AvgCost      float64            `bson:"avgCost" json:"avgCost"`     // Average purchase price per unit
	TotalCost    float64            `bson:"totalCost" json:"totalCost"` // Total amount invested
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// NewPosition creates a new position with calculated totals.
func NewPosition(portfolioID, instrumentID primitive.ObjectID, symbol string, qty, price float64) *Position {
	now := time.Now()
	return &Position{
		ID:           primitive.NewObjectID(),
		PortfolioID:  portfolioID,
		InstrumentID: instrumentID,
		Symbol:       symbol,
		Quantity:     qty,
		AvgCost:      price,
		TotalCost:    qty * price,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// AddShares adds more shares to an existing position.
// Recalculates average cost using weighted average formula.
func (p *Position) AddShares(qty, price float64) {
	newTotal := p.TotalCost + (qty * price)
	newQty := p.Quantity + qty
	p.Quantity = newQty
	p.AvgCost = newTotal / newQty
	p.TotalCost = newTotal
	p.UpdatedAt = time.Now()
}

// RemoveShares removes shares from the position.
// Returns error if trying to remove more than available.
func (p *Position) RemoveShares(qty float64) error {
	if qty > p.Quantity {
		return ErrInsufficientShares
	}
	p.Quantity -= qty
	p.TotalCost = p.Quantity * p.AvgCost
	p.UpdatedAt = time.Now()
	return nil
}

// MarketValue returns the current market value based on given price.
func (p *Position) MarketValue(currentPrice float64) float64 {
	return p.Quantity * currentPrice
}

// UnrealizedPnL calculates unrealized profit/loss.
func (p *Position) UnrealizedPnL(currentPrice float64) float64 {
	return p.MarketValue(currentPrice) - p.TotalCost
}

// UnrealizedPnLPercent calculates P&L as a percentage.
func (p *Position) UnrealizedPnLPercent(currentPrice float64) float64 {
	if p.TotalCost == 0 {
		return 0
	}
	return (p.UnrealizedPnL(currentPrice) / p.TotalCost) * 100
}

// IsEmpty returns true if there are no shares in the position.
func (p *Position) IsEmpty() bool {
	return p.Quantity <= 0
}
