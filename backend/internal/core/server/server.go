package server

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"plm/internal/core/config"
	"plm/internal/core/middleware"
	"plm/internal/core/module"
)

type Server struct {
	engine   *gin.Engine
	config   *config.Config
	registry *module.Registry
}

func New(cfg *config.Config) *Server {
	// 设置Gin运行模式
	gin.SetMode(cfg.Server.Mode)

	engine := gin.New()

	// 注册全局中间件
	engine.Use(gin.Recovery())
	engine.Use(middleware.Logger())
	engine.Use(middleware.CORS())
	engine.Use(middleware.RequestID())
	engine.Use(middleware.XSSProtection())

	// 初始化认证中间件
	middleware.InitAuth(&cfg.JWT)

	return &Server{
		engine: engine,
		config: cfg,
	}
}

func (s *Server) RegisterModules(registry *module.Registry) {
	s.registry = registry
}

func (s *Server) Run() error {
	// 健康检查
	s.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API路由组
	api := s.engine.Group("/api/v1")

	// 注册所有模块的路由
	if s.registry != nil {
		container := module.NewContainer(nil) // TODO: 注入DB
		// 设置配置
		container.SetConfig(&module.Config{
			JWTSecret: s.config.JWT.Secret,
			JWTExpire: s.config.JWT.Expire,
		})
		for _, m := range s.registry.All() {
			if err := m.Init(container); err != nil {
				log.Fatalf("Failed to init module %s: %v", m.Name(), err)
			}
			m.Routes(&ginRouterGroup{api.Group(m.Name())})
			fmt.Printf("Module %s registered\n", m.Name())
		}
	}

	addr := fmt.Sprintf(":%d", s.config.Server.Port)
	fmt.Printf("Server starting on %s\n", addr)
	return s.engine.Run(addr)
}

// ginRouterGroup 适配gin.RouterGroup到module.RouterGroup接口
type ginRouterGroup struct {
	*gin.RouterGroup
}

func (g *ginRouterGroup) Group(relativePath string, handlers ...interface{}) module.RouterGroup {
	ginHandlers := make([]gin.HandlerFunc, 0, len(handlers))
	for _, h := range handlers {
		if gh, ok := h.(gin.HandlerFunc); ok {
			ginHandlers = append(ginHandlers, gh)
		}
	}
	return &ginRouterGroup{g.RouterGroup.Group(relativePath, ginHandlers...)}
}

func (g *ginRouterGroup) GET(relativePath string, handlers ...interface{}) {
	ginHandlers := toGinHandlers(handlers)
	g.RouterGroup.GET(relativePath, ginHandlers...)
}

func (g *ginRouterGroup) POST(relativePath string, handlers ...interface{}) {
	ginHandlers := toGinHandlers(handlers)
	g.RouterGroup.POST(relativePath, ginHandlers...)
}

func (g *ginRouterGroup) PUT(relativePath string, handlers ...interface{}) {
	ginHandlers := toGinHandlers(handlers)
	g.RouterGroup.PUT(relativePath, ginHandlers...)
}

func (g *ginRouterGroup) PATCH(relativePath string, handlers ...interface{}) {
	ginHandlers := toGinHandlers(handlers)
	g.RouterGroup.PATCH(relativePath, ginHandlers...)
}

func (g *ginRouterGroup) DELETE(relativePath string, handlers ...interface{}) {
	ginHandlers := toGinHandlers(handlers)
	g.RouterGroup.DELETE(relativePath, ginHandlers...)
}

func toGinHandlers(handlers []interface{}) []gin.HandlerFunc {
	ginHandlers := make([]gin.HandlerFunc, 0, len(handlers))
	for _, h := range handlers {
		// 首先尝试直接类型断言
		if gh, ok := h.(gin.HandlerFunc); ok {
			ginHandlers = append(ginHandlers, gh)
			continue
		}
		// 尝试断言为 func(*gin.Context) 然后转换为 gin.HandlerFunc
		if fn, ok := h.(func(*gin.Context)); ok {
			ginHandlers = append(ginHandlers, fn)
		}
	}
	return ginHandlers
}
