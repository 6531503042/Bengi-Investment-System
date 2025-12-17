package dto

import "time"

type (
	CreateAccountRequest struct {
		Currency string `json:"currency" validate:"required, len=3, oneof=USD RUB EUR"`
	}

	AccountResponse struct {
		ID        string    `json:"id"`
		Currency  string    `json:"currency"`
		Balance   float64   `json:"balance"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}

	DepositRequest struct {
		Amount      float64 `json:"amount" validate:"required,gt=0"`
		Description string  `json:"description"`
	}

	WithdrawRequest struct {
		Amount      float64 `json:"amount" validate:"required,gt=0"`
		Description string  `json:"description"`
	}

	TransferRequest struct {
		ToAccountID string  `json:"toAccountId" validate:"required"`
		Amount      float64 `json:"amount" validate:"required,gt=0"`
		Description string  `json:"description"`
	}
)
