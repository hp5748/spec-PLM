package document

import (
	"plm/internal/core/middleware"
	"plm/internal/core/module"
	"plm/internal/modules/document/handler"
	"plm/internal/modules/document/repository"
	"plm/internal/modules/document/service"
	"plm/pkg/storage"
)

type Module struct {
	handler    *handler.Handler
	service    *service.Service
	repository *repository.Repository
	storage    *storage.MinIOStorage
}

func NewModule(storage *storage.MinIOStorage) *Module {
	return &Module{storage: storage}
}

func (m *Module) Name() string {
	return "document"
}

func (m *Module) Init(container *module.Container) error {
	// 初始化Repository层
	m.repository = repository.NewRepository()

	// 初始化Service层
	m.service = service.NewService(m.repository, m.storage)
	m.service.SetDB(container.DB)

	// 初始化Handler层
	m.handler = handler.NewHandler(m.service)

	return nil
}

func (m *Module) Routes(r module.RouterGroup) {
	// 文档管理路由（需要登录）
	documents := r.Group("/documents", middleware.Auth())
	{
		documents.GET("", m.handler.List)
		documents.POST("", m.handler.Upload)
		documents.GET("/search", m.handler.Search)
		documents.GET("/:id", m.handler.Get)
		documents.PUT("/:id", m.handler.Update)
		documents.DELETE("/:id", m.handler.Delete)
		documents.GET("/:id/download", m.handler.Download)
		documents.GET("/:id/preview", m.handler.Preview)
		documents.POST("/:id/version", m.handler.IncrementVersion)
		documents.GET("/:id/versions", m.handler.GetVersionHistory)
	}
}
