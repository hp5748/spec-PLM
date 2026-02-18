package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"plm/internal/modules/material/model"
	"plm/internal/modules/material/service"
	"plm/pkg/response"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

// List 物料列表
func (h *Handler) List(c *gin.Context) {
	var query model.MaterialListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 设置默认值
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 20
	}

	materials, total, err := h.svc.ListMaterials(c, &query)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessPage(c, materials, total, query.Page, query.PageSize)
}

// Get 获取物料详情
func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的物料ID")
		return
	}

	material, err := h.svc.GetMaterial(c, uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, material)
}

// Create 创建物料
func (h *Handler) Create(c *gin.Context) {
	var req model.CreateMaterialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 从上下文获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	material, err := h.svc.CreateMaterial(c, &req, userID.(uint))
	if err != nil {
		response.Error(c, 40000, err.Error())
		return
	}

	response.SuccessWithMessage(c, "创建成功", material)
}

// Update 更新物料
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的物料ID")
		return
	}

	var req model.UpdateMaterialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	material, err := h.svc.UpdateMaterial(c, uint(id), &req)
	if err != nil {
		response.Error(c, 40000, err.Error())
		return
	}

	response.SuccessWithMessage(c, "更新成功", material)
}

// Delete 删除物料
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的物料ID")
		return
	}

	if err := h.svc.DeleteMaterial(c, uint(id)); err != nil {
		response.Error(c, 40000, err.Error())
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

// Search 搜索物料
func (h *Handler) Search(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		response.BadRequest(c, "请输入搜索关键字")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	materials, total, err := h.svc.SearchMaterials(c, keyword, page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessPage(c, materials, total, page, pageSize)
}

// IncrementVersion 版本升级
func (h *Handler) IncrementVersion(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的物料ID")
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	c.ShouldBindJSON(&req)

	userID := c.GetUint("user_id")

	material, err := h.svc.IncrementVersion(c, uint(id), req.Reason, userID)
	if err != nil {
		response.Error(c, 40000, err.Error())
		return
	}

	response.SuccessWithMessage(c, "版本升级成功", material)
}

// GetVersionHistory 获取物料版本历史
func (h *Handler) GetVersionHistory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的物料ID")
		return
	}

	histories, err := h.svc.GetVersionHistory(c, uint(id))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, histories)
}
