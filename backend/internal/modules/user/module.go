package user

import (
	"time"

	"plm/internal/core/module"
	"plm/internal/modules/user/handler"
	"plm/internal/modules/user/repository"
	"plm/internal/modules/user/service"
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
	return "user"
}

func (m *Module) Init(container *module.Container) error {
	// 初始化Repository层
	m.repository = repository.NewRepository()

	// 从容器获取配置
	jwtSecret := "plm-default-secret-key"
	jwtExpire := 24 * time.Hour // 默认一天

	if container.Config != nil {
		if container.Config.JWTSecret != "" {
			jwtSecret = container.Config.JWTSecret
		}
		if expire, ok := container.Config.JWTExpire.(time.Duration); ok {
			jwtExpire = expire
		}
	}

	// 初始化Service层
	m.service = service.NewService(m.repository, jwtSecret, jwtExpire)

	// 初始化Handler层
	m.handler = handler.NewHandler(m.service)

	return nil
}

func (m *Module) Routes(r module.RouterGroup) {
	// 认证相关路由（无需登录）
	auth := r.Group("/auth")
	{
		auth.POST("/login", m.handler.Login)
		auth.POST("/register", m.handler.Register)
		auth.POST("/refresh", m.handler.RefreshToken)
	}

	// 用户管理路由（需要登录）
	users := r.Group("/users")
	{
		users.GET("", m.handler.List)
		users.POST("", m.handler.Create)
		users.GET("/:id", m.handler.Get)
		users.PUT("/:id", m.handler.Update)
		users.DELETE("/:id", m.handler.Delete)
	}

	// 组织管理路由
	orgs := r.Group("/organizations")
	{
		orgs.GET("", m.handler.ListOrganizations)
		orgs.POST("", m.handler.CreateOrganization)
		orgs.GET("/:id", m.handler.GetOrganization)
		orgs.PUT("/:id", m.handler.UpdateOrganization)
		orgs.DELETE("/:id", m.handler.DeleteOrganization)
	}

	// 部门管理路由
	depts := r.Group("/departments")
	{
		depts.GET("", m.handler.ListDepartments)
		depts.POST("", m.handler.CreateDepartment)
		depts.GET("/:id", m.handler.GetDepartment)
		depts.PUT("/:id", m.handler.UpdateDepartment)
		depts.DELETE("/:id", m.handler.DeleteDepartment)
	}

	// 角色管理路由
	roles := r.Group("/roles")
	{
		roles.GET("", m.handler.ListRoles)
		roles.POST("", m.handler.CreateRole)
		roles.GET("/:id", m.handler.GetRole)
		roles.PUT("/:id", m.handler.UpdateRole)
		roles.DELETE("/:id", m.handler.DeleteRole)
	}
}
