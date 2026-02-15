package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"cinema-booking/internal/models"
	"cinema-booking/internal/mq"
	"cinema-booking/internal/services"
	wsHub "cinema-booking/internal/ws"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookingHandler struct {
	Mongo   *services.MongoService
	Redis   *services.RedisService
	Hub     *wsHub.Hub
	MQ      *mq.MQService
	LockTTL time.Duration
}

type LockRequest struct {
	Seats []string `json:"seats" binding:"required"`
}

func (h *BookingHandler) LockSeats(c *gin.Context) {
	showtimeIdStr := c.Param("id")
	showtimeID, err := primitive.ObjectIDFromHex(showtimeIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid showtime id"})
		return
	}

	userIdStr, _ := c.Get("user_id")
	userId := userIdStr.(string)
	userOID, _ := primitive.ObjectIDFromHex(userId)

	var req LockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.Seats) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no seats selected"})
		return
	}

	ctx := context.Background()

	// Try to acquire Redis locks for all seats
	var lockedSeats []string
	for _, seat := range req.Seats {
		ok, err := h.Redis.AcquireLock(ctx, showtimeIdStr, seat, userId, h.LockTTL)
		if err != nil || !ok {
			// Release all acquired locks
			for _, ls := range lockedSeats {
				h.Redis.ReleaseLock(ctx, showtimeIdStr, ls, userId)
			}
			c.JSON(http.StatusConflict, gin.H{"error": "seat " + seat + " is already locked"})
			return
		}
		lockedSeats = append(lockedSeats, seat)
	}

	lockExpiresAt := time.Now().Add(h.LockTTL)

	// Update seat_reservations in MongoDB
	for _, seat := range req.Seats {
		filter := bson.M{
			"showtime_id": showtimeID,
			"seat_code":   seat,
			"state":       models.SeatStateAvailable,
		}
		update := bson.M{
			"$set": bson.M{
				"state":             models.SeatStateLocked,
				"locked_by_user_id": userOID,
				"lock_expires_at":   lockExpiresAt,
				"updated_at":        time.Now(),
			},
		}
		result, err := h.Mongo.Collection("seat_reservations").UpdateOne(ctx, filter, update)
		if err != nil || result.MatchedCount == 0 {
			// Rollback: release Redis locks and revert any MongoDB changes
			for _, ls := range req.Seats {
				h.Redis.ReleaseLock(ctx, showtimeIdStr, ls, userId)
				h.Mongo.Collection("seat_reservations").UpdateOne(ctx,
					bson.M{"showtime_id": showtimeID, "seat_code": ls, "locked_by_user_id": userOID},
					bson.M{"$set": bson.M{"state": models.SeatStateAvailable, "locked_by_user_id": nil, "lock_expires_at": nil}},
				)
			}
			c.JSON(http.StatusConflict, gin.H{"error": "seat " + seat + " is not available"})
			return
		}
	}

	// Create booking
	bookingID := primitive.NewObjectID()
	booking := models.Booking{
		ID:            bookingID,
		UserID:        userOID,
		ShowtimeID:    showtimeID,
		Seats:         req.Seats,
		Status:        models.BookingStatusLocked,
		LockExpiresAt: &lockExpiresAt,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	_, err = h.Mongo.Collection("bookings").InsertOne(ctx, booking)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create booking"})
		return
	}

	// Update seat_reservations with booking_id
	for _, seat := range req.Seats {
		h.Mongo.Collection("seat_reservations").UpdateOne(ctx,
			bson.M{"showtime_id": showtimeID, "seat_code": seat, "locked_by_user_id": userOID},
			bson.M{"$set": bson.M{"booking_id": bookingID}},
		)
	}

	// Broadcast SEAT_LOCKED for each seat
	for _, seat := range req.Seats {
		msg, _ := json.Marshal(map[string]interface{}{
			"type":           "SEAT_LOCKED",
			"seatCode":       seat,
			"lockedByUserId": userId,
			"lockExpiresAt":  lockExpiresAt.Format(time.RFC3339),
		})
		h.Hub.BroadcastToRoom(showtimeIdStr, msg)
	}

	c.JSON(http.StatusOK, gin.H{
		"bookingId":     bookingID.Hex(),
		"lockExpiresAt": lockExpiresAt.Format(time.RFC3339),
	})
}

