package dto

import "time"

// DemoDepositRequest is the request for depositing demo funds
type DemoDepositRequest struct {
	Amount float64 `json:"amount" validate:"required,min=100,max=1000000"`
}

// DemoDepositResponse is the response after depositing demo funds
type DemoDepositResponse struct {
	AccountID     string  `json:"accountId"`
	NewBalance    float64 `json:"newBalance"`
	TotalDeposits float64 `json:"totalDeposits"`
	Message       string  `json:"message"`
}

// DemoResetRequest is the request for resetting demo account
type DemoResetRequest struct {
	InitialBalance float64 `json:"initialBalance,omitempty"` // Default 10000
}

// DemoResetResponse is the response after resetting demo account
type DemoResetResponse struct {
	AccountID      string  `json:"accountId"`
	NewBalance     float64 `json:"newBalance"`
	InitialBalance float64 `json:"initialBalance"`
	Message        string  `json:"message"`
}

// CreateDemoAccountRequest is for creating a new demo account
type CreateDemoAccountRequest struct {
	Currency       string  `json:"currency,omitempty"`       // Default USD
	Leverage       int     `json:"leverage,omitempty"`       // Default 10
	InitialBalance float64 `json:"initialBalance,omitempty"` // Default 10000
}

// CreateDemoAccountResponse is the response after creating demo account
type CreateDemoAccountResponse struct {
	AccountID      string  `json:"accountId"`
	Currency       string  `json:"currency"`
	Balance        float64 `json:"balance"`
	Leverage       int     `json:"leverage"`
	InitialBalance float64 `json:"initialBalance"`
	Message        string  `json:"message"`
}

// DemoAccountStats shows demo account statistics
type DemoAccountStats struct {
	AccountID      string    `json:"accountId"`
	Currency       string    `json:"currency"`
	Balance        float64   `json:"balance"`
	InitialBalance float64   `json:"initialBalance"`
	TotalDeposits  float64   `json:"totalDeposits"`
	TotalPnL       float64   `json:"totalPnL"`
	Leverage       int       `json:"leverage"`
	PnLPercentage  float64   `json:"pnlPercentage"`
	CreatedAt      time.Time `json:"createdAt"`
}
