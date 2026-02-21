package main

import (
	"log"

	"plm/internal/core/config"
	"plm/internal/core/module"
	"plm/internal/core/server"
	attributeModule "plm/internal/modules/attribute"
	userModule "plm/internal/modules/user"
	userRepository "plm/internal/modules/user/repository"
	"plm/internal/modules/user/seeder"
	materialModule "plm/internal/modules/material"
	materialRepository "plm/internal/modules/material/repository"
	documentModule "plm/internal/modules/document"
	documentRepository "plm/internal/modules/document/repository"
	bomModule "plm/internal/modules/bom"
	workflowModule "plm/internal/modules/workflow"
	scriptModule "plm/internal/modules/script"
	"plm/pkg/database"
	"plm/pkg/redis"
	"plm/pkg/storage"
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

	// 同步权限（自动添加新权限，不删除已有权限）
	if err := seeder.SyncPermissions(db); err != nil {
		log.Fatalf("Failed to sync permissions: %v", err)
	}

	// 设置Repository的数据库连接
	userRepository.SetDB(db)
	materialRepository.SetDB(db)
	documentRepository.SetDB(db)

	// 初始化Redis
	_, err = redis.Init(&cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to init redis: %v", err)
	}

	// 初始化MinIO存储
	minioStorage, err := storage.NewMinIOStorage(&storage.Config{
		Endpoint:  cfg.Storage.Endpoint,
		AccessKey: cfg.Storage.AccessKey,
		SecretKey: cfg.Storage.SecretKey,
		UseSSL:    cfg.Storage.UseSSL,
		Bucket:    cfg.Storage.Bucket,
	})
	if err != nil {
		log.Fatalf("Failed to init minio storage: %v", err)
	}

	// 创建模块注册中心
	registry := module.NewRegistry()

	// 注册模块
	registry.Register(userModule.NewModule())
	registry.Register(materialModule.NewModule())
	registry.Register(documentModule.NewModule(minioStorage))
	registry.Register(bomModule.NewModule())
	registry.Register(attributeModule.NewModule())
	registry.Register(workflowModule.NewModule())
	registry.Register(scriptModule.NewModule())

	// 创建服务器并注册模块
	srv := server.New(cfg)
	srv.SetDB(db)
	srv.RegisterModules(registry)

	// 启动服务器
	if err := srv.Run(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
