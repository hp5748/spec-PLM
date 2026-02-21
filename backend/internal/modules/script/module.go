package script

import (
	"plm/internal/core/middleware"
	"plm/internal/core/module"
	"plm/internal/modules/script/handler"
	"plm/internal/modules/script/model"
	"plm/internal/modules/script/repository"
	"plm/internal/modules/script/service"

	"gorm.io/gorm"
)

// Module 脚本模块
type Module struct {
	db      *gorm.DB
	handler *handler.Handler
	service *service.Service
	repo    *repository.Repository
}

// NewModule 创建脚本模块
func NewModule() *Module {
	return &Module{}
}

// Name 模块名称
func (m *Module) Name() string {
	return "script"
}

// Init 初始化模块
func (m *Module) Init(container *module.Container) error {
	m.db = container.DB
	m.repo = repository.NewRepository(container.DB)
	m.service = service.NewService(m.repo)
	m.service.SetDB(container.DB)
	m.handler = handler.NewHandler(m.service)

	// 自动迁移数据库表
	return m.db.AutoMigrate(
		&model.Script{},
		&model.ScriptExecutionLog{},
	)
}

// Routes 注册路由
func (m *Module) Routes(r module.RouterGroup) {
	// 脚本管理路由（需要登录）
	scripts := r.Group("/scripts", middleware.Auth())
	{
		// 脚本管理
		scripts.GET("", m.handler.List)
		scripts.POST("", m.handler.Create)
		scripts.GET("/:id", m.handler.Get)
		scripts.PUT("/:id", m.handler.Update)
		scripts.DELETE("/:id", m.handler.Delete)
		scripts.POST("/:id/release", m.handler.Release)
		scripts.POST("/:id/obsolete", m.handler.Obsolete)

		// 脚本执行
		scripts.POST("/execute", m.handler.Execute)
		scripts.POST("/test", m.handler.Test)
	}

	// 执行日志
	logs := r.Group("/script-logs", middleware.Auth())
	{
		logs.GET("", m.handler.ListLogs)
		logs.GET("/:id", m.handler.GetLog)
	}
}

// GetService 获取服务
func (m *Module) GetService() *service.Service {
	return m.service
}
