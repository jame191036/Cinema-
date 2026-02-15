package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Movie struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	DurationMin int                `bson:"duration_min" json:"durationMin"`
	Rating      string             `bson:"rating" json:"rating"`
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
}
