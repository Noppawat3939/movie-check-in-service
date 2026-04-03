package redis

import (
	"context"
	"log"
	"os"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type Client struct {
	Conn *goredis.Client
}

func NewClient() (*Client, error) {
	host := os.Getenv("REIDS_HOST")
	port := os.Getenv("REDIS_EXTERNAL_PORT")

	client := goredis.NewClient(&goredis.Options{
		Addr: host + ":" + port,
		DB:   0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Minute)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	log.Println("connected to redis")

	return &Client{Conn: client}, nil
}

func (c *Client) Close() error {
	return c.Conn.Close()
}
