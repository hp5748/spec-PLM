package attribute

import (
	"plm/internal/core/middleware"
	"plm/internal/core/module"
	"plm/internal/modules/attribute/handler"
	"plm/internal/modules/attribute/repository"
	"plm/internal/modules/attribute/service"
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
	return "attribute"
}

func (m *Module) Init(container *module.Container) error {
	// 初始化Repository层
	m.repository = repository.NewRepository()
	m.repository.SetDB(container.DB)

	// 初始化Service层
	m.service = service.NewService(m.repository)

	// 初始化Handler层
	m.handler = handler.NewHandler(m.service)

	return nil
}

func (m *Module) Routes(r module.RouterGroup) {
	// 属性Schema管理路由（需要登录 + 管理员权限）
	schemas := r.Group("/attribute-schemas", middleware.Auth())
	{
		schemas.GET("", m.handler.List)
		schemas.GET("/active", m.handler.GetActive)
		schemas.GET("/all-active", m.handler.GetAllActive)
		schemas.GET("/:id", m.handler.Get)
		schemas.POST("", m.handler.Create)
		schemas.PUT("/:id", m.handler.Update)
		schemas.DELETE("/:id", m.handler.Delete)
		schemas.POST("/:id/release", m.handler.Release)
		schemas.POST("/:id/activate", m.handler.Activate)
	}
}