func (h *BookingHandler) MockPayment(c *gin.Context) {
	bookingIdStr := c.Param("id")
	bookingID, err := primitive.ObjectIDFromHex(bookingIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking id"})
		return
	}

	ctx := context.Background()

	// Verify booking exists and belongs to user
	userIdStr, _ := c.Get("user_id")
	userOID, _ := primitive.ObjectIDFromHex(userIdStr.(string))

	var booking models.Booking
	err = h.Mongo.Collection("bookings").FindOne(ctx, bson.M{"_id": bookingID, "user_id": userOID}).Decode(&booking)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
		return
	}

	// Create mock payment (always succeeds)
	paymentID := primitive.NewObjectID()
	payment := models.Payment{
		ID:        paymentID,
		BookingID: bookingID,
		Amount:    float64(len(booking.Seats)) * 250.0, // 250 per seat
		Status:    models.PaymentStatusSuccess,
		Provider:  "MOCK",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = h.Mongo.Collection("payments").InsertOne(ctx, payment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create payment"})
		return
	}

	// Update booking with payment_id
	h.Mongo.Collection("bookings").UpdateOne(ctx,
		bson.M{"_id": bookingID},
		bson.M{"$set": bson.M{"payment_id": paymentID, "updated_at": time.Now()}},
	)

	c.JSON(http.StatusOK, gin.H{
		"paymentId": paymentID.Hex(),
		"status":    models.PaymentStatusSuccess,
		"amount":    payment.Amount,
	})
}

func (h *BookingHandler) ConfirmBooking(c *gin.Context) {
	bookingIdStr := c.Param("id")
	bookingID, err := primitive.ObjectIDFromHex(bookingIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking id"})
		return
	}

	ctx := context.Background()
	userIdStr, _ := c.Get("user_id")
	userId := userIdStr.(string)
	userOID, _ := primitive.ObjectIDFromHex(userId)

	// Get booking
	var booking models.Booking
	err = h.Mongo.Collection("bookings").FindOne(ctx, bson.M{
		"_id":     bookingID,
		"user_id": userOID,
		"status":  models.BookingStatusLocked,
	}).Decode(&booking)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "booking not found or not in LOCKED state"})
		return
	}

	// Check lock not expired
	if booking.LockExpiresAt != nil && booking.LockExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusConflict, gin.H{"error": "lock has expired"})
		return
	}

	// Check payment exists and is SUCCESS
	if booking.PaymentID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payment required"})
		return
	}

	var payment models.Payment
	err = h.Mongo.Collection("payments").FindOne(ctx, bson.M{
		"_id":    *booking.PaymentID,
		"status": models.PaymentStatusSuccess,
	}).Decode(&payment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payment not found or not successful"})
		return
	}

	// Atomic update: seat_reservations -> BOOKED
	showtimeIdStr := booking.ShowtimeID.Hex()
	for _, seat := range booking.Seats {
		result, err := h.Mongo.Collection("seat_reservations").UpdateOne(ctx,
			bson.M{
				"showtime_id":       booking.ShowtimeID,
				"seat_code":         seat,
				"state":             models.SeatStateLocked,
				"locked_by_user_id": userOID,
			},
			bson.M{"$set": bson.M{
				"state":      models.SeatStateBooked,
				"updated_at": time.Now(),
			}},
		)
		if err != nil || result.MatchedCount == 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "seat " + seat + " confirmation failed - may have been released"})
			return
		}
	}

	// Update booking -> BOOKED
	h.Mongo.Collection("bookings").UpdateOne(ctx,
		bson.M{"_id": bookingID},
		bson.M{"$set": bson.M{"status": models.BookingStatusBooked, "updated_at": time.Now()}},
	)

	// Release Redis locks
	for _, seat := range booking.Seats {
		h.Redis.ReleaseLock(ctx, showtimeIdStr, seat, userId)
	}

	// Broadcast SEAT_BOOKED
	for _, seat := range booking.Seats {
		msg, _ := json.Marshal(map[string]interface{}{
			"type":     "SEAT_BOOKED",
			"seatCode": seat,
		})
		h.Hub.BroadcastToRoom(showtimeIdStr, msg)
	}

	// ✅ ดึงข้อมูล User
	var user models.User
	err = h.Mongo.Collection("users").FindOne(ctx, bson.M{
		"_id": userOID,
	}).Decode(&user)

	if err != nil {
		log.Printf("Warning: Failed to get user: %v", err)
		// Fallback: ส่ง event โดยไม่มี email/name
		user.Email = ""
		user.Name = "Unknown User"
	}

	// ตอนนี้ใช้ได้แล้ว
	log.Printf("User: %s (%s)", user.Name, user.Email)

	// Publish MQ event
	event := mq.NewBookingConfirmedEvent(bookingIdStr, userId, showtimeIdStr, booking.Seats)
	h.MQ.SafePublish(event)

	c.JSON(http.StatusOK, gin.H{"status": "BOOKED"})
}

