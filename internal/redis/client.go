package redisClient

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func New(host string) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return &RedisClient{rdb}
}
func (c *RedisClient) Increment(pageID string) error {
	if pageID == "" {
		return errors.New("pageID is empty")
	}
	c.client.Incr(context.Background(), pageID)
	return nil
}

func (c *RedisClient) GetCount(pageID string) (int, error) {
	if pageID == "" {
		return 0, errors.New("pageID is empty")
	}
	count, err := c.client.Get(context.Background(), pageID).Int()
	if err == redis.Nil {
		return 0, errors.New("pageID not found")
	}
	if err != nil {
		return 0, err
	}
	return count, nil
}
