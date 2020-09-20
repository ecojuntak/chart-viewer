package repository

import (
	"github.com/go-redis/redis"
	"os"
)

type Repository interface {
	Set(string, string)
	Get(string) string
}

type repository struct {
	redisClient *redis.Client
}

func NewRepository() Repository {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")

	return repository{
		redisClient: redis.NewClient(&redis.Options{
			Addr: host + ":" + port,
		}),
	}
}

func (r repository) Set(key string, value string) {
	_ = r.redisClient.Set(key, value, 0)
}

func (r repository) Get(key string) string {
	value, _ := r.redisClient.Get(key).Result()
	return value
}