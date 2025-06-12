package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(addr, password string, db int) *RedisClient {
	if addr == "" {
		addr = "localhost:6379" // default Redis address
	}

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Test koneksi
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}

	return &RedisClient{
		client: client,
	}
}

func (r *RedisClient) SetToken(ctx context.Context, userID uint, token string, expiration time.Duration) error {
	if r.client == nil {
		return redis.ErrClosed
	}
	key := getTokenKey(userID)
	return r.client.Set(ctx, key, token, expiration).Err()
}

func (r *RedisClient) GetToken(ctx context.Context, userID uint) (string, error) {
	if r.client == nil {
		return "", redis.ErrClosed
	}
	key := getTokenKey(userID)
	return r.client.Get(ctx, key).Result()
}

func (r *RedisClient) DeleteToken(ctx context.Context, userID uint) error {
	if r.client == nil {
		return redis.ErrClosed
	}
	key := getTokenKey(userID)
	return r.client.Del(ctx, key).Err()
}

func getTokenKey(userID uint) string {
	return "user_token:" + string(userID)
}
