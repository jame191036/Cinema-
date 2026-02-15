package handlers

import (
	"context"
	"fmt"
	"log"
	"time"

	"cinema-booking/internal/models"
	"cinema-booking/internal/services"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SeedData(mongo *services.MongoService) {
	ctx := context.Background()

	// Check if already seeded
	count, _ := mongo.Collection("movies").CountDocuments(ctx, bson.M{})
	if count > 0 {
		log.Println("Database already seeded, skipping")
		return
	}

	log.Println("Seeding database...")

	// Create seatmap
	seatmap := models.Seatmap{
		ID: "AUDI-01-V1",
	}

	rowLabels := []string{"A", "B", "C", "D", "E"}
	for _, label := range rowLabels {
		row := models.Row{RowLabel: label}
		for i := 1; i <= 8; i++ {
			seatType := "NORMAL"
			if label == "E" {
				seatType = "VIP"
			}
			row.Seats = append(row.Seats, models.Seat{
				SeatCode: fmt.Sprintf("%s%d", label, i),
				Type:     seatType,
				Active:   true,
			})
		}
		seatmap.Rows = append(seatmap.Rows, row)
	}

	mongo.Collection("seatmaps").InsertOne(ctx, seatmap)

	// Create movies
	movies := []models.Movie{
		{ID: primitive.NewObjectID(), Title: "Inception", DurationMin: 148, Rating: "PG-13", CreatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Title: "The Dark Knight", DurationMin: 152, Rating: "PG-13", CreatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Title: "Interstellar", DurationMin: 169, Rating: "PG-13", CreatedAt: time.Now()},
	}

	for _, m := range movies {
		mongo.Collection("movies").InsertOne(ctx, m)
	}

	// Create showtimes (today + tomorrow, 2 per movie)
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)

	var showtimes []models.Showtime
	times := []time.Time{
		today.Add(14 * time.Hour),    // 2 PM today
		today.Add(18 * time.Hour),    // 6 PM today
		tomorrow.Add(14 * time.Hour), // 2 PM tomorrow
		tomorrow.Add(18 * time.Hour), // 6 PM tomorrow
	}

	for _, movie := range movies {
		for i, t := range times {
			if i >= 2 {
				break // 2 showtimes per movie
			}
			st := models.Showtime{
				ID:           primitive.NewObjectID(),
				MovieID:      movie.ID,
				StartTime:    t,
				AuditoriumID: "Auditorium 1",
				SeatmapID:    "AUDI-01-V1",
				CreatedAt:    time.Now(),
			}
			showtimes = append(showtimes, st)
			mongo.Collection("showtimes").InsertOne(ctx, st)
		}
	}

	// Create seat_reservations for each showtime
	for _, st := range showtimes {
		for _, row := range seatmap.Rows {
			for _, seat := range row.Seats {
				sr := models.SeatReservation{
					ID:         primitive.NewObjectID(),
					ShowtimeID: st.ID,
					SeatCode:   seat.SeatCode,
					State:      models.SeatStateAvailable,
					UpdatedAt:  time.Now(),
				}
				mongo.Collection("seat_reservations").InsertOne(ctx, sr)
			}
		}
	}

	log.Printf("Seeded: %d movies, %d showtimes, %d seats per showtime", len(movies), len(showtimes), 40)
}
