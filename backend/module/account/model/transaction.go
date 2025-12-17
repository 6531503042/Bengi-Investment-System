package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const TransactionCollection = "transactions"

type TransactionType string
type TransactionStatus string

const (
	TransactionTypeDeposit  TransactionType = "DEPOSIT"
	TransactionTypeWithdraw TransactionType = "WITHDRAW"
	TransactionTypeTrade    TransactionType = "TRADE"
	TransactionTypeFee      TransactionType = "FEE"
	TransactionTypeDividend TransactionType = "DIVIDEND"
	TransactionTypeTransfer TransactionType = "TRANSFER"
)

const (
	TransactionStatusPending   TransactionStatus = "PENDING"
	TransactionStatusCompleted TransactionStatus = "COMPLETED"
	TransactionStatusFailed    TransactionStatus = "FAILED"
	TransactionStatusCancelled TransactionStatus = "CANCELLED"
)

type Transaction struct {
	ID            primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	AccountID     primitive.ObjectID  `bson:"accountId" json:"accountId"`
	Type          TransactionType     `bson:"type" json:"type"`
	Status        TransactionStatus   `bson:"status" json:"status"`
	Amount        float64             `bson:"amount" json:"amount"`
	BalanceBefore float64             `bson:"balanceBefore" json:"balanceBefore"`
	BalanceAfter  float64             `bson:"balanceAfter" json:"balanceAfter"`
	ReferenceType string              `bson:"referenceType,omitempty" json:"referenceType,omitempty"`
	ReferenceID   *primitive.ObjectID `bson:"referenceId,omitempty" json:"referenceId,omitempty"`
	Description   string              `bson:"description" json:"description"`
	CreatedAt     time.Time           `bson:"createdAt" json:"createdAt"`
}
