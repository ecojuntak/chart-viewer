package repository

import (
	"github.com/go-redis/redis"
)

type Repository interface {
	Set(string, string)
	Get(string) string
}

type repository struct {
	redisClient *redis.Client
}

func NewRepository(redisAddress string) (error, Repository) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddress,
	})

	status := redisClient.Ping()
	err := status.Err()
	if err != nil {
		return err, nil
	}

	return nil, repository{
		redisClient: redisClient,
	}
}

func (r repository) Set(key string, value string) {
	_ = r.redisClient.Set(key, value, 0)
}

func (r repository) Get(key string) string {
	value, _ := r.redisClient.Get(key).Result()
	return value
}
