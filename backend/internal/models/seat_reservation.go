package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	SeatStateAvailable = "AVAILABLE"
	SeatStateLocked    = "LOCKED"
	SeatStateBooked    = "BOOKED"
)

type SeatReservation struct {
	ID             primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	ShowtimeID     primitive.ObjectID  `bson:"showtime_id" json:"showtimeId"`
	SeatCode       string              `bson:"seat_code" json:"seatCode"`
	State          string              `bson:"state" json:"state"`
	LockedByUserID *primitive.ObjectID `bson:"locked_by_user_id,omitempty" json:"lockedByUserId,omitempty"`
	LockExpiresAt  *time.Time          `bson:"lock_expires_at,omitempty" json:"lockExpiresAt,omitempty"`
	BookingID      *primitive.ObjectID `bson:"booking_id,omitempty" json:"bookingId,omitempty"`
	UpdatedAt      time.Time           `bson:"updated_at" json:"updatedAt"`
}
