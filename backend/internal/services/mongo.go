package services

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoService struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func NewMongoService(uri, dbName string) *MongoService {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	db := client.Database(dbName)
	svc := &MongoService{Client: client, DB: db}
	svc.createIndexes()
	return svc
}

func (s *MongoService) createIndexes() {
	ctx := context.Background()

	// users indexes
	s.DB.Collection("users").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "provider_sub", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "role", Value: 1}}},
	})

	// seat_reservations indexes
	s.DB.Collection("seat_reservations").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "showtime_id", Value: 1}, {Key: "seat_code", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "showtime_id", Value: 1}, {Key: "state", Value: 1}}},
	})

	// bookings indexes
	s.DB.Collection("bookings").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "created_at", Value: 1}}},
		{Keys: bson.D{{Key: "showtime_id", Value: 1}, {Key: "created_at", Value: 1}}},
	})

	// audit_logs indexes
	s.DB.Collection("audit_logs").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "event_type", Value: 1}, {Key: "created_at", Value: 1}}},
		{Keys: bson.D{{Key: "showtime_id", Value: 1}, {Key: "created_at", Value: 1}}},
	})

	log.Println("MongoDB indexes created")
}

func (s *MongoService) Collection(name string) *mongo.Collection {
	return s.DB.Collection(name)
}
