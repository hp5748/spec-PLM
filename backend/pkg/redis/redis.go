package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"plm/internal/core/config"
)

var Client *redis.Client

// Init 初始化Redis连接
func Init(cfg *config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// 测试连接
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect redis: %w", err)
	}

	Client = client
	fmt.Println("Redis connected successfully")
	return client, nil
}

// Get 获取Redis客户端
func Get() *redis.Client {
	return Client
}
