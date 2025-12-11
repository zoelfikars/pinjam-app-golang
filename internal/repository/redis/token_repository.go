package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type TokenRepository interface {
	BlacklistToken(ctx context.Context, token string, expiration time.Duration) error
	IsTokenBlacklisted(ctx context.Context, token string) (bool, error)
}
type tokenRepository struct {
	Client *redis.Client
}

func NewTokenRepository(client *redis.Client) TokenRepository {
	return &tokenRepository{Client: client}
}
func (r *tokenRepository) BlacklistToken(ctx context.Context, token string, expiration time.Duration) error {
	return r.Client.Set(ctx, getTokenKey(token), "blacklisted", expiration).Err()
}
func (r *tokenRepository) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	val, err := r.Client.Get(ctx, getTokenKey(token)).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return val == "blacklisted", nil
}
func getTokenKey(token string) string {
	return fmt.Sprintf("token:blacklist:%s", token)
}
