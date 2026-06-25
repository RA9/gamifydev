// Package rdb wraps the Redis client. Redis is optional today (no feature
// depends on it yet) but wired in for upcoming work — caching, rate limiting,
// session/leaderboard storage, and pub/sub for live competitions.
package rdb

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Open parses a Redis connection URL (e.g. redis://:pass@host:6379/0, or the
// rediss:// TLS form Railway provides), connects, and verifies with a PING.
func Open(ctx context.Context, url string) (*redis.Client, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("parse REDIS_URL: %w", err)
	}
	client := redis.NewClient(opt)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("ping redis: %w", err)
	}
	return client, nil
}
