package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"plm/internal/core/config"
	"plm/pkg/auth"
	"plm/pkg/response"
)

var jwtInstance *auth.JWT

// InitAuth 初始化认证中间件
func InitAuth(cfg *config.JWTConfig) {
	jwtInstance = auth.NewJWT(&auth.JWTConfig{
		Secret:     cfg.Secret,
		ExpireTime: cfg.Expire,
	})
}

// Auth JWT认证中间件
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "未提供认证令牌")
			c.Abort()
			return
		}

		// 解析Bearer Token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "认证令牌格式错误")
			c.Abort()
			return
		}

		token := parts[1]

		// 验证JWT Token
		claims, err := jwtInstance.ParseToken(token)
		if err != nil {
			if err == auth.ErrTokenExpired {
				response.Unauthorized(c, "认证令牌已过期")
			} else {
				response.Unauthorized(c, "认证令牌无效")
			}
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("real_name", claims.RealName)
		c.Set("roles", claims.Roles)
		c.Set("token", token)

		c.Next()
	}
}

// GetUserID 从上下文获取用户ID
func GetUserID(c *gin.Context) uint {
	if userID, exists := c.Get("user_id"); exists {
		return userID.(uint)
	}
	return 0
}

// GetUsername 从上下文获取用户名
func GetUsername(c *gin.Context) string {
	if username, exists := c.Get("username"); exists {
		return username.(string)
	}
	return ""
}

// GetRealName 从上下文获取真实姓名
func GetRealName(c *gin.Context) string {
	if realName, exists := c.Get("real_name"); exists {
		return realName.(string)
	}
	return ""
}

// GetRoles 从上下文获取角色
func GetRoles(c *gin.Context) string {
	if roles, exists := c.Get("roles"); exists {
		return roles.(string)
	}
	return ""
}
