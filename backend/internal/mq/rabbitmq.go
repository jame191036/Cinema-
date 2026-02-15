package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type BookingEvent struct {
	EventID    string   `json:"eventId"`
	EventType  string   `json:"eventType"`
	OccurredAt string   `json:"occurredAt"`
	BookingID  string   `json:"bookingId"`
	UserID     string   `json:"userId"`
	ShowtimeID string   `json:"showtimeId"`
	Seats      []string `json:"seats"`
}

type MQService struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

func NewMQService(url string) *MQService {
	var conn *amqp.Connection
	var err error

	// Retry connection (RabbitMQ might not be ready yet in Docker)
	for i := 0; i < 30; i++ {
		conn, err = amqp.Dial(url)
		if err == nil {
			break
		}
		log.Printf("RabbitMQ not ready, retrying in 2s... (%d/30)", i+1)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open channel: %v", err)
	}

	q, err := ch.QueueDeclare("booking.events", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	return &MQService{conn: conn, channel: ch, queue: q}
}

func (s *MQService) Publish(event BookingEvent) error {
	if event.EventID == "" {
		event.EventID = uuid.New().String()
	}
	if event.OccurredAt == "" {
		event.OccurredAt = time.Now().UTC().Format(time.RFC3339)
	}

	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.channel.PublishWithContext(ctx, "", s.queue.Name, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}

func (s *MQService) Consume(handler func(BookingEvent)) {
	msgs, err := s.channel.Consume(s.queue.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to consume: %v", err)
	}

	go func() {
		for msg := range msgs {
			var event BookingEvent
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				log.Printf("Failed to unmarshal event: %v", err)
				continue
			}
			handler(event)
		}
	}()

	log.Println("MQ Consumer started")
}

func (s *MQService) Close() {
	if s.channel != nil {
		s.channel.Close()
	}
	if s.conn != nil {
		s.conn.Close()
	}
}

func NewBookingConfirmedEvent(bookingID, userID, showtimeID string, seats []string) BookingEvent {
	return BookingEvent{
		EventID:    uuid.New().String(),
		EventType:  "BookingConfirmed",
		OccurredAt: time.Now().UTC().Format(time.RFC3339),
		BookingID:  bookingID,
		UserID:     userID,
		ShowtimeID: showtimeID,
		Seats:      seats,
	}

}

// IgnoreConnectionError is used to make MQ optional during development
func IgnoreConnectionError() *MQService {
	return &MQService{}
}

func (s *MQService) IsConnected() bool {
	return s.conn != nil && !s.conn.IsClosed()
}

// SafePublish publishes if connected, logs warning otherwise
func (s *MQService) SafePublish(event BookingEvent) {
	if !s.IsConnected() {
		fmt.Printf("MQ not connected, skipping event: %s\n", event.EventType)
		return
	}
	if err := s.Publish(event); err != nil {
		fmt.Printf("Failed to publish event: %v\n", err)
	}
}
