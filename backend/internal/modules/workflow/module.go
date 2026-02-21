package workflow

import (
	"plm/internal/core/middleware"
	"plm/internal/core/module"
	"plm/internal/modules/workflow/handler"
	"plm/internal/modules/workflow/repository"
	"plm/internal/modules/workflow/service"
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
	return "workflow"
}

func (m *Module) Init(container *module.Container) error {
	m.repo = repository.NewRepository(container.DB)
	m.service = service.NewService(m.repo)
	m.service.SetDB(container.DB)
	m.handler = handler.NewHandler(m.service)
	return nil
}

func (m *Module) Routes(r module.RouterGroup) {
	// 流程管理路由（需要登录）
	workflows := r.Group("/workflows", middleware.Auth())
	{
		// 流程定义管理
		workflows.GET("/definitions", m.handler.ListDefinitions)
		workflows.POST("/definitions", m.handler.CreateDefinition)
		workflows.GET("/definitions/:id", m.handler.GetDefinition)
		workflows.PUT("/definitions/:id", m.handler.UpdateDefinition)
		workflows.DELETE("/definitions/:id", m.handler.DeleteDefinition)
		workflows.POST("/definitions/:id/release", m.handler.ReleaseDefinition)
		workflows.GET("/definitions/active", m.handler.GetActiveDefinitionByType)

		// 流程实例管理
		workflows.GET("/instances", m.handler.ListInstances)
		workflows.POST("/instances", m.handler.InitiateWorkflow)
		workflows.GET("/instances/:id", m.handler.GetInstance)

		// 待办管理
		workflows.GET("/todos", m.handler.GetMyTodos)
		workflows.GET("/todos/count", m.handler.GetTodoCount)

		// 审批操作
		workflows.POST("/instances/:id/approve", m.handler.Approve)
		workflows.POST("/instances/:id/reject", m.handler.Reject)
		workflows.POST("/instances/:id/transfer", m.handler.Transfer)
		workflows.POST("/instances/:id/withdraw", m.handler.Withdraw)

		// 流程历史
		workflows.GET("/instances/:id/history", m.handler.GetHistories)
	}
}
