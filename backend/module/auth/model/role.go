package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const RoleCollection = "roles"

const (
	RoleAdmin = "ADMIN"
	RoleUser  = "USER"
)

type Role struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}
