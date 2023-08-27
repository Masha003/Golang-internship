package data

import (
	"context"
	"log"

	"github.com/Masha003/Golang-internship.git/internal/config"
	"github.com/redis/go-redis/v9"
)

func NewRDB(cfg config.Config) (*redis.Client, error) {
	log.Print("Initializing redis")

	opt, err := redis.ParseURL(cfg.RedisUrl)
	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(opt)
	redisCtx := context.Background()

	status := rdb.Ping(redisCtx)
	if status.Err() != nil {
		return nil, status.Err()
	}

	return rdb, nil
}
