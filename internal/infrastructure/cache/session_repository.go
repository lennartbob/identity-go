package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vondr/identity-go/internal/core"
)

type SessionData struct {
	MemberID       string `json:"member_id"`
	Email          string `json:"email"`
	OrganizationID string `json:"organization_id"`
	MicrosoftID    string `json:"microsoft_id"`
}

type SessionRepository interface {
	CreateSession(ctx context.Context, token string, sessionData SessionData, ttl time.Duration) error
	GetSession(ctx context.Context, token string) (*SessionData, error)
	DeleteSession(ctx context.Context, token string) error
}

type RedisSessionRepository struct {
	redisClient *redis.Client
}

func NewRedisSessionRepository(redisClient *redis.Client) *RedisSessionRepository {
	return &RedisSessionRepository{redisClient: redisClient}
}

func (r *RedisSessionRepository) CreateSession(ctx context.Context, token string, sessionData SessionData, ttl time.Duration) error {
	data, err := json.Marshal(sessionData)
	if err != nil {
		return err
	}

	return r.redisClient.Set(ctx, "session:"+token, data, ttl).Err()
}

func (r *RedisSessionRepository) GetSession(ctx context.Context, token string) (*SessionData, error) {
	data, err := r.redisClient.Get(ctx, "session:"+token).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, core.ErrInvalidSession
		}
		return nil, err
	}

	var sessionData SessionData
	if err := json.Unmarshal([]byte(data), &sessionData); err != nil {
		return nil, err
	}

	return &sessionData, nil
}

func (r *RedisSessionRepository) DeleteSession(ctx context.Context, token string) error {
	return r.redisClient.Del(ctx, "session:"+token).Err()
}
