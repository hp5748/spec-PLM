package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"plm/internal/modules/user/model"
	"plm/internal/modules/user/repository"
	"plm/pkg/auth"
	"plm/pkg/utils"
)

type Service struct {
	repo *repository.Repository
	jwt  *auth.JWT
}

func NewService(repo *repository.Repository, jwtSecret string, jwtExpire time.Duration) *Service {
	return &Service{
		repo: repo,
		jwt: auth.NewJWT(&auth.JWTConfig{
			Secret:     jwtSecret,
			ExpireTime: jwtExpire,
		}),
	}
}

// Login 用户登录
func (s *Service) Login(ctx context.Context, req *model.LoginRequest) (string, *model.User, error) {
	user, err := s.repo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return "", nil, errors.New("用户名或密码错误")
	}

	// 检查用户状态
	if user.Status == "locked" {
		return "", nil, errors.New("用户已被锁定")
	}
	if user.Status == "inactive" {
		return "", nil, errors.New("用户已停用")
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		return "", nil, errors.New("用户名或密码错误")
	}

	// 获取用户的组织和部门名称
	var orgName, deptName string
	if user.OrganizationID != nil {
		org, err := s.repo.GetOrganizationByID(ctx, *user.OrganizationID)
		if err == nil {
			orgName = org.Name
		}
	}
	if user.DepartmentID != nil {
		dept, err := s.repo.GetDepartmentByID(ctx, *user.DepartmentID)
		if err == nil {
			deptName = dept.Name
		}
	}

	// 获取用户角色
	userWithRoles, _ := s.repo.GetUserByID(ctx, user.ID)
	var roleNames []string
	if userWithRoles != nil {
		for _, role := range userWithRoles.Roles {
			roleNames = append(roleNames, role.Name)
		}
	}
	rolesStr := strings.Join(roleNames, ",")

	// 生成JWT Token
	token, err := s.jwt.GenerateToken(user.ID, user.Username, user.RealName, orgName, deptName, rolesStr)
	if err != nil {
		return "", nil, errors.New("生成Token失败")
	}

	return token, user, nil
}

// RefreshToken 刷新Token
func (s *Service) RefreshToken(ctx context.Context, oldToken string) (string, error) {
	newToken, err := s.jwt.RefreshToken(oldToken)
	if err != nil {
		return "", err
	}
	return newToken, nil
}

// Register 用户注册
func (s *Service) Register(ctx context.Context, req *model.RegisterRequest) (*model.User, error) {
	// 检查用户名是否已存在
	exists, err := s.repo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("用户名已存在")
	}

	// 哈希密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
		Phone:    req.Phone,
		RealName: req.RealName,
		Status:   "active",
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// ListUsers 用户列表
func (s *Service) ListUsers(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	return s.repo.ListUsers(ctx, page, pageSize)
}

// GetUser 获取用户详情
func (s *Service) GetUser(ctx context.Context, id uint) (*model.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

// CreateUser 创建用户
func (s *Service) CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
	// 检查用户名是否已存在
	exists, err := s.repo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("用户名已存在")
	}

	// 哈希密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username:       req.Username,
		Password:       hashedPassword,
		Email:          req.Email,
		Phone:          req.Phone,
		RealName:       req.RealName,
		OrganizationID: req.OrganizationID,
		DepartmentID:   req.DepartmentID,
		Status:         "active",
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	// 分配角色
	if len(req.RoleIDs) > 0 {
		if err := s.repo.AssignRoles(ctx, user.ID, req.RoleIDs); err != nil {
			return nil, err
		}
	}

	return user, nil
}

// UpdateUser 更新用户
func (s *Service) UpdateUser(ctx context.Context, id uint, req *model.UpdateUserRequest) (*model.User, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.RealName != "" {
		user.RealName = req.RealName
	}
	if req.OrganizationID != nil {
		user.OrganizationID = req.OrganizationID
	}
	if req.DepartmentID != nil {
		user.DepartmentID = req.DepartmentID
	}
	if req.Status != "" {
		user.Status = req.Status
	}

	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	// 更新角色
	if req.RoleIDs != nil {
		if err := s.repo.AssignRoles(ctx, user.ID, req.RoleIDs); err != nil {
			return nil, err
		}
	}

	return user, nil
}

// DeleteUser 删除用户
func (s *Service) DeleteUser(ctx context.Context, id uint) error {
	return s.repo.DeleteUser(ctx, id)
}

// Organization methods

// ListOrganizations 组织列表
func (s *Service) ListOrganizations(ctx context.Context) ([]*model.Organization, error) {
	return s.repo.ListOrganizations(ctx)
}

// CreateOrganization 创建组织
func (s *Service) CreateOrganization(ctx context.Context, req *model.CreateOrganizationRequest) (*model.Organization, error) {
	org := &model.Organization{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Logo:        req.Logo,
		Status:      "active",
	}

	if err := s.repo.CreateOrganization(ctx, org); err != nil {
		return nil, err
	}

	return org, nil
}

