package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const PortfolioCollection = "portfolios"

type Portfolio struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"userId" json:"userId"`
	AccountID primitive.ObjectID `bson:"accountId" json:"accountId"`
	Name      string             `bson:"name" json:"name"`
	IsDefault bool               `bson:"isDefault" json:"isDefault"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}
