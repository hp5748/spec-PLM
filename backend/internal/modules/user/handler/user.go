package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"plm/internal/modules/user/model"
	"plm/internal/modules/user/service"
	"plm/pkg/response"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

// Login 登录
func (h *Handler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	token, user, err := h.svc.Login(c, &req)
	if err != nil {
		response.Error(c, 40100, err.Error())
		return
	}

	response.Success(c, gin.H{
		"token": token,
		"user":  user,
	})
}

// Register 注册
func (h *Handler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	user, err := h.svc.Register(c, &req)
	if err != nil {
		response.Error(c, 40000, err.Error())
		return
	}

	response.SuccessWithMessage(c, "注册成功", user)
}

// RefreshToken 刷新Token
func (h *Handler) RefreshToken(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	newToken, err := h.svc.RefreshToken(c, req.Token)
	if err != nil {
		response.Unauthorized(c, "Token刷新失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"token": newToken,
	})
}

// List 用户列表
func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	users, total, err := h.svc.ListUsers(c, page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessPage(c, users, total, page, pageSize)
}

// Get 获取用户详情
func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	user, err := h.svc.GetUser(c, uint(id))
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	response.Success(c, user)
}

// Create 创建用户
func (h *Handler) Create(c *gin.Context) {
	var req model.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	user, err := h.svc.CreateUser(c, &req)
	if err != nil {
		response.Error(c, 40000, err.Error())
		return
	}

	response.SuccessWithMessage(c, "创建成功", user)
}

// Update 更新用户
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	var req model.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	user, err := h.svc.UpdateUser(c, uint(id), &req)
	if err != nil {
		response.Error(c, 40000, err.Error())
		return
	}

	response.SuccessWithMessage(c, "更新成功", user)
}

// Delete 删除用户
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	if err := h.svc.DeleteUser(c, uint(id)); err != nil {
		response.Error(c, 40000, err.Error())
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

// Organization handlers

// ListOrganizations 组织列表
func (h *Handler) ListOrganizations(c *gin.Context) {
	orgs, err := h.svc.ListOrganizations(c)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, orgs)
}

// CreateOrganization 创建组织
func (h *Handler) CreateOrganization(c *gin.Context) {
	var req model.CreateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	org, err := h.svc.CreateOrganization(c, &req)
	if err != nil {
		response.Error(c, 40000, err.Error())
		return
	}

	response.SuccessWithMessage(c, "创建成功", org)
}

// GetOrganization 获取组织详情
func (h *Handler) GetOrganization(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的组织ID")
		return
	}

	org, err := h.svc.GetOrganization(c, uint(id))
	if err != nil {
		response.NotFound(c, "组织不存在")
		return
	}

	response.Success(c, org)
}

// UpdateOrganization 更新组织
func (h *Handler) UpdateOrganization(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的组织ID")
		return
	}

	var req model.UpdateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	org, err := h.svc.UpdateOrganization(c, uint(id), &req)
	if err != nil {
		response.Error(c, 40000, err.Error())
		return
	}

	response.SuccessWithMessage(c, "更新成功", org)
}

// DeleteOrganization 删除组织
func (h *Handler) DeleteOrganization(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的组织ID")
		return
	}

	if err := h.svc.DeleteOrganization(c, uint(id)); err != nil {
		response.Error(c, 40000, err.Error())
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

// Department handlers

// ListDepartments 部门列表（树形）
func (h *Handler) ListDepartments(c *gin.Context) {
	orgID := c.Query("organization_id")
	depts, err := h.svc.ListDepartments(c, orgID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, depts)
}

// CreateDepartment 创建部门
func (h *Handler) CreateDepartment(c *gin.Context) {
	var req model.CreateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	dept, err := h.svc.CreateDepartment(c, &req)
	if err != nil {
		response.Error(c, 40000, err.Error())
		return
	}

	response.SuccessWithMessage(c, "创建成功", dept)
}

// GetDepartment 获取部门详情
func (h *Handler) GetDepartment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的部门ID")
		return
	}

	dept, err := h.svc.GetDepartment(c, uint(id))
	if err != nil {
		response.NotFound(c, "部门不存在")
		return
	}

	response.Success(c, dept)
}

// UpdateDepartment 更新部门
func (h *Handler) UpdateDepartment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的部门ID")
		return
	}

	var req model.UpdateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	dept, err := h.svc.UpdateDepartment(c, uint(id), &req)
	if err != nil {
		response.Error(c, 40000, err.Error())
		return
	}

	response.SuccessWithMessage(c, "更新成功", dept)
}

// DeleteDepartment 删除部门
func (h *Handler) DeleteDepartment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的部门ID")
		return
	}

	if err := h.svc.DeleteDepartment(c, uint(id)); err != nil {
		response.Error(c, 40000, err.Error())
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

// Role handlers

// ListRoles 角色列表
func (h *Handler) ListRoles(c *gin.Context) {
	roles, err := h.svc.ListRoles(c)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, roles)
}

// CreateRole 创建角色
func (h *Handler) CreateRole(c *gin.Context) {
	var req model.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	role, err := h.svc.CreateRole(c, &req)
	if err != nil {
		response.Error(c, 40000, err.Error())
		return
	}

	response.SuccessWithMessage(c, "创建成功", role)
}

// GetRole 获取角色详情
func (h *Handler) GetRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的角色ID")
		return
	}

	role, err := h.svc.GetRole(c, uint(id))
	if err != nil {
		response.NotFound(c, "角色不存在")
		return
	}

	response.Success(c, role)
}

// UpdateRole 更新角色
func (h *Handler) UpdateRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的角色ID")
		return
	}

	var req model.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	role, err := h.svc.UpdateRole(c, uint(id), &req)
	if err != nil {
		response.Error(c, 40000, err.Error())
		return
	}

	response.SuccessWithMessage(c, "更新成功", role)
}

// DeleteRole 删除角色
func (h *Handler) DeleteRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的角色ID")
		return
	}

	if err := h.svc.DeleteRole(c, uint(id)); err != nil {
		response.Error(c, 40000, err.Error())
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}