// GetOrganization 获取组织详情
func (s *Service) GetOrganization(ctx context.Context, id uint) (*model.Organization, error) {
	return s.repo.GetOrganizationByID(ctx, id)
}

// UpdateOrganization 更新组织
func (s *Service) UpdateOrganization(ctx context.Context, id uint, req *model.UpdateOrganizationRequest) (*model.Organization, error) {
	org, err := s.repo.GetOrganizationByID(ctx, id)
	if err != nil {
		return nil, errors.New("组织不存在")
	}

	if req.Name != "" {
		org.Name = req.Name
	}
	if req.Description != "" {
		org.Description = req.Description
	}
	if req.Logo != "" {
		org.Logo = req.Logo
	}
	if req.Status != "" {
		org.Status = req.Status
	}

	if err := s.repo.UpdateOrganization(ctx, org); err != nil {
		return nil, err
	}

	return org, nil
}

// DeleteOrganization 删除组织
func (s *Service) DeleteOrganization(ctx context.Context, id uint) error {
	return s.repo.DeleteOrganization(ctx, id)
}

// Department methods

// ListDepartments 部门列表
func (s *Service) ListDepartments(ctx context.Context, orgID string) ([]*model.Department, error) {
	return s.repo.ListDepartments(ctx, orgID)
}

// CreateDepartment 创建部门
func (s *Service) CreateDepartment(ctx context.Context, req *model.CreateDepartmentRequest) (*model.Department, error) {
	dept := &model.Department{
		OrganizationID: req.OrganizationID,
		ParentID:       req.ParentID,
		Name:           req.Name,
		Code:           req.Code,
		SortOrder:      req.SortOrder,
		ManagerID:      req.ManagerID,
		Status:         "active",
	}

	if err := s.repo.CreateDepartment(ctx, dept); err != nil {
		return nil, err
	}

	return dept, nil
}

// GetDepartment 获取部门详情
func (s *Service) GetDepartment(ctx context.Context, id uint) (*model.Department, error) {
	return s.repo.GetDepartmentByID(ctx, id)
}

// UpdateDepartment 更新部门
func (s *Service) UpdateDepartment(ctx context.Context, id uint, req *model.UpdateDepartmentRequest) (*model.Department, error) {
	dept, err := s.repo.GetDepartmentByID(ctx, id)
	if err != nil {
		return nil, errors.New("部门不存在")
	}

	if req.Name != "" {
		dept.Name = req.Name
	}
	if req.Code != "" {
		dept.Code = req.Code
	}
	if req.SortOrder != 0 {
		dept.SortOrder = req.SortOrder
	}
	if req.ManagerID != nil {
		dept.ManagerID = req.ManagerID
	}
	if req.Status != "" {
		dept.Status = req.Status
	}

	if err := s.repo.UpdateDepartment(ctx, dept); err != nil {
		return nil, err
	}

	return dept, nil
}

// DeleteDepartment 删除部门
func (s *Service) DeleteDepartment(ctx context.Context, id uint) error {
	return s.repo.DeleteDepartment(ctx, id)
}

// Role methods

// ListRoles 角色列表
func (s *Service) ListRoles(ctx context.Context) ([]*model.Role, error) {
	return s.repo.ListRoles(ctx)
}

// CreateRole 创建角色
func (s *Service) CreateRole(ctx context.Context, req *model.CreateRoleRequest) (*model.Role, error) {
	role := &model.Role{
		OrganizationID: req.OrganizationID,
		Name:           req.Name,
		Code:           req.Code,
		Description:    req.Description,
		IsSystem:       false,
	}

	if err := s.repo.CreateRole(ctx, role); err != nil {
		return nil, err
	}

	// 分配权限
	if len(req.PermissionIDs) > 0 {
		if err := s.repo.AssignPermissions(ctx, role.ID, req.PermissionIDs); err != nil {
			return nil, err
		}
	}

	return role, nil
}

// GetRole 获取角色详情
func (s *Service) GetRole(ctx context.Context, id uint) (*model.Role, error) {
	return s.repo.GetRoleByID(ctx, id)
}

// UpdateRole 更新角色
func (s *Service) UpdateRole(ctx context.Context, id uint, req *model.UpdateRoleRequest) (*model.Role, error) {
	role, err := s.repo.GetRoleByID(ctx, id)
	if err != nil {
		return nil, errors.New("角色不存在")
	}

	if req.Name != "" {
		role.Name = req.Name
	}
	if req.Description != "" {
		role.Description = req.Description
	}

	if err := s.repo.UpdateRole(ctx, role); err != nil {
		return nil, err
	}

	// 更新权限
	if req.PermissionIDs != nil {
		if err := s.repo.AssignPermissions(ctx, role.ID, req.PermissionIDs); err != nil {
			return nil, err
		}
	}

	return role, nil
}

// DeleteRole 删除角色
func (s *Service) DeleteRole(ctx context.Context, id uint) error {
	role, err := s.repo.GetRoleByID(ctx, id)
	if err != nil {
		return err
	}

	if role.IsSystem {
		return errors.New("系统内置角色不能删除")
	}

	return s.repo.DeleteRole(ctx, id)
}
