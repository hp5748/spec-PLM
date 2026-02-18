package material

import (
	"plm/internal/core/middleware"
	"plm/internal/core/module"
	"plm/internal/modules/material/handler"
	"plm/internal/modules/material/repository"
	"plm/internal/modules/material/service"
)

type Module struct {
	handler    *handler.Handler
	service    *service.Service
	repository *repository.Repository
}

func NewModule() *Module {
	return &Module{}
}

func (m *Module) Name() string {
	return "material"
}

func (m *Module) Init(container *module.Container) error {
	// 初始化Repository层
	m.repository = repository.NewRepository()

	// 初始化Service层
	m.service = service.NewService(m.repository)
	m.service.SetDB(container.DB)

	// 初始化Handler层
	m.handler = handler.NewHandler(m.service)

	return nil
}

func (m *Module) Routes(r module.RouterGroup) {
	// 物料管理路由（需要登录）
	materials := r.Group("/materials", middleware.Auth())
	{
		materials.GET("", m.handler.List)
		materials.POST("", m.handler.Create)
		materials.GET("/search", m.handler.Search)
		materials.GET("/:id", m.handler.Get)
		materials.PUT("/:id", m.handler.Update)
		materials.DELETE("/:id", m.handler.Delete)
		materials.POST("/:id/version", m.handler.IncrementVersion)
		materials.GET("/:id/versions", m.handler.GetVersionHistory)
	}
}
