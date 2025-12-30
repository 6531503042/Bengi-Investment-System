package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AccountCollection is the MongoDB collection name for trading accounts.
const AccountCollection = "accounts"

// AccountStatus represents the current state of an account.
type AccountStatus string

const (
	AccountStatusActive AccountStatus = "ACTIVE" // Normal trading enabled
	AccountStatusFrozen AccountStatus = "FROZEN" // Temporarily suspended
	AccountStatusClosed AccountStatus = "CLOSED" // Permanently closed
)

// AccountType distinguishes between demo and real trading accounts.
type AccountType string

const (
	AccountTypeDemo AccountType = "DEMO" // Paper trading with virtual money
	AccountTypeLive AccountType = "LIVE" // Real trading with actual funds
)

// Account represents a user's trading account.
// Users can have both demo and live accounts.
type Account struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID         primitive.ObjectID `bson:"userId" json:"userId"`
	Currency       string             `bson:"currency" json:"currency"` // USD, THB, etc.
	Balance        float64            `bson:"balance" json:"balance"`   // Current available balance
	Status         AccountStatus      `bson:"status" json:"status"`
	Type           AccountType        `bson:"type" json:"type"`
	Leverage       int                `bson:"leverage" json:"leverage"`             // Max leverage (1 = no leverage)
	InitialBalance float64            `bson:"initialBalance" json:"initialBalance"` // Starting balance (for reset)
	TotalDeposits  float64            `bson:"totalDeposits" json:"totalDeposits"`
	TotalPnL       float64            `bson:"totalPnL" json:"totalPnL"` // Cumulative profit/loss
	CreatedAt      time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// NewDemoAccount creates a demo account with $50,000 virtual balance.
func NewDemoAccount(userID primitive.ObjectID, currency string) *Account {
	now := time.Now()
	initialBalance := 50000.0
	return &Account{
		ID:             primitive.NewObjectID(),
		UserID:         userID,
		Currency:       currency,
		Balance:        initialBalance,
		Status:         AccountStatusActive,
		Type:           AccountTypeDemo,
		Leverage:       10,
		InitialBalance: initialBalance,
		TotalDeposits:  0,
		TotalPnL:       0,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// ResetBalance resets a demo account to its initial balance.
func (a *Account) ResetBalance() {
	a.Balance = a.InitialBalance
	a.TotalDeposits = 0
	a.TotalPnL = 0
	a.UpdatedAt = time.Now()
}
