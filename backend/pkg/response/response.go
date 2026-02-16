package response

import (
	"time"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
}

// PageData 分页数据结构
type PageData struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Code:      0,
		Message:   "success",
		Data:      data,
		Timestamp: time.Now().UnixMilli(),
	})
}

// 成功响应带消息
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(200, Response{
		Code:      0,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UnixMilli(),
	})
}

// 分页成功响应
func SuccessPage(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	c.JSON(200, Response{
		Code:    0,
		Message: "success",
		Data: PageData{
			List:     list,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
		Timestamp: time.Now().UnixMilli(),
	})
}

// 错误响应
func Error(c *gin.Context, code int, message string) {
	c.JSON(200, Response{
		Code:      code,
		Message:   message,
		Data:      nil,
		Timestamp: time.Now().UnixMilli(),
	})
}

// 参数错误
func BadRequest(c *gin.Context, message string) {
	Error(c, 40000, message)
}

// 未认证
func Unauthorized(c *gin.Context, message string) {
	c.JSON(401, Response{
		Code:      40100,
		Message:   message,
		Data:      nil,
		Timestamp: time.Now().UnixMilli(),
	})
}

// 无权限
func Forbidden(c *gin.Context, message string) {
	c.JSON(403, Response{
		Code:      40300,
		Message:   message,
		Data:      nil,
		Timestamp: time.Now().UnixMilli(),
	})
}

// 资源不存在
func NotFound(c *gin.Context, message string) {
	Error(c, 40400, message)
}

// 服务器错误
func ServerError(c *gin.Context, message string) {
	Error(c, 50000, message)
}
