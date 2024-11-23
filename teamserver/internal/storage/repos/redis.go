package repos

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/ksel172/Meduza/teamserver/conf"
	"log"
	"strings"
)

type Service interface {
	StringGet(ctx context.Context, key string) (string, error)
	StringSet(ctx context.Context, key, value string) error
	JsonSet(ctx context.Context, key string, value interface{}) error
	JsonGet(ctx context.Context, key string) (string, error)
	JsonDelete(ctx context.Context, key string) error
	GetAllByPartial(ctx context.Context, partialKey string) ([]interface{}, error)
	DeleteAllByPartial(ctx context.Context, partialKey string) error
}

type redisService struct {
	client *redis.Client
}

func NewRedisService() Service {
	client := redis.NewClient(&redis.Options{
		Addr:     conf.GetMeduzaRedisAddress(),
		Password: conf.GetMeduzaRedisPassword(),
		DB:       0, // TODO Default DB
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		log.Fatalf("redis connect err: %v", err)
	}

	return &redisService{client: client}
}

func (r *redisService) GetClient() *redis.Client {
	return r.client
}

func (r *redisService) StringGet(ctx context.Context, key string) (string, error) {
	if strings.TrimSpace(key) == "" {
		return "", fmt.Errorf("key cannot be empty")
	}

	result, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", fmt.Errorf("failed StringGet: %w", err)
	}

	return result, nil
}

func (r *redisService) StringSet(ctx context.Context, key, value string) error {
	if strings.TrimSpace(key) == "" || strings.TrimSpace(value) == "" {
		return fmt.Errorf("key and value cannot be empty")
	}

	err := r.client.Set(ctx, key, value, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to set key and value: %w", err)
	}

	return nil
}

func (r *redisService) JsonSet(ctx context.Context, key string, value interface{}) error {
	if strings.TrimSpace(key) == "" || value == nil {
		return fmt.Errorf("key and value cannot be empty")
	}

	serialized, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	err = r.client.Do(ctx, "JSON.SET", key, ".", serialized).Err()
	if err != nil {
		return fmt.Errorf("failed to set key and value: %w", err)
	}

	return nil
}

func (r *redisService) JsonGet(ctx context.Context, key string) (string, error) {
	if strings.TrimSpace(key) == "" {
		return "", fmt.Errorf("key cannot be empty")
	}

	result, err := r.client.Do(ctx, "JSON.GET", key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", fmt.Errorf("failed JSONGet: %w", err)
	}

	return result.(string), nil
}

func (r *redisService) JsonDelete(ctx context.Context, key string) error {
	if strings.TrimSpace(key) == "" {
		return fmt.Errorf("key cannot be empty")
	}

	deleted, err := r.client.Do(ctx, "JSON.DEL", key).Int()
	if err != nil {
		return fmt.Errorf("failed to delete key: %w", err)
	}
	if deleted == 0 {
		return fmt.Errorf("key not found")
	}

	return nil
}

func (r *redisService) GetAllByPartial(ctx context.Context, partialKey string) ([]interface{}, error) {
	if strings.TrimSpace(partialKey) == "" {
		return nil, fmt.Errorf("partialKey cannot be empty")
	}

	var results []interface{}

	iter := r.client.Scan(ctx, 0, partialKey+"*", 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		value, err := r.JsonGet(ctx, key)
		if err != nil {
			log.Printf("Failed to get value by partial at key: %s, error: %v", key, err)
			continue
		}
		results = append(results, value)
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *redisService) DeleteAllByPartial(ctx context.Context, partialKey string) error {
	if strings.TrimSpace(partialKey) == "" {
		return fmt.Errorf("partialKey cannot be empty")
	}

	iter := r.client.Scan(ctx, 0, partialKey+"*", 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		if err := r.JsonDelete(ctx, key); err != nil {
			log.Printf("Failed to delete by partial at key: %s, error: %v", key, err)
		}
	}

	if err := iter.Err(); err != nil {
		return err
	}

	return nil
}
