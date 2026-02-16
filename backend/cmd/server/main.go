package main

import (
	"log"

	"plm/internal/core/config"
	"plm/internal/core/module"
	"plm/internal/core/server"
	"plm/internal/modules/user/repository"
	userModule "plm/internal/modules/user"
	"plm/pkg/database"
	"plm/pkg/redis"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库
	db, err := database.Init(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}

	// 设置Repository的数据库连接
	repository.SetDB(db)

	// 初始化Redis
	_, err = redis.Init(&cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to init redis: %v", err)
	}

	// 创建模块注册中心
	registry := module.NewRegistry()

	// 注册模块
	registry.Register(userModule.NewModule())

	// 创建服务器并注册模块
	srv := server.New(cfg)
	srv.RegisterModules(registry)

	// 启动服务器
	if err := srv.Run(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
