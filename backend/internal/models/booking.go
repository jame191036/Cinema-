package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	BookingStatusLocked    = "LOCKED"
	BookingStatusBooked    = "BOOKED"
	BookingStatusCancelled = "CANCELLED"
	BookingStatusExpired   = "EXPIRED"
)

type Booking struct {
	ID            primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	UserID        primitive.ObjectID  `bson:"user_id" json:"userId"`
	ShowtimeID    primitive.ObjectID  `bson:"showtime_id" json:"showtimeId"`
	Seats         []string            `bson:"seats" json:"seats"`
	Status        string              `bson:"status" json:"status"`
	LockExpiresAt *time.Time          `bson:"lock_expires_at,omitempty" json:"lockExpiresAt,omitempty"`
	PaymentID     *primitive.ObjectID `bson:"payment_id,omitempty" json:"paymentId,omitempty"`
	CreatedAt     time.Time           `bson:"created_at" json:"createdAt"`
	UpdatedAt     time.Time           `bson:"updated_at" json:"updatedAt"`
}
