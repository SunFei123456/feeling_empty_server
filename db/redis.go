package db

import (
	"context"
	"github.com/go-redis/redis/v8"
	"os"
)

var RedisClient *redis.Client

// InitRedis 初始化Redis连接
func InitRedis() error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),     // Redis地址
		Password: os.Getenv("REDIS_PASSWORD"), // Redis密码
		DB:       0,                          // 默认DB 0
	})

	// 测试连接
	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return err
	}

	return nil
}