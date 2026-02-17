package config

import (
	"os"
	"strconv"
)

type Config struct {
	MongoURI          string
	MongoDB           string
	RedisAddr         string
	RedisPassword     string
	RabbitMQURL       string
	JWTSecret         string
	GoogleClientID    string
	FirebaseProjectID string
	BackendPort       string
	SeatLockTTL       int
	WorkerInterval    int

	// Email (SMTP)
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	SMTPFrom     string
}

func Load() *Config {
	return &Config{
		MongoURI:          getEnv("MONGO_URI", "mongodb://mongo:27017/cinema"),
		MongoDB:           getEnv("MONGO_DB", "cinema"),
		RedisAddr:         getEnv("REDIS_ADDR", "redis:6379"),
		RedisPassword:     getEnv("REDIS_PASSWORD", ""),
		RabbitMQURL:       getEnv("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/"),
		JWTSecret:         getEnv("JWT_SECRET", "change-me-in-production"),
		GoogleClientID:    getEnv("GOOGLE_CLIENT_ID", ""),
		FirebaseProjectID: getEnv("FIREBASE_PROJECT_ID", "cinema-25e75"),
		BackendPort:       getEnv("BACKEND_PORT", "8080"),
		SeatLockTTL:       getEnvInt("SEAT_LOCK_TTL", 300),
		WorkerInterval:    getEnvInt("WORKER_INTERVAL", 5),

		// Email config
		SMTPHost:     getEnv("SMTP_HOST", ""),
		SMTPPort:     getEnv("SMTP_PORT", "587"),
		SMTPUsername: getEnv("SMTP_USERNAME", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:     getEnv("SMTP_FROM", "noreply@cinema.local"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}
