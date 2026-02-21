package handler

import (
	"strconv"

	"plm/internal/modules/workflow/model"
	"plm/internal/modules/workflow/service"
	"plm/pkg/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

// ==================== 流程定义相关 ====================

// ListDefinitions 流程定义列表
func (h *Handler) ListDefinitions(c *gin.Context) {
	var query model.WorkflowDefinitionListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, 40001, "参数错误: "+err.Error())
		return
	}

	list, total, err := h.service.ListDefinitions(c.Request.Context(), &query)
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":       list,
		"total":      total,
		"page":       query.Page,
		"page_size":  query.PageSize,
	})
}

// GetDefinition 获取流程定义详情
func (h *Handler) GetDefinition(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 40001, "无效的ID")
		return
	}

	def, err := h.service.GetDefinition(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, 40401, err.Error())
		return
	}

	response.Success(c, def)
}

// CreateDefinition 创建流程定义
func (h *Handler) CreateDefinition(c *gin.Context) {
	var req model.CreateWorkflowDefinitionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, "参数错误: "+err.Error())
		return
	}

	userID := c.GetUint("user_id")
	def, err := h.service.CreateDefinition(c.Request.Context(), &req, userID)
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}

	response.Success(c, def)
}

// UpdateDefinition 更新流程定义
func (h *Handler) UpdateDefinition(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 40001, "无效的ID")
		return
	}

	var req model.UpdateWorkflowDefinitionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, "参数错误: "+err.Error())
		return
	}

	def, err := h.service.UpdateDefinition(c.Request.Context(), uint(id), &req)
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}

	response.Success(c, def)
}

// ReleaseDefinition 发布流程定义
func (h *Handler) ReleaseDefinition(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 40001, "无效的ID")
		return
	}

	def, err := h.service.ReleaseDefinition(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}

	response.Success(c, def)
}

// DeleteDefinition 删除流程定义
func (h *Handler) DeleteDefinition(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 40001, "无效的ID")
		return
	}

	if err := h.service.DeleteDefinition(c.Request.Context(), uint(id)); err != nil {
		response.Error(c, 50001, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetActiveDefinitionByType 获取指定类型的激活流程定义
func (h *Handler) GetActiveDefinitionByType(c *gin.Context) {
	workflowType := c.Query("type")
	if workflowType == "" {
		response.Error(c, 40001, "缺少类型参数")
		return
	}

	def, err := h.service.GetActiveDefinitionByType(c.Request.Context(), workflowType)
	if err != nil {
		response.Error(c, 40401, "未找到对应的流程定义")
		return
	}

	response.Success(c, def)
}

// ==================== 流程实例相关 ====================

// ListInstances 流程实例列表
func (h *Handler) ListInstances(c *gin.Context) {
	var query model.WorkflowListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, 40001, "参数错误: "+err.Error())
		return
	}

	list, total, err := h.service.ListInstances(c.Request.Context(), &query)
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":       list,
		"total":      total,
		"page":       query.Page,
		"page_size":  query.PageSize,
	})
}

// GetInstance 获取流程实例详情
func (h *Handler) GetInstance(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 40001, "无效的ID")
		return
	}

	instance, err := h.service.GetInstance(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, 40401, err.Error())
		return
	}

	response.Success(c, instance)
}

// InitiateWorkflow 发起流程
func (h *Handler) InitiateWorkflow(c *gin.Context) {
	var req model.InitiateWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, "参数错误: "+err.Error())
		return
	}

	userID := c.GetUint("user_id")
	instance, err := h.service.InitiateWorkflow(c.Request.Context(), &req, userID)
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}

	response.Success(c, instance)
}

// GetMyTodos 获取我的待办
func (h *Handler) GetMyTodos(c *gin.Context) {
	var query model.TodoListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, 40001, "参数错误: "+err.Error())
		return
	}

	userID := c.GetUint("user_id")
	list, total, err := h.service.GetMyTodos(c.Request.Context(), userID, &query)
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":       list,
		"total":      total,
		"page":       query.Page,
		"page_size":  query.PageSize,
	})
}

// GetTodoCount 获取待办数量
func (h *Handler) GetTodoCount(c *gin.Context) {
	userID := c.GetUint("user_id")
	count, err := h.service.GetTodoCount(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}

	response.Success(c, gin.H{
		"count": count,
	})
}

// GetHistories 获取流程历史
func (h *Handler) GetHistories(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 40001, "无效的ID")
		return
	}

	histories, err := h.service.GetHistories(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}

	response.Success(c, histories)
}

// Approve 同意审批
func (h *Handler) Approve(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 40001, "无效的ID")
		return
	}

	var req model.ApproveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, "参数错误: "+err.Error())
		return
	}

	userID := c.GetUint("user_id")
	instance, err := h.service.Approve(c.Request.Context(), uint(id), userID, &req)
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}

	response.Success(c, instance)
}

// Reject 驳回审批
func (h *Handler) Reject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 40001, "无效的ID")
		return
	}

	var req model.RejectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, "参数错误: "+err.Error())
		return
	}

	userID := c.GetUint("user_id")
	instance, err := h.service.Reject(c.Request.Context(), uint(id), userID, &req)
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}

	response.Success(c, instance)
}

// Transfer 转交审批
func (h *Handler) Transfer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 40001, "无效的ID")
		return
	}

	var req model.TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, "参数错误: "+err.Error())
		return
	}

	userID := c.GetUint("user_id")
	instance, err := h.service.Transfer(c.Request.Context(), uint(id), userID, &req)
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}

	response.Success(c, instance)
}

// Withdraw 撤回流程
func (h *Handler) Withdraw(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 40001, "无效的ID")
		return
	}

	userID := c.GetUint("user_id")
	instance, err := h.service.Withdraw(c.Request.Context(), uint(id), userID)
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}

	response.Success(c, instance)
}
