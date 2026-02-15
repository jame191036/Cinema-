package worker

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"cinema-booking/internal/models"
	"cinema-booking/internal/services"
	wsHub "cinema-booking/internal/ws"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TimeoutWorker struct {
	Mongo    *services.MongoService
	Redis    *services.RedisService
	Hub      *wsHub.Hub
	Interval time.Duration
}

func NewTimeoutWorker(mongo *services.MongoService, redis *services.RedisService, hub *wsHub.Hub, interval time.Duration) *TimeoutWorker {
	return &TimeoutWorker{
		Mongo:    mongo,
		Redis:    redis,
		Hub:      hub,
		Interval: interval,
	}
}

func (w *TimeoutWorker) Start() {
	go func() {
		ticker := time.NewTicker(w.Interval)
		defer ticker.Stop()

		log.Printf("Timeout worker started (interval: %v)", w.Interval)
		for range ticker.C {
			w.cleanup()
		}
	}()
}

func (w *TimeoutWorker) cleanup() {
	ctx := context.Background()

	// Find expired locked seats
	filter := bson.M{
		"state":           models.SeatStateLocked,
		"lock_expires_at": bson.M{"$lt": time.Now()},
	}

	cursor, err := w.Mongo.Collection("seat_reservations").Find(ctx, filter)
	if err != nil {
		log.Printf("Worker: failed to query expired locks: %v", err)
		return
	}
	defer cursor.Close(ctx)

	var expiredSeats []models.SeatReservation
	if err := cursor.All(ctx, &expiredSeats); err != nil {
		log.Printf("Worker: failed to decode: %v", err)
		return
	}

	if len(expiredSeats) == 0 {
		return
	}

	log.Printf("Worker: found %d expired seat locks", len(expiredSeats))

	// Group by booking for batch processing
	bookingSeats := make(map[primitive.ObjectID][]models.SeatReservation)
	for _, seat := range expiredSeats {
		if seat.BookingID != nil {
			bookingSeats[*seat.BookingID] = append(bookingSeats[*seat.BookingID], seat)
		}

		// Release seat
		w.Mongo.Collection("seat_reservations").UpdateOne(ctx,
			bson.M{"_id": seat.ID, "state": models.SeatStateLocked},
			bson.M{"$set": bson.M{
				"state":             models.SeatStateAvailable,
				"locked_by_user_id": nil,
				"lock_expires_at":   nil,
				"booking_id":        nil,
				"updated_at":        time.Now(),
			}},
		)

		// Force release Redis lock
		w.Redis.ForceReleaseLock(ctx, seat.ShowtimeID.Hex(), seat.SeatCode)

		// Broadcast SEAT_RELEASED
		msg, _ := json.Marshal(map[string]interface{}{
			"type":     "SEAT_RELEASED",
			"seatCode": seat.SeatCode,
		})
		w.Hub.BroadcastToRoom(seat.ShowtimeID.Hex(), msg)

		// Create audit log
		auditLog := models.AuditLog{
			ID:         primitive.NewObjectID(),
			EventType:  "SEAT_RELEASED",
			ShowtimeID: &seat.ShowtimeID,
			SeatCode:   seat.SeatCode,
			BookingID:  seat.BookingID,
			Payload: map[string]interface{}{
				"reason": "lock_expired",
			},
			CreatedAt: time.Now(),
		}
		if seat.LockedByUserID != nil {
			auditLog.UserID = seat.LockedByUserID
		}
		w.Mongo.Collection("audit_logs").InsertOne(ctx, auditLog)
	}

	// Update expired bookings
	for bookingID := range bookingSeats {
		w.Mongo.Collection("bookings").UpdateOne(ctx,
			bson.M{"_id": bookingID, "status": models.BookingStatusLocked},
			bson.M{"$set": bson.M{"status": models.BookingStatusExpired, "updated_at": time.Now()}},
		)

		// Audit log for booking timeout
		auditLog := models.AuditLog{
			ID:        primitive.NewObjectID(),
			EventType: "BOOKING_TIMEOUT",
			BookingID: &bookingID,
			Payload: map[string]interface{}{
				"reason": "lock_expired",
			},
			CreatedAt: time.Now(),
		}
		w.Mongo.Collection("audit_logs").InsertOne(ctx, auditLog)
	}
}
