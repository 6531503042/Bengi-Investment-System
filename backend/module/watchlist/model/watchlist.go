package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const WatchlistCollection = "watchlists"

type Watchlist struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"userId" json:"userId"`
	Name      string             `bson:"name" json:"name"`
	Symbols   []string           `bson:"symbols" json:"symbols"`
	IsDefault bool               `bson:"isDefault" json:"isDefault"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}
