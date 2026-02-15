package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Showtime struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	MovieID      primitive.ObjectID `bson:"movie_id" json:"movieId"`
	StartTime    time.Time          `bson:"start_time" json:"startTime"`
	AuditoriumID string             `bson:"auditorium_id" json:"auditoriumId"`
	SeatmapID    string             `bson:"seatmap_id" json:"seatmapId"`
	CreatedAt    time.Time          `bson:"created_at" json:"createdAt"`
}
