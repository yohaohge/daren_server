package store

import (
	"LittleVideo/config"
	"time"

	"github.com/go-redis/redis"
)

type RedisClient struct {
	c *redis.Client
}

var RC *RedisClient

func ConnectRedis() {
	RC = new(RedisClient)
	RC.c = newLoginRedisClient(100)
}

func newLoginRedisClient(poolSize int) *redis.Client {
	c := redis.NewClient(&redis.Options{
		Addr:         config.GetLoginRedisAddr(),
		Password:     config.GetLoginRedisPassword(),
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     poolSize,
		PoolTimeout:  30 * time.Second,
	})
	_, err := c.Ping().Result()
	if err != nil {
		panic(err)
	}
	return c
}

func (rc *RedisClient) IsSessionValid() bool {
	return true
}

func (rc *RedisClient) CreateSession() string {
	return ""
}
