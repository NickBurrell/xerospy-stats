package app

import (
	"github.com/go-redis/redis"
)

var cache *redis.Client

func InitRedis(address string, password string, db int) {

	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})

	cache = client
}

func GetCache() *redis.Client {
	return cache
}
