package repository

import (
	"errors"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func (r *productRepository) InjectDB(db *gorm.DB) error {
	if db == nil {
		return errors.New("invalid db")
	}
	r.db = db
	return nil
}

func (r *productRepository) InjectRedisClient(client *redis.Client) error {
	if client == nil {
		return errors.New("invalid redis client")
	}
	r.redisClient = client
	return nil
}
