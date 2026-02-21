package handler

import (
	"net/http"
	"strconv"

	"plm/internal/modules/bom/model"
	"plm/internal/modules/bom/service"
	"plm/pkg/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

// List BOM视图列表
func (h *Handler) List(c *gin.Context) {
	var query model.BOMListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, 40001, "参数错误: "+err.Error())
		return
	}

	list, total, err := h.service.ListBOMViews(c.Request.Context(), &query)
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

// Create 创建BOM视图
func (h *Handler) Create(c *gin.Context) {
	var req model.CreateBOMViewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, "参数错误: "+err.Error())
		return
	}

	userID := c.GetUint("user_id")
	bom, err := h.service.CreateBOMView(c.Request.Context(), &req, userID)
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}

	response.Success(c, bom)
}

// Get 获取BOM视图详情
func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 40001, "无效的ID")
		return
	}

	bom, err := h.service.GetBOMView(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, 40401, err.Error())
		return
	}

	response.Success(c, bom)
}

// Update 更新BOM视图
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 40001, "无效的ID")
		return
	}

	var req model.UpdateBOMViewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, "参数错误: "+err.Error())
		return
	}

	bom, err := h.service.UpdateBOMView(c.Request.Context(), uint(id), &req)
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}

	response.Success(c, bom)
}

// Delete 删除BOM视图
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 40001, "无效的ID")
		return
	}

	if err := h.service.DeleteBOMView(c.Request.Context(), uint(id)); err != nil {
		response.Error(c, 50001, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetTree 获取BOM树形结构
func (h *Handler) GetTree(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 40001, "无效的ID")
		return
	}

	tree, err := h.service.GetBOMTree(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}

	response.Success(c, tree)
}

// AddItem 添加BOM项
func (h *Handler) AddItem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 40001, "无效的ID")
		return
	}

	var req model.AddBOMItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, "参数错误: "+err.Error())
		return
	}

	item, err := h.service.AddBOMItem(c.Request.Context(), uint(id), &req)
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}

	response.Success(c, item)
}

// UpdateItem 更新BOM项
func (h *Handler) UpdateItem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 40001, "无效的BOM ID")
		return
	}

	itemID, err := strconv.ParseUint(c.Param("itemId"), 10, 64)
	if err != nil {
		response.Error(c, 40001, "无效的BOM项ID")
		return
	}

	var req model.UpdateBOMItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, "参数错误: "+err.Error())
		return
	}

	item, err := h.service.UpdateBOMItem(c.Request.Context(), uint(id), uint(itemID), &req)
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}

	response.Success(c, item)
}

// DeleteItem 删除BOM项
func (h *Handler) DeleteItem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 40001, "无效的BOM ID")
		return
	}

	itemID, err := strconv.ParseUint(c.Param("itemId"), 10, 64)
	if err != nil {
		response.Error(c, 40001, "无效的BOM项ID")
		return
	}

	if err := h.service.DeleteBOMItem(c.Request.Context(), uint(id), uint(itemID)); err != nil {
		response.Error(c, 50001, err.Error())
		return
	}

	response.Success(c, nil)
}

// Convert 转换BOM类型
func (h *Handler) Convert(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 40001, "无效的ID")
		return
	}

	var req model.ConvertBOMRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, "参数错误: "+err.Error())
		return
	}

	if err := h.service.ConvertBOMType(c.Request.Context(), uint(id), req.ToExact); err != nil {
		response.Error(c, 50001, err.Error())
		return
	}

	response.Success(c, nil)
}

// Export 导出BOM
func (h *Handler) Export(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 40001, "无效的ID")
		return
	}

	buf, filename, err := h.service.ExportBOM(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}

	// 设置响应头
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	// 返回文件内容
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
}

// Import 导入BOM
func (h *Handler) Import(c *gin.Context) {
	// 获取BOM视图ID
	bomViewIDStr := c.PostForm("bom_view_id")
	if bomViewIDStr == "" {
		response.Error(c, 40001, "缺少bom_view_id参数")
		return
	}
	bomViewID, err := strconv.ParseUint(bomViewIDStr, 10, 64)
	if err != nil {
		response.Error(c, 40001, "无效的BOM视图ID")
		return
	}

	// 获取上传的文件
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		response.Error(c, 40001, "请上传Excel文件")
		return
	}
	defer file.Close()

	// 读取文件内容
	fileData := make([]byte, 0)
	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if err != nil {
			break
		}
		fileData = append(fileData, buf[:n]...)
	}

	userID := c.GetUint("user_id")
	result, err := h.service.ImportBOM(c.Request.Context(), uint(bomViewID), fileData, userID)
	if err != nil {
		response.Error(c, 50001, err.Error())
		return
	}

	response.Success(c, result)
}

// DownloadTemplate 下载导入模板
func (h *Handler) DownloadTemplate(c *gin.Context) {
	buf, err := h.service.GenerateImportTemplate()
	if err != nil {
		response.Error(c, 50001, "生成模板失败")
		return
	}

	// 设置响应头
	filename := "BOM_导入模板.xlsx"
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	// 返回文件内容
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
}
