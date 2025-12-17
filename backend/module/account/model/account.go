package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const AccountCollection = "accounts"

type AccountStatus string

const (
	AccountStatusActive AccountStatus = "ACTIVE"
	AccountStatusFrozen AccountStatus = "FROZEN"
	AccountStatusClosed AccountStatus = "CLOSED"
)

type Account struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"userId" json:"userId"`
	Currency  string             `bson:"currency" json:"currency"`
	Balance   float64            `bson:"balance" json:"balance"`
	Status    AccountStatus      `bson:"status" json:"status"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}
