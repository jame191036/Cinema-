package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"cinema-booking/internal/config"
	"cinema-booking/internal/handlers"
	"cinema-booking/internal/middleware"
	"cinema-booking/internal/models"
	"cinema-booking/internal/mq"
	"cinema-booking/internal/services"
	"cinema-booking/internal/worker"
	wsHub "cinema-booking/internal/ws"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	cfg := config.Load()

	// Connect services
	mongoSvc := services.NewMongoService(cfg.MongoURI, cfg.MongoDB)
	redisSvc := services.NewRedisService(cfg.RedisAddr, cfg.RedisPassword)

	// Try MQ connection (optional in dev)
	var mqSvc *mq.MQService
	mqSvc = mq.NewMQService(cfg.RabbitMQURL)

	// Seed data
	handlers.SeedData(mongoSvc)

	// WebSocket hub
	hub := wsHub.NewHub()
	go hub.Run()

	// Start timeout worker
	lockTTL := time.Duration(cfg.SeatLockTTL) * time.Second
	workerInterval := time.Duration(cfg.WorkerInterval) * time.Second
	tw := worker.NewTimeoutWorker(mongoSvc, redisSvc, hub, workerInterval)
	tw.Start()
	emailSvc := services.NewEmailService(
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.SMTPUsername,
		cfg.SMTPPassword,
		cfg.SMTPFrom,
	)

	// MQ Consumer - writes audit logs on BookingConfirmed
	if mqSvc.IsConnected() {
		mqSvc.Consume(func(event mq.BookingEvent) {

			if event.EventType == "BookingConfirmed" {
				bookingOID, _ := primitive.ObjectIDFromHex(event.BookingID)
				userOID, _ := primitive.ObjectIDFromHex(event.UserID)
				showtimeOID, _ := primitive.ObjectIDFromHex(event.ShowtimeID)

				auditLog := models.AuditLog{
					ID:         primitive.NewObjectID(),
					EventType:  "BOOKING_SUCCESS",
					UserID:     &userOID,
					ShowtimeID: &showtimeOID,
					BookingID:  &bookingOID,
					Payload: map[string]interface{}{
						"seats":    event.Seats,
						"event_id": event.EventID,
					},
					CreatedAt: time.Now(),
				}
				mongoSvc.Collection("audit_logs").InsertOne(context.Background(), auditLog)

				if event.UserEmail != "" {
					go func(e mq.BookingEvent) {
						err := emailSvc.SendBookingConfirmation(services.BookingConfirmationData{
							UserName:   e.UserName,
							UserEmail:  e.UserEmail,
							BookingID:  e.BookingID,
							Seats:      e.Seats,
							ShowtimeID: e.ShowtimeID,
							OccurredAt: e.OccurredAt,
						})
						if err != nil {
							log.Printf("Email send error: %v", err)
						}
					}(event)
				}
			}
		})
	}

	// Setup Gin router
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Firebase Auth
	firebaseAuth := services.NewFirebaseAuth(cfg.FirebaseProjectID)
	log.Printf("Firebase Auth initialized for project: %s", cfg.FirebaseProjectID)

	// Handlers
	authHandler := &handlers.AuthHandler{Mongo: mongoSvc, JWTSecret: cfg.JWTSecret, Firebase: firebaseAuth}
	movieHandler := &handlers.MovieHandler{Mongo: mongoSvc}
	showtimeHandler := &handlers.ShowtimeHandler{Mongo: mongoSvc}
	bookingHandler := &handlers.BookingHandler{
		Mongo:   mongoSvc,
		Redis:   redisSvc,
		Hub:     hub,
		MQ:      mqSvc,
		LockTTL: lockTTL,
	}
	adminHandler := &handlers.AdminHandler{Mongo: mongoSvc}

	// Public routes
	r.POST("/api/auth/login", authHandler.Login)
	r.POST("/api/auth/google", authHandler.GoogleLogin)

	// Auth protected routes
	auth := r.Group("/api")
	auth.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		auth.GET("/movies", movieHandler.ListMovies)
		auth.GET("/movies/:id", movieHandler.GetMovie)
		auth.GET("/showtimes", showtimeHandler.ListShowtimes)
		auth.GET("/showtimes/:id/seats", showtimeHandler.GetSeats)

		auth.POST("/showtimes/:id/seats/lock", bookingHandler.LockSeats)
		auth.POST("/bookings/:id/pay", bookingHandler.MockPayment)
		auth.POST("/bookings/:id/confirm", bookingHandler.ConfirmBooking)
		auth.POST("/bookings/:id/cancel", bookingHandler.CancelBooking)
		auth.GET("/bookings/:id", bookingHandler.GetBooking)
	}

	// Admin routes
	admin := r.Group("/api/admin")
	admin.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	admin.Use(middleware.AdminMiddleware())
	{
		admin.GET("/bookings", adminHandler.ListBookings)
		admin.GET("/audit-logs", adminHandler.ListAuditLogs)
	}

	// WebSocket route
	r.GET("/ws/showtimes/:showtimeId", func(c *gin.Context) {
		showtimeId := c.Param("showtimeId")

		// Upgrade connection
		conn, err := wsHub.Upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		client := &wsHub.Client{
			Hub:        hub,
			Conn:       conn,
			ShowtimeID: showtimeId,
			Send:       make(chan []byte, 256),
		}

		hub.Register <- client
		go client.WritePump()

		// Send initial snapshot
		go func() {
			showtimeOID, err := primitive.ObjectIDFromHex(showtimeId)
			if err != nil {
				return
			}

			var seats []models.SeatReservation
			cursor, err := mongoSvc.Collection("seat_reservations").Find(
				context.Background(),
				bson.M{"showtime_id": showtimeOID},
			)
			if err != nil {
				return
			}
			defer cursor.Close(context.Background())
			cursor.All(context.Background(), &seats)

			snapshot, _ := json.Marshal(map[string]interface{}{
				"type":  "SYNC_SNAPSHOT",
				"seats": seats,
			})

			select {
			case client.Send <- snapshot:
			default:
			}
		}()

		client.ReadPump()
	})

	// Health check
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	port := cfg.BackendPort
	log.Printf("Server starting on port %s", port)
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
