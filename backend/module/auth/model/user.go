package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const UserCollection = "users"

type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	RoleID      primitive.ObjectID `bson:"roleId" json:"roleId"`
	Email       string             `bson:"email" json:"email"`
	Password    string             `bson:"password" json:"-"` // Never expose password
	FullName    string             `bson:"fullName" json:"fullName"`
	Phone       string             `bson:"phone,omitempty" json:"phone,omitempty"`
	LastLoginAt time.Time          `bson:"lastLoginAt" json:"lastLoginAt"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}
