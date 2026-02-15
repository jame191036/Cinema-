package handlers

import (
	"context"
	"net/http"
	"time"

	"cinema-booking/internal/models"
	"cinema-booking/internal/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AdminHandler struct {
	Mongo *services.MongoService
}

func (h *AdminHandler) ListBookings(c *gin.Context) {
	ctx := context.Background()
	filter := bson.M{}

	if status := c.Query("status"); status != "" {
		filter["status"] = status
	}
	if userID := c.Query("user_id"); userID != "" {
		if oid, err := primitive.ObjectIDFromHex(userID); err == nil {
			filter["user_id"] = oid
		}
	}
	if dateStr := c.Query("date"); dateStr != "" {
		if t, err := time.Parse("2006-01-02", dateStr); err == nil {
			filter["created_at"] = bson.M{
				"$gte": t,
				"$lt":  t.Add(24 * time.Hour),
			}
		}
	}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(100)

	var bookings []models.Booking
	cursor, err := h.Mongo.Collection("bookings").Find(ctx, filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch bookings"})
		return
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &bookings); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode"})
		return
	}

	if bookings == nil {
		bookings = []models.Booking{}
	}
	c.JSON(http.StatusOK, bookings)
}

func (h *AdminHandler) ListAuditLogs(c *gin.Context) {
	ctx := context.Background()
	filter := bson.M{}

	if eventType := c.Query("event_type"); eventType != "" {
		filter["event_type"] = eventType
	}

	dateFilter := bson.M{}
	if from := c.Query("date_from"); from != "" {
		if t, err := time.Parse("2006-01-02", from); err == nil {
			dateFilter["$gte"] = t
		}
	}
	if to := c.Query("date_to"); to != "" {
		if t, err := time.Parse("2006-01-02", to); err == nil {
			dateFilter["$lt"] = t.Add(24 * time.Hour)
		}
	}
	if len(dateFilter) > 0 {
		filter["created_at"] = dateFilter
	}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(200)

	var logs []models.AuditLog
	cursor, err := h.Mongo.Collection("audit_logs").Find(ctx, filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch audit logs"})
		return
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &logs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode"})
		return
	}

	if logs == nil {
		logs = []models.AuditLog{}
	}
	c.JSON(http.StatusOK, logs)
}
