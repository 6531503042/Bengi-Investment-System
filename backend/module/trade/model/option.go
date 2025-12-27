package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const OptionCollection = "options"

type OptionType string
type OptionStatus string

const (
	OptionTypeCall OptionType = "CALL" // Bet price goes UP
	OptionTypePut  OptionType = "PUT"  // Bet price goes DOWN
)

const (
	OptionStatusOpen    OptionStatus = "OPEN"
	OptionStatusWon     OptionStatus = "WON"
	OptionStatusLost    OptionStatus = "LOST"
	OptionStatusExpired OptionStatus = "EXPIRED"
)

// Option represents a binary option trade
type Option struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      primitive.ObjectID `bson:"userId" json:"userId"`
	AccountID   primitive.ObjectID `bson:"accountId" json:"accountId"`
	Symbol      string             `bson:"symbol" json:"symbol"`
	OptionType  OptionType         `bson:"optionType" json:"optionType"`             // CALL or PUT
	StrikePrice float64            `bson:"strikePrice" json:"strikePrice"`           // Price at time of purchase
	Investment  float64            `bson:"investment" json:"investment"`             // Amount invested
	PayoutRate  float64            `bson:"payoutRate" json:"payoutRate"`             // e.g., 0.85 for 85% payout
	Payout      float64            `bson:"payout" json:"payout"`                     // Investment * (1 + PayoutRate) if won
	ExpiryTime  time.Time          `bson:"expiryTime" json:"expiryTime"`             // When option expires
	ExpiryPrice *float64           `bson:"expiryPrice,omitempty" json:"expiryPrice"` // Price at expiry
	Status      OptionStatus       `bson:"status" json:"status"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	SettledAt   *time.Time         `bson:"settledAt,omitempty" json:"settledAt"`
}

// DefaultPayoutRate is the default payout for winning options (85%)
const DefaultPayoutRate = 0.85

// ExpiryDurations available for options
var ExpiryDurations = map[string]time.Duration{
	"1m":  1 * time.Minute,
	"5m":  5 * time.Minute,
	"15m": 15 * time.Minute,
	"30m": 30 * time.Minute,
	"1h":  1 * time.Hour,
	"4h":  4 * time.Hour,
	"1d":  24 * time.Hour,
}

// CalculatePayout calculates the payout if option wins
func (o *Option) CalculatePayout() float64 {
	o.Payout = o.Investment * (1 + o.PayoutRate)
	return o.Payout
}

// Settle settles the option based on expiry price
func (o *Option) Settle(expiryPrice float64) {
	now := time.Now()
	o.ExpiryPrice = &expiryPrice
	o.SettledAt = &now

	// Determine if won or lost
	if o.OptionType == OptionTypeCall {
		// CALL wins if price goes UP
		if expiryPrice > o.StrikePrice {
			o.Status = OptionStatusWon
		} else {
			o.Status = OptionStatusLost
		}
	} else {
		// PUT wins if price goes DOWN
		if expiryPrice < o.StrikePrice {
			o.Status = OptionStatusWon
		} else {
			o.Status = OptionStatusLost
		}
	}
}

// IsExpired checks if option has expired
func (o *Option) IsExpired() bool {
	return time.Now().After(o.ExpiryTime)
}

// GetResult returns the profit/loss from this option
func (o *Option) GetResult() float64 {
	switch o.Status {
	case OptionStatusWon:
		return o.Payout - o.Investment // Net profit
	case OptionStatusLost:
		return -o.Investment // Lost everything
	default:
		return 0
	}
}
