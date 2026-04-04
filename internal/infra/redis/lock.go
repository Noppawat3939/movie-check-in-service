package redis

import (
	"context"
	"time"
)

type LockService struct {
	client *Client
}

func NewLockService(c *Client) *LockService {
	return &LockService{client: c}
}

// AcquireLock tries to acquire a distributed lock with TTL
func (s *LockService) AcquireLock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	return false, nil
}
