package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	PaymentStatusPending = "PENDING"
	PaymentStatusSuccess = "SUCCESS"
	PaymentStatusFailed  = "FAILED"
)

type Payment struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BookingID primitive.ObjectID `bson:"booking_id" json:"bookingId"`
	Amount    float64            `bson:"amount" json:"amount"`
	Status    string             `bson:"status" json:"status"`
	Provider  string             `bson:"provider" json:"provider"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updatedAt"`
}
