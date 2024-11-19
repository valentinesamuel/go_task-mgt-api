package services

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisService struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisService() *RedisService {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return &RedisService{
		client: client,
		ctx:    context.Background(),
	}
}

// Set stores data in Redis with optional expiration
func (rs *RedisService) Set(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return rs.client.Set(rs.ctx, key, data, expiration).Err()
}

// Get retrieves data from Redis and unmarshals into dest
func (rs *RedisService) Get(key string, dest interface{}) error {
	data, err := rs.client.Get(rs.ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), dest)
}

// Delete removes a key from Redis
func (rs *RedisService) Delete(key string) error {
	return rs.client.Del(rs.ctx, key).Err()
}

// Close closes the Redis connection
func (rs *RedisService) Close() error {
	return rs.client.Close()
}
