package infrastructure

import (
	goredis "github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"
	"github.com/krobus00/product-service/internal/config"
	"github.com/sirupsen/logrus"
)

func NewAsynqClient() (*asynq.Client, error) {
	redisURL, err := goredis.ParseURL(config.RedisAsynqHost())
	if err != nil {
		return nil, err
	}
	client := asynq.NewClient(asynq.RedisClientOpt{
		Network:      redisURL.Network,
		Addr:         redisURL.Addr,
		DB:           redisURL.DB,
		Username:     redisURL.Username,
		Password:     redisURL.Password,
		DialTimeout:  config.RedisDialTimeout(),
		WriteTimeout: config.RedisWriteTimeout(),
		ReadTimeout:  config.RedisReadTimeout(),
	})

	return client, nil
}

func NewAsynqServer() (*asynq.Server, error) {
	redisURL, err := goredis.ParseURL(config.RedisAsynqHost())
	if err != nil {
		return nil, err
	}
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Network:      redisURL.Network,
			Addr:         redisURL.Addr,
			DB:           redisURL.DB,
			Username:     redisURL.Username,
			Password:     redisURL.Password,
			DialTimeout:  config.RedisDialTimeout(),
			WriteTimeout: config.RedisWriteTimeout(),
			ReadTimeout:  config.RedisReadTimeout(),
		},
		asynq.Config{
			Concurrency: config.AsynqConcurrency(),
			Logger:      logrus.New(),
		},
	)
	return srv, nil
}
