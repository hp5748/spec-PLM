package middleware

import (
	"crypto/subtle"
	"github.com/gin-gonic/gin"
	"plm/pkg/response"
)

const (
	csrfHeader = "X-CSRF-Token"
	csrfCookie = "csrf_token"
)

// CSRF CSRF防护中间件
func CSRF() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 只对需要修改数据的请求进行CSRF检查
		if c.Request.Method != "GET" && c.Request.Method != "HEAD" && c.Request.Method != "OPTIONS" {
			// 从Header获取CSRF Token
			token := c.GetHeader(csrfHeader)
			if token == "" {
				// 也可以从表单获取
				token = c.PostForm("_csrf")
			}

			// 从Cookie获取期望的Token
			cookieToken, err := c.Cookie(csrfCookie)
			if err != nil || cookieToken == "" {
				response.Error(c, 40300, "CSRF Token缺失")
				c.Abort()
				return
			}

			// 使用恒定时间比较防止时序攻击
			if subtle.ConstantTimeCompare([]byte(token), []byte(cookieToken)) != 1 {
				response.Error(c, 40300, "CSRF Token无效")
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// GenerateCSRFToken 生成CSRF Token（在登录时调用）
func GenerateCSRFToken(c *gin.Context, token string) {
	// 设置Cookie
	c.SetCookie(csrfCookie, token, 86400, "/", "", false, true)
	// 同时在响应中返回
	c.Header(csrfHeader, token)
}
