package handler

import (
	"net/http"
	"strconv"

	"plm/internal/modules/attribute/model"
	"plm/internal/modules/attribute/service"
	"plm/pkg/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

// List 获取属性Schema列表
// @Summary 获取属性Schema列表
// @Tags 属性Schema管理
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param entity_type query string false "实体类型(MATERIAL/DOCUMENT)"
// @Param type_code query string false "类型编码"
// @Param attr_type query string false "属性类型"
// @Param status query string false "状态"
// @Success 200 {object} response.Response
// @Router /attribute-schemas [get]
func (h *Handler) List(c *gin.Context) {
	var query model.AttributeSchemaQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	list, total, err := h.svc.ListAttributeSchemas(&query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取列表失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":  list,
		"total": total,
	})
}

// Get 获取属性Schema详情
// @Summary 获取属性Schema详情
// @Tags 属性Schema管理
// @Accept json
// @Produce json
// @Param id path int true "Schema ID"
// @Success 200 {object} response.Response
// @Router /attribute-schemas/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的ID")
		return
	}

	schema, err := h.svc.GetAttributeSchema(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "Schema不存在")
		return
	}

	response.Success(c, schema)
}

// Create 创建属性Schema
// @Summary 创建属性Schema
// @Tags 属性Schema管理
// @Accept json
// @Produce json
// @Param request body model.CreateAttributeSchemaRequest true "创建请求"
// @Success 200 {object} response.Response
// @Router /attribute-schemas [post]
func (h *Handler) Create(c *gin.Context) {
	var req model.CreateAttributeSchemaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	// 获取当前用户ID
	userID, _ := c.Get("userID")

	schema, err := h.svc.CreateAttributeSchema(&req, userID.(uint))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "创建失败: "+err.Error())
		return
	}

	response.Success(c, schema)
}

// Update 更新属性Schema
// @Summary 更新属性Schema
// @Tags 属性Schema管理
// @Accept json
// @Produce json
// @Param id path int true "Schema ID"
// @Param request body model.UpdateAttributeSchemaRequest true "更新请求"
// @Success 200 {object} response.Response
// @Router /attribute-schemas/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的ID")
		return
	}

	var req model.UpdateAttributeSchemaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	// 获取当前用户ID
	userID, _ := c.Get("userID")

	schema, err := h.svc.UpdateAttributeSchema(uint(id), &req, userID.(uint))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "更新失败: "+err.Error())
		return
	}

	response.Success(c, schema)
}

// Delete 删除属性Schema
// @Summary 删除属性Schema
// @Tags 属性Schema管理
// @Accept json
// @Produce json
// @Param id path int true "Schema ID"
// @Success 200 {object} response.Response
// @Router /attribute-schemas/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的ID")
		return
	}

	if err := h.svc.DeleteAttributeSchema(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// Release 发布属性Schema
// @Summary 发布属性Schema
// @Tags 属性Schema管理
// @Accept json
// @Produce json
// @Param id path int true "Schema ID"
// @Success 200 {object} response.Response
// @Router /attribute-schemas/{id}/release [post]
func (h *Handler) Release(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的ID")
		return
	}

	if err := h.svc.ReleaseAttributeSchema(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "发布失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// Activate 激活指定版本的Schema
// @Summary 激活指定版本的Schema
// @Tags 属性Schema管理
// @Accept json
// @Produce json
// @Param id path int true "Schema ID"
// @Success 200 {object} response.Response
// @Router /attribute-schemas/{id}/activate [post]
func (h *Handler) Activate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的ID")
		return
	}

	if err := h.svc.ActivateSchema(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "激活失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// GetActive 获取当前激活的Schema
// @Summary 获取当前激活的Schema
// @Tags 属性Schema管理
// @Accept json
// @Produce json
// @Param entity_type query string true "实体类型(MATERIAL/DOCUMENT)"
// @Param type_code query string true "类型编码"
// @Param attr_type query string true "属性类型"
// @Success 200 {object} response.Response
// @Router /attribute-schemas/active [get]
func (h *Handler) GetActive(c *gin.Context) {
	entityType := c.Query("entity_type")
	typeCode := c.Query("type_code")
	attrType := c.Query("attr_type")

	if entityType == "" || typeCode == "" || attrType == "" {
		response.Error(c, http.StatusBadRequest, "缺少必要参数")
		return
	}

	schema, err := h.svc.GetActiveSchema(entityType, typeCode, attrType)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取失败: "+err.Error())
		return
	}

	response.Success(c, schema)
}

// GetAllActive 获取指定类型的所有激活Schema
// @Summary 获取指定类型的所有激活Schema
// @Tags 属性Schema管理
// @Accept json
// @Produce json
// @Param entity_type query string true "实体类型(MATERIAL/DOCUMENT)"
// @Param type_code query string true "类型编码"
// @Success 200 {object} response.Response
// @Router /attribute-schemas/all-active [get]
func (h *Handler) GetAllActive(c *gin.Context) {
	entityType := c.Query("entity_type")
	typeCode := c.Query("type_code")

	if entityType == "" || typeCode == "" {
		response.Error(c, http.StatusBadRequest, "缺少必要参数")
		return
	}

	schemas, err := h.svc.GetAllActiveSchemas(entityType, typeCode)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取失败: "+err.Error())
		return
	}

	response.Success(c, schemas)
}
