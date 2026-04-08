package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type LockRepository struct {
	client *Client
}

func NewLockRepository(client *Client) *LockRepository {
	return &LockRepository{client}
}

func (r *LockRepository) AcquireLock(ctx context.Context, key string, value string, ttl time.Duration) (bool, error) {
	result, err := r.client.Conn.SetArgs(ctx, key, value, redis.SetArgs{Mode: "NX", TTL: ttl}).Result()

	if err != nil {
		return false, fmt.Errorf("acquire lock failed %w", err)
	}

	return result == "OK", nil
}

// ReleaseLock - remove key out of redis
func (r *LockRepository) ReleaseLock(ctx context.Context, key string, value string) error {
	script := `
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("DEL", KEYS[1])
		end
		return 0
	`

	_, err := r.client.Conn.Eval(ctx, script, []string{key}, value).Result()

	if err != nil {
		return fmt.Errorf("release lock failed %w", err)
	}

	return nil
}
