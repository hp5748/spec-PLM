package handler

import (
	"strconv"

	"plm/internal/modules/script/model"
	"plm/internal/modules/script/service"
	"plm/pkg/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

// ==================== 脚本管理 ====================

// List 脚本列表
// @Summary 脚本列表
// @Tags 脚本管理
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param type query string false "脚本类型"
// @Param trigger_type query string false "触发类型"
// @Param business_type query string false "业务类型"
// @Param status query string false "状态"
// @Param keyword query string false "关键字"
// @Success 200 {object} response.Response
// @Router /api/v1/script/scripts [get]
func (h *Handler) List(c *gin.Context) {
	var query model.ScriptListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, 40001, "参数错误: "+err.Error())
		return
	}

	scripts, total, err := h.service.ListScripts(c.Request.Context(), &query)
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":  scripts,
		"total": total,
		"page":  query.Page,
		"page_size": query.PageSize,
	})
}

// Get 获取脚本详情
// @Summary 获取脚本详情
// @Tags 脚本管理
// @Accept json
// @Produce json
// @Param id path int true "脚本ID"
// @Success 200 {object} response.Response
// @Router /api/v1/script/scripts/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 40001, "无效的ID")
		return
	}

	script, err := h.service.GetScript(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}

	response.Success(c, script)
}

// Create 创建脚本
// @Summary 创建脚本
// @Tags 脚本管理
// @Accept json
// @Produce json
// @Param body body model.CreateScriptRequest true "请求体"
// @Success 200 {object} response.Response
// @Router /api/v1/script/scripts [post]
func (h *Handler) Create(c *gin.Context) {
	var req model.CreateScriptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, "参数错误: "+err.Error())
		return
	}

	userID := c.GetUint("user_id")
	script, err := h.service.CreateScript(c.Request.Context(), &req, userID)
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}

	response.Success(c, script)
}

// Update 更新脚本
// @Summary 更新脚本
// @Tags 脚本管理
// @Accept json
// @Produce json
// @Param id path int true "脚本ID"
// @Param body body model.UpdateScriptRequest true "请求体"
// @Success 200 {object} response.Response
// @Router /api/v1/script/scripts/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 40001, "无效的ID")
		return
	}

	var req model.UpdateScriptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, "参数错误: "+err.Error())
		return
	}

	script, err := h.service.UpdateScript(c.Request.Context(), uint(id), &req)
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}

	response.Success(c, script)
}

// Release 发布脚本
// @Summary 发布脚本
// @Tags 脚本管理
// @Accept json
// @Produce json
// @Param id path int true "脚本ID"
// @Success 200 {object} response.Response
// @Router /api/v1/script/scripts/{id}/release [post]
func (h *Handler) Release(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 40001, "无效的ID")
		return
	}

	script, err := h.service.ReleaseScript(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}

	response.Success(c, script)
}

// Obsolete 废弃脚本
// @Summary 废弃脚本
// @Tags 脚本管理
// @Accept json
// @Produce json
// @Param id path int true "脚本ID"
// @Success 200 {object} response.Response
// @Router /api/v1/script/scripts/{id}/obsolete [post]
func (h *Handler) Obsolete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 40001, "无效的ID")
		return
	}

	script, err := h.service.ObsoleteScript(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}

	response.Success(c, script)
}

// Delete 删除脚本
// @Summary 删除脚本
// @Tags 脚本管理
// @Accept json
// @Produce json
// @Param id path int true "脚本ID"
// @Success 200 {object} response.Response
// @Router /api/v1/script/scripts/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 40001, "无效的ID")
		return
	}

	if err := h.service.DeleteScript(c.Request.Context(), uint(id)); err != nil {
		response.Error(c, 50001, err.Error())
		return
	}

	response.Success(c, nil)
}

// ==================== 脚本执行 ====================

// Execute 执行脚本
// @Summary 执行脚本
// @Tags 脚本管理
// @Accept json
// @Produce json
// @Param body body model.ExecuteScriptRequest true "请求体"
// @Success 200 {object} response.Response
// @Router /api/v1/script/scripts/execute [post]
func (h *Handler) Execute(c *gin.Context) {
	var req model.ExecuteScriptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, "参数错误: "+err.Error())
		return
	}

	userID := c.GetUint("user_id")
	result, err := h.service.ExecuteScript(
		c.Request.Context(),
		req.ScriptID,
		req.BusinessType,
		req.BusinessID,
		req.Context,
		userID,
	)
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}

	response.Success(c, result)
}

// Test 测试脚本
// @Summary 测试脚本
// @Tags 脚本管理
// @Accept json
// @Produce json
// @Param body body model.ScriptTestRequest true "请求体"
// @Success 200 {object} response.Response
// @Router /api/v1/script/scripts/test [post]
func (h *Handler) Test(c *gin.Context) {
	var req model.ScriptTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, "参数错误: "+err.Error())
		return
	}

	result, err := h.service.TestScript(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}

	response.Success(c, result)
}

// ==================== 执行日志 ====================

// ListLogs 执行日志列表
// @Summary 执行日志列表
// @Tags 脚本管理
// @Accept json
// @Produce json
// @Param script_id query int false "脚本ID"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} response.Response
// @Router /api/v1/script/logs [get]
func (h *Handler) ListLogs(c *gin.Context) {
	scriptID, _ := strconv.ParseUint(c.Query("script_id"), 10, 64)
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))

	logs, total, err := h.service.GetExecutionLogs(c.Request.Context(), uint(scriptID), page, pageSize)
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":  logs,
		"total": total,
		"page":  page,
		"page_size": pageSize,
	})
}

// GetLog 获取执行日志详情
// @Summary 获取执行日志详情
// @Tags 脚本管理
// @Accept json
// @Produce json
// @Param id path int true "日志ID"
// @Success 200 {object} response.Response
// @Router /api/v1/script/logs/{id} [get]
func (h *Handler) GetLog(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 40001, "无效的ID")
		return
	}

	log, err := h.service.GetExecutionLog(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}

	response.Success(c, log)
}

// RegisterRoutes 注册路由
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	scripts := r.Group("/scripts")
	{
		scripts.GET("", h.List)
		scripts.POST("", h.Create)
		scripts.GET("/:id", h.Get)
		scripts.PUT("/:id", h.Update)
		scripts.DELETE("/:id", h.Delete)
		scripts.POST("/:id/release", h.Release)
		scripts.POST("/:id/obsolete", h.Obsolete)
		scripts.POST("/execute", h.Execute)
		scripts.POST("/test", h.Test)
	}

	logs := r.Group("/logs")
	{
		logs.GET("", h.ListLogs)
		logs.GET("/:id", h.GetLog)
	}
}
