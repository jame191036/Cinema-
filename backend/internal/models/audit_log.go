package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuditLog struct {
	ID         primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	EventType  string              `bson:"event_type" json:"eventType"`
	UserID     *primitive.ObjectID `bson:"user_id,omitempty" json:"userId,omitempty"`
	ShowtimeID *primitive.ObjectID `bson:"showtime_id,omitempty" json:"showtimeId,omitempty"`
	SeatCode   string              `bson:"seat_code,omitempty" json:"seatCode,omitempty"`
	BookingID  *primitive.ObjectID `bson:"booking_id,omitempty" json:"bookingId,omitempty"`
	Payload    map[string]any      `bson:"payload,omitempty" json:"payload,omitempty"`
	CreatedAt  time.Time           `bson:"created_at" json:"createdAt"`
}
