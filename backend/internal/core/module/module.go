package module

import (
	"gorm.io/gorm"
)

// Container 依赖容器，管理模块间的依赖
type Container struct {
	DB     *gorm.DB
	Config *Config
}

// Config 容器配置
type Config struct {
	JWTSecret     string
	JWTExpire     interface{} // time.Duration
	RedisAddr     string
	RedisPassword string
	RedisDB       int
}

func NewContainer(db *gorm.DB) *Container {
	return &Container{
		DB: db,
	}
}

func (c *Container) SetConfig(cfg *Config) {
	c.Config = cfg
}

// Module 模块接口定义
type Module interface {
	// Name 模块名称
	Name() string
	// Init 初始化模块依赖
	Init(container *Container) error
	// Routes 注册路由
	Routes(r RouterGroup)
}

// RouterGroup 路由组接口（适配gin.RouterGroup）
type RouterGroup interface {
	Group(relativePath string, handlers ...interface{}) RouterGroup
	GET(relativePath string, handlers ...interface{})
	POST(relativePath string, handlers ...interface{})
	PUT(relativePath string, handlers ...interface{})
	PATCH(relativePath string, handlers ...interface{})
	DELETE(relativePath string, handlers ...interface{})
}
