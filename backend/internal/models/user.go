package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	RoleUser  = "USER"
	RoleAdmin = "ADMIN"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AuthProvider string             `bson:"auth_provider" json:"authProvider"`
	ProviderSub  string             `bson:"provider_sub" json:"providerSub"`
	Email        string             `bson:"email" json:"email"`
	Name         string             `bson:"name" json:"name"`
	Role         string             `bson:"role" json:"role"`
	CreatedAt    time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updatedAt"`
}
