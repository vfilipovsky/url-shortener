package cache

//go:generate mockgen -source=./cache.go -destination=./mock/mock_cache.go

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/vfilipovsky/url-shortener/pkg/config"
	"github.com/vfilipovsky/url-shortener/pkg/logger"
)

type Repository interface {
	Get(key string) (string, error)
	GetStruct(key string, entity any) error
	Set(key, value string, exp time.Duration) error
	SetStruct(key string, value any, exp time.Duration) error
	Delete(key string) error
	Close() error
}

type redisService struct {
	instance *redis.Client
}

var (
	ctx = context.Background()
)

func NewInstance(config *config.RedisCredentials) (Repository, error) {
	r := &redisService{
		instance: redis.NewClient(
			&redis.Options{
				Addr:     fmt.Sprintf("%s:%s", config.Host, config.Port),
				Password: config.Pass,
				DB:       0,
			}),
	}

	_, err := r.instance.Ping(context.Background()).Result()

	if err != nil {
		return &redisService{}, err
	}

	logger.Infof("Cache connected")

	return r, nil
}

func (r *redisService) Get(key string) (string, error) {
	return r.instance.Get(ctx, key).Result()
}

func (r *redisService) Set(key, value string, exp time.Duration) error {
	return r.instance.Set(ctx, key, value, exp).Err()
}

func (r *redisService) Delete(key string) error {
	return r.instance.Del(ctx, key).Err()
}

func (r *redisService) Close() error {
	return r.instance.Close()
}

func (r *redisService) GetStruct(key string, entity any) error {
	val, err := r.instance.Get(ctx, key).Result()

	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), entity)
}

func (r *redisService) SetStruct(key string, val any, exp time.Duration) error {
	data, err := json.Marshal(val)

	if err != nil {
		return err
	}

	return r.instance.Set(ctx, key, data, exp).Err()
}
