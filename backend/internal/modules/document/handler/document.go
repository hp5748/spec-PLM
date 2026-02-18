package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"plm/internal/modules/document/model"
	"plm/internal/modules/document/service"
	"plm/pkg/response"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

// List 文档列表
func (h *Handler) List(c *gin.Context) {
	var query model.DocumentListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 20
	}

	docs, total, err := h.svc.ListDocuments(c, &query)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessPage(c, docs, total, query.Page, query.PageSize)
}

// Get 获取文档详情
func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的文档ID")
		return
	}

	doc, err := h.svc.GetDocument(c, uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, doc)
}

// Upload 上传文档
func (h *Handler) Upload(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "请上传文件")
		return
	}

	// 检查文件大小（100MB）
	if file.Size > 100*1024*1024 {
		response.BadRequest(c, "文件大小不能超过100MB")
		return
	}

	// 获取表单数据
	req := &model.CreateDocumentRequest{
		DocID:      c.PostForm("doc_id"),
		DocName:    c.PostForm("doc_name"),
		DocType:    c.PostForm("doc_type"),
		Attributes: c.PostForm("attributes"),
	}

	// 如果没有提供doc_type，从文件名推断
	if req.DocType == "" {
		req.DocType = service.GetDocTypeFromExt(file.Filename)
	}

	// 验证必填字段
	if req.DocID == "" {
		response.BadRequest(c, "文档编码不能为空")
		return
	}
	if req.DocName == "" {
		req.DocName = file.Filename
	}

	// 从上下文获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	doc, err := h.svc.UploadDocument(c, req, file, userID.(uint))
	if err != nil {
		response.Error(c, 40000, err.Error())
		return
	}

	response.SuccessWithMessage(c, "上传成功", doc)
}

// Update 更新文档
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的文档ID")
		return
	}

	var req model.UpdateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	doc, err := h.svc.UpdateDocument(c, uint(id), &req)
	if err != nil {
		response.Error(c, 40000, err.Error())
		return
	}

	response.SuccessWithMessage(c, "更新成功", doc)
}

// Delete 删除文档
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的文档ID")
		return
	}

	if err := h.svc.DeleteDocument(c, uint(id)); err != nil {
		response.Error(c, 40000, err.Error())
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

// Download 下载文档
func (h *Handler) Download(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的文档ID")
		return
	}

	url, err := h.svc.GetDocumentURL(c, uint(id))
	if err != nil {
		response.Error(c, 40000, err.Error())
		return
	}

	response.Success(c, gin.H{
		"download_url": url,
	})
}

// Search 搜索文档
func (h *Handler) Search(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		response.BadRequest(c, "请输入搜索关键字")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	docs, total, err := h.svc.SearchDocuments(c, keyword, page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessPage(c, docs, total, page, pageSize)
}

// IncrementVersion 文档版本升级
func (h *Handler) IncrementVersion(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的文档ID")
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	c.ShouldBindJSON(&req)

	userID := c.GetUint("user_id")

	doc, err := h.svc.IncrementVersion(c, uint(id), req.Reason, userID)
	if err != nil {
		response.Error(c, 40000, err.Error())
		return
	}

	response.SuccessWithMessage(c, "版本升级成功", doc)
}

// GetVersionHistory 获取文档版本历史
func (h *Handler) GetVersionHistory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的文档ID")
		return
	}

	histories, err := h.svc.GetVersionHistory(c, uint(id))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, histories)
}

// Preview 文档预览
func (h *Handler) Preview(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的文档ID")
		return
	}

	url, err := h.svc.GetPreviewURL(c, uint(id))
	if err != nil {
		response.Error(c, 40000, err.Error())
		return
	}

	response.Success(c, gin.H{
		"preview_url": url,
	})
}
