package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PortfolioCollection is the MongoDB collection name for portfolios.
const PortfolioCollection = "portfolios"

// Portfolio represents a user's investment portfolio.
// A user can have multiple portfolios (e.g., "Long-term", "Day Trading").
type Portfolio struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      primitive.ObjectID `bson:"userId" json:"userId"`
	AccountID   primitive.ObjectID `bson:"accountId" json:"accountId"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	IsDefault   bool               `bson:"isDefault" json:"isDefault"` // Primary portfolio for this account
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// NewPortfolio creates a portfolio with default timestamps.
func NewPortfolio(userID, accountID primitive.ObjectID, name string, isDefault bool) *Portfolio {
	now := time.Now()
	return &Portfolio{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		AccountID: accountID,
		Name:      name,
		IsDefault: isDefault,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