func (h *BookingHandler) CancelBooking(c *gin.Context) {
	bookingIdStr := c.Param("id")
	bookingID, err := primitive.ObjectIDFromHex(bookingIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking id"})
		return
	}

	ctx := context.Background()
	userIdStr, _ := c.Get("user_id")
	userId := userIdStr.(string)
	userOID, _ := primitive.ObjectIDFromHex(userId)

	var booking models.Booking
	err = h.Mongo.Collection("bookings").FindOne(ctx, bson.M{
		"_id":     bookingID,
		"user_id": userOID,
		"status":  models.BookingStatusLocked,
	}).Decode(&booking)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "booking not found or not in LOCKED state"})
		return
	}

	showtimeIdStr := booking.ShowtimeID.Hex()

	// Release Redis locks + update MongoDB
	for _, seat := range booking.Seats {
		h.Redis.ReleaseLock(ctx, showtimeIdStr, seat, userId)
		h.Mongo.Collection("seat_reservations").UpdateOne(ctx,
			bson.M{"showtime_id": booking.ShowtimeID, "seat_code": seat},
			bson.M{"$set": bson.M{
				"state":             models.SeatStateAvailable,
				"locked_by_user_id": nil,
				"lock_expires_at":   nil,
				"booking_id":        nil,
				"updated_at":        time.Now(),
			}},
		)
	}

	// Update booking -> CANCELLED
	h.Mongo.Collection("bookings").UpdateOne(ctx,
		bson.M{"_id": bookingID},
		bson.M{"$set": bson.M{"status": models.BookingStatusCancelled, "updated_at": time.Now()}},
	)

	// Broadcast SEAT_RELEASED
	for _, seat := range booking.Seats {
		msg, _ := json.Marshal(map[string]interface{}{
			"type":     "SEAT_RELEASED",
			"seatCode": seat,
		})
		h.Hub.BroadcastToRoom(showtimeIdStr, msg)
	}

	c.JSON(http.StatusOK, gin.H{"status": "CANCELLED"})
}

func (h *BookingHandler) GetBooking(c *gin.Context) {
	bookingIdStr := c.Param("id")
	bookingID, err := primitive.ObjectIDFromHex(bookingIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking id"})
		return
	}

	var booking models.Booking
	err = h.Mongo.Collection("bookings").FindOne(context.Background(), bson.M{"_id": bookingID}).Decode(&booking)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
		return
	}

	c.JSON(http.StatusOK, booking)
}
