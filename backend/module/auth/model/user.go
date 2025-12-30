package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserCollection is the MongoDB collection name for users.
const UserCollection = "users"

// User represents a registered user in the system.
// Contains authentication info and user profile data.
type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	RoleID      primitive.ObjectID `bson:"roleId" json:"roleId"`
	Email       string             `bson:"email" json:"email"`
	Password    string             `bson:"password" json:"-"` // Never sent to client
	FullName    string             `bson:"fullName" json:"fullName"`
	Phone       string             `bson:"phone,omitempty" json:"phone,omitempty"`
	IsVerified  bool               `bson:"isVerified" json:"isVerified"`
	LastLoginAt time.Time          `bson:"lastLoginAt" json:"lastLoginAt"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// NewUser creates a user with default timestamps.
func NewUser(email, passwordHash, fullName string, roleID primitive.ObjectID) *User {
	now := time.Now()
	return &User{
		ID:         primitive.NewObjectID(),
		RoleID:     roleID,
		Email:      email,
		Password:   passwordHash,
		FullName:   fullName,
		IsVerified: false,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}
