package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const LeveragePositionCollection = "leverage_positions"

type PositionSide string
type PositionStatus string

const (
	PositionSideLong  PositionSide = "LONG"
	PositionSideShort PositionSide = "SHORT"
)

const (
	PositionStatusOpen       PositionStatus = "OPEN"
	PositionStatusClosed     PositionStatus = "CLOSED"
	PositionStatusLiquidated PositionStatus = "LIQUIDATED"
)

// LeveragePosition represents a leveraged trading position
type LeveragePosition struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID           primitive.ObjectID `bson:"userId" json:"userId"`
	AccountID        primitive.ObjectID `bson:"accountId" json:"accountId"`
	Symbol           string             `bson:"symbol" json:"symbol"`
	Side             PositionSide       `bson:"side" json:"side"`         // LONG or SHORT
	Leverage         int                `bson:"leverage" json:"leverage"` // 1-100x
	EntryPrice       float64            `bson:"entryPrice" json:"entryPrice"`
	CurrentPrice     float64            `bson:"currentPrice" json:"currentPrice"`
	Quantity         float64            `bson:"quantity" json:"quantity"`
	Margin           float64            `bson:"margin" json:"margin"`                   // Collateral used
	StopLoss         *float64           `bson:"stopLoss,omitempty" json:"stopLoss"`     // Optional
	TakeProfit       *float64           `bson:"takeProfit,omitempty" json:"takeProfit"` // Optional
	LiquidationPrice float64            `bson:"liquidationPrice" json:"liquidationPrice"`
	UnrealizedPnL    float64            `bson:"unrealizedPnL" json:"unrealizedPnL"`
	RealizedPnL      float64            `bson:"realizedPnL" json:"realizedPnL"`
	Status           PositionStatus     `bson:"status" json:"status"`
	OpenedAt         time.Time          `bson:"openedAt" json:"openedAt"`
	ClosedAt         *time.Time         `bson:"closedAt,omitempty" json:"closedAt"`
	UpdatedAt        time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// CalculateUnrealizedPnL calculates the unrealized P&L based on current price
func (p *LeveragePosition) CalculateUnrealizedPnL(currentPrice float64) float64 {
	p.CurrentPrice = currentPrice

	if p.Side == PositionSideLong {
		// Long: profit when price goes up
		p.UnrealizedPnL = (currentPrice - p.EntryPrice) * p.Quantity
	} else {
		// Short: profit when price goes down
		p.UnrealizedPnL = (p.EntryPrice - currentPrice) * p.Quantity
	}

	// Apply leverage effect (P&L is already magnified by position size)
	return p.UnrealizedPnL
}

// CalculateLiquidationPrice calculates the price at which position gets liquidated
func (p *LeveragePosition) CalculateLiquidationPrice() float64 {
	// Liquidation occurs when loss equals margin (simplified)
	marginRatio := 1.0 / float64(p.Leverage)

	if p.Side == PositionSideLong {
		// Long liquidates when price drops
		p.LiquidationPrice = p.EntryPrice * (1 - marginRatio*0.9) // 90% of margin = liquidation
	} else {
		// Short liquidates when price rises
		p.LiquidationPrice = p.EntryPrice * (1 + marginRatio*0.9)
	}

	return p.LiquidationPrice
}

// IsLiquidated checks if position should be liquidated at current price
func (p *LeveragePosition) IsLiquidated(currentPrice float64) bool {
	if p.Side == PositionSideLong {
		return currentPrice <= p.LiquidationPrice
	}
	return currentPrice >= p.LiquidationPrice
}

// ShouldTriggerStopLoss checks if stop loss should trigger
func (p *LeveragePosition) ShouldTriggerStopLoss(currentPrice float64) bool {
	if p.StopLoss == nil {
		return false
	}
	if p.Side == PositionSideLong {
		return currentPrice <= *p.StopLoss
	}
	return currentPrice >= *p.StopLoss
}

// ShouldTriggerTakeProfit checks if take profit should trigger
func (p *LeveragePosition) ShouldTriggerTakeProfit(currentPrice float64) bool {
	if p.TakeProfit == nil {
		return false
	}
	if p.Side == PositionSideLong {
		return currentPrice >= *p.TakeProfit
	}
	return currentPrice <= *p.TakeProfit
}
