package model

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	RealName string `json:"real_name"`
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username       string `json:"username" binding:"required"`
	Password       string `json:"password" binding:"required,min=8"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	RealName       string `json:"real_name"`
	OrganizationID *uint  `json:"organization_id"`
	DepartmentID   *uint  `json:"department_id"`
	RoleIDs        []uint `json:"role_ids"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	RealName       string `json:"real_name"`
	OrganizationID *uint  `json:"organization_id"`
	DepartmentID   *uint  `json:"department_id"`
	RoleIDs        []uint `json:"role_ids"`
	Status         string `json:"status"`
}

// CreateOrganizationRequest 创建组织请求
type CreateOrganizationRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Description string `json:"description"`
	Logo        string `json:"logo"`
}

// UpdateOrganizationRequest 更新组织请求
type UpdateOrganizationRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Logo        string `json:"logo"`
	Status      string `json:"status"`
}

// CreateDepartmentRequest 创建部门请求
type CreateDepartmentRequest struct {
	OrganizationID uint   `json:"organization_id" binding:"required"`
	ParentID       *uint  `json:"parent_id"`
	Name           string `json:"name" binding:"required"`
	Code           string `json:"code"`
	SortOrder      int    `json:"sort_order"`
	ManagerID      *uint  `json:"manager_id"`
}

// UpdateDepartmentRequest 更新部门请求
type UpdateDepartmentRequest struct {
	Name      string `json:"name"`
	Code      string `json:"code"`
	SortOrder int    `json:"sort_order"`
	ManagerID *uint  `json:"manager_id"`
	Status    string `json:"status"`
}

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	OrganizationID *uint  `json:"organization_id"`
	Name           string `json:"name" binding:"required"`
	Code           string `json:"code" binding:"required"`
	Description    string `json:"description"`
	PermissionIDs  []uint `json:"permission_ids"`
}

// UpdateRoleRequest 更新角色请求
type UpdateRoleRequest struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	PermissionIDs []uint `json:"permission_ids"`
}
