package services

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	Client *redis.Client
}

func NewRedisService(addr, password string) *RedisService {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		// Log but don't fatal - redis might still be starting in docker
		fmt.Printf("Warning: Redis ping failed: %v\n", err)
	}

	return &RedisService{Client: client}
}

func lockKey(showtimeId, seatCode string) string {
	return fmt.Sprintf("lock:showtime:%s:seat:%s", showtimeId, seatCode)
}

func (s *RedisService) AcquireLock(ctx context.Context, showtimeId, seatCode, userId string, ttl time.Duration) (bool, error) {
	key := lockKey(showtimeId, seatCode)
	ok, err := s.Client.SetNX(ctx, key, userId, ttl).Result()
	return ok, err
}

func (s *RedisService) ReleaseLock(ctx context.Context, showtimeId, seatCode, userId string) error {
	key := lockKey(showtimeId, seatCode)
	// Lua script: only delete if value matches
	script := redis.NewScript(`
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("DEL", KEYS[1])
		end
		return 0
	`)
	_, err := script.Run(ctx, s.Client, []string{key}, userId).Result()
	return err
}

func (s *RedisService) GetLockOwner(ctx context.Context, showtimeId, seatCode string) (string, error) {
	key := lockKey(showtimeId, seatCode)
	val, err := s.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func (s *RedisService) ForceReleaseLock(ctx context.Context, showtimeId, seatCode string) error {
	key := lockKey(showtimeId, seatCode)
	return s.Client.Del(ctx, key).Err()
}
