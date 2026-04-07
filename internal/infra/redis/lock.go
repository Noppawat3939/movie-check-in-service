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

// AcquireLock — atomic SET NX via Lua script
// return true = lock (can lock), false = key exits (cannot lock)
func (r *LockRepository) AcquireLock(ctx context.Context, key string, value string, ttl time.Duration) (bool, error) {
	// // check key if not exits can allow to set then lock same key
	// script := `
	// 	if redis.call("EXISTS", KEYS[1]) == 0 then
	// 		redis.call("SET", KEYS[1], ARGV[1])
	// 		redis.call("PEXPIRE", KEYS[1], ARGV[2])
	// 		return 1
	// 	end
	// 	return 0
	// `

	// ttlMs := ttl.Milliseconds()

	result, err := r.client.Conn.SetArgs(ctx, key, value, redis.SetArgs{Mode: "NX", TTL: ttl}).Result()

	// result, err := r.client.Conn.Eval(ctx, script, []string{key}, 1, ttlMs).Int()
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
