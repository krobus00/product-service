package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/goccy/go-json"
	"github.com/krobus00/product-service/internal/config"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func WithPagination(page int, limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * limit
		return db.Limit(limit).Offset(offset)
	}
}

func WithSearch(value string, columns []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if value == "" {
			return db
		}
		for _, column := range columns {
			db.Where(fmt.Sprintf("%s LIKE ?", column), fmt.Sprintf("%%%s%%", value))
		}
		return db
	}
}

func HSetWithExpiry(ctx context.Context, redisClient *redis.Client, bucketCacheKey string, field string, data any) error {
	if config.DisableCaching() {
		return nil
	}
	cacheData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = redisClient.HSet(ctx, bucketCacheKey, field, cacheData).Err()
	if err != nil {
		return err
	}
	err = redisClient.ExpireNX(ctx, bucketCacheKey, config.RedisCacheTTL()).Err()
	if err != nil {
		return err
	}
	return nil
}

func SetWithExpiry(ctx context.Context, redisClient *redis.Client, cacheKey string, data any) error {
	if config.DisableCaching() {
		return nil
	}
	cacheData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = redisClient.Set(ctx, cacheKey, cacheData, config.RedisCacheTTL()).Err()
	if err != nil {
		return err
	}
	return nil
}

func Get(ctx context.Context, redisClient *redis.Client, cacheKey string) ([]byte, error) {
	if config.DisableCaching() {
		return nil, nil
	}
	cachedData, err := redisClient.Get(ctx, cacheKey).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		logrus.WithField("cacheKey", cacheKey).Error(err.Error())
		return nil, err
	}
	return cachedData, nil
}

func DeleteByKeys(ctx context.Context, redisClient *redis.Client, cacheKeys []string) error {
	if config.DisableCaching() {
		return nil
	}
	for _, cacheKey := range cacheKeys {
		err := redisClient.Del(ctx, cacheKey).Err()
		if err != nil && !errors.Is(err, redis.Nil) {
			logrus.WithField("cacheKey", cacheKey).Error(err.Error())
			return err
		}
	}
	return nil
}

func HGet(ctx context.Context, redisClient *redis.Client, bucketCacheKey string, field string) ([]byte, error) {
	if config.DisableCaching() {
		return nil, nil
	}
	cachedData, err := redisClient.HGet(ctx, bucketCacheKey, field).Bytes()
	if err != nil {
		return nil, err
	}
	return cachedData, nil
}
