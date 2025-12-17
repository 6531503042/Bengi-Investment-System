package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const RoleCollection = "roles"

// Role names
const (
	RoleAdmin  = "ADMIN"
	RoleUser   = "USER"
	RoleTrader = "TRADER"
)

type Role struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Permissions []string           `bson:"permissions" json:"permissions"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}
