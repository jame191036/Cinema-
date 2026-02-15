package handlers

import (
	"context"
	"net/http"

	"cinema-booking/internal/models"
	"cinema-booking/internal/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShowtimeHandler struct {
	Mongo *services.MongoService
}

func (h *ShowtimeHandler) ListShowtimes(c *gin.Context) {
	filter := bson.M{}
	if movieID := c.Query("movie_id"); movieID != "" {
		oid, err := primitive.ObjectIDFromHex(movieID)
		if err == nil {
			filter["movie_id"] = oid
		}
	}

	var showtimes []models.Showtime
	cursor, err := h.Mongo.Collection("showtimes").Find(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch showtimes"})
		return
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &showtimes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode"})
		return
	}

	if showtimes == nil {
		showtimes = []models.Showtime{}
	}
	c.JSON(http.StatusOK, showtimes)
}

func (h *ShowtimeHandler) GetSeats(c *gin.Context) {
	showtimeID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid showtime id"})
		return
	}

	var seats []models.SeatReservation
	cursor, err := h.Mongo.Collection("seat_reservations").Find(
		context.Background(),
		bson.M{"showtime_id": showtimeID},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch seats"})
		return
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &seats); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode"})
		return
	}

	if seats == nil {
		seats = []models.SeatReservation{}
	}
	c.JSON(http.StatusOK, seats)
}
