package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const AccountCollection = "accounts"

type AccountStatus string
type AccountType string

const (
	AccountStatusActive AccountStatus = "ACTIVE"
	AccountStatusFrozen AccountStatus = "FROZEN"
	AccountStatusClosed AccountStatus = "CLOSED"
)

const (
	AccountTypeDemo AccountType = "DEMO"
	AccountTypeLive AccountType = "LIVE"
)

type Account struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID         primitive.ObjectID `bson:"userId" json:"userId"`
	Currency       string             `bson:"currency" json:"currency"`
	Balance        float64            `bson:"balance" json:"balance"`
	Status         AccountStatus      `bson:"status" json:"status"`
	Type           AccountType        `bson:"type" json:"type"`                     // DEMO or LIVE
	Leverage       int                `bson:"leverage" json:"leverage"`             // Default leverage (1-100)
	InitialBalance float64            `bson:"initialBalance" json:"initialBalance"` // Starting balance for demo
	TotalDeposits  float64            `bson:"totalDeposits" json:"totalDeposits"`   // Total deposited (demo)
	TotalPnL       float64            `bson:"totalPnL" json:"totalPnL"`             // Cumulative P&L
	CreatedAt      time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt" json:"updatedAt"`
}
