package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type Cache struct {
	client *Client
}

func NewCache(client *Client) *Cache {
	return &Cache{client}
}

func (c *Cache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("cache marshal failed %w", err)
	}

	if err := c.client.Conn.Set(ctx, key, b, ttl).Err(); err != nil {
		return fmt.Errorf("cache set failed %w", err)
	}

	return nil
}

func (c *Cache) Get(ctx context.Context, key string, dest any) error {
	val, err := c.client.Conn.Get(ctx, key).Result()
	if err == goredis.Nil {
		return fmt.Errorf("cache missing %w", err)
	}

	if err != nil {
		return fmt.Errorf("cache get failed %w", err)
	}

	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return fmt.Errorf("cache unmarshal failed %w", err)
	}

	return nil
}

func (c *Cache) Del(ctx context.Context, key string) error {
	if err := c.client.Conn.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("cache delete failed %w", err)
	}

	return nil
}
