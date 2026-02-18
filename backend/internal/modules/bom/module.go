package bom

import (
	"plm/internal/core/middleware"
	"plm/internal/core/module"
	"plm/internal/modules/bom/handler"
	"plm/internal/modules/bom/repository"
	"plm/internal/modules/bom/service"
)

type Module struct {
	handler *handler.Handler
	service *service.Service
	repo    *repository.Repository
}

func NewModule() *Module {
	return &Module{}
}

func (m *Module) Name() string {
	return "bom"
}

func (m *Module) Init(container *module.Container) error {
	m.repo = repository.NewRepository(container.DB)
	m.service = service.NewService(m.repo)
	m.service.SetDB(container.DB)
	m.handler = handler.NewHandler(m.service)
	return nil
}

func (m *Module) Routes(r module.RouterGroup) {
	// BOM管理路由（需要登录）
	boms := r.Group("/boms", middleware.Auth())
	{
		// BOM视图管理
		boms.GET("", m.handler.List)
		boms.POST("", m.handler.Create)
		boms.GET("/:id", m.handler.Get)
		boms.PUT("/:id", m.handler.Update)
		boms.DELETE("/:id", m.handler.Delete)

		// BOM树形结构
		boms.GET("/:id/tree", m.handler.GetTree)

		// BOM项管理
		boms.POST("/:id/items", m.handler.AddItem)
		boms.PUT("/:id/items/:itemId", m.handler.UpdateItem)
		boms.DELETE("/:id/items/:itemId", m.handler.DeleteItem)

		// BOM类型转换
		boms.POST("/:id/convert", m.handler.Convert)

		// 导入导出
		boms.GET("/:id/export", m.handler.Export)
		boms.POST("/import", m.handler.Import)
		boms.GET("/template", m.handler.DownloadTemplate)
	}
}
