package middleware

import (
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

// XSSProtection XSS防护中间件
func XSSProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置安全相关的HTTP头
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Content-Security-Policy", "default-src 'self'")

		c.Next()
	}
}

// SanitizeInput 清理输入中的潜在XSS攻击字符串
var xssPatterns = []*regexp.Regexp{
	regexp.MustCompile(`<script[^>]*>.*?</script>`),
	regexp.MustCompile(`<script[^>]*>`),
	regexp.MustCompile(`javascript:`),
	regexp.MustCompile(`on\w+\s*=`),
	regexp.MustCompile(`<iframe[^>]*>`),
	regexp.MustCompile(`<object[^>]*>`),
	regexp.MustCompile(`<embed[^>]*>`),
	regexp.MustCompile(`<img[^>]+onerror\s*=`),
}

// SanitizeString 清理字符串中的XSS攻击代码
func SanitizeString(input string) string {
	output := input
	for _, pattern := range xssPatterns {
		output = pattern.ReplaceAllString(output, "")
	}
	return strings.TrimSpace(output)
}

// XSSSanitizeMiddleware XSS清理中间件
func XSSSanitizeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 只处理POST/PUT/PATCH请求
		method := c.Request.Method
		if method == "POST" || method == "PUT" || method == "PATCH" {
			// 读取请求体并清理
			// 注意：这里简化处理，实际项目中可能需要更复杂的处理逻辑
			// 对于JSON请求体，应在绑定后逐字段清理
		}
		c.Next()
	}
}
