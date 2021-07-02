package model

import "github.com/go-redis/redis/v7"

var RedisClient *redis.Client

func RedisInit() error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "xx:6379",
		Password: "xx",
	})

	_, err := RedisClient.Ping().Result()
	if err != nil {
		return err
	}

	return nil
}
