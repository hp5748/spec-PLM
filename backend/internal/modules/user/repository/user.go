package repository

import (
	"context"

	"gorm.io/gorm"
	"plm/internal/modules/user/model"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

// User methods

func (r *Repository) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	err := getDB(ctx).Preload("Roles.Permissions").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := getDB(ctx).Preload("Roles.Permissions").Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	err := getDB(ctx).Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

func (r *Repository) CreateUser(ctx context.Context, user *model.User) error {
	return getDB(ctx).Create(user).Error
}

func (r *Repository) UpdateUser(ctx context.Context, user *model.User) error {
	return getDB(ctx).Save(user).Error
}

func (r *Repository) DeleteUser(ctx context.Context, id uint) error {
	return getDB(ctx).Delete(&model.User{}, id).Error
}

func (r *Repository) ListUsers(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	db := getDB(ctx).Model(&model.User{})

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Preload("Organization").Preload("Department").Preload("Roles").
		Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *Repository) AssignRoles(ctx context.Context, userID uint, roleIDs []uint) error {
	// 先清除旧角色
	if err := getDB(ctx).Exec("DELETE FROM user_roles WHERE user_id = ?", userID).Error; err != nil {
		return err
	}

	// 分配新角色
	if len(roleIDs) > 0 {
		roleUsers := make([]map[string]interface{}, len(roleIDs))
		for i, roleID := range roleIDs {
			roleUsers[i] = map[string]interface{}{
				"user_id": userID,
				"role_id": roleID,
			}
		}
		return getDB(ctx).Table("user_roles").Create(roleUsers).Error
	}
	return nil
}

// Organization methods

func (r *Repository) GetOrganizationByID(ctx context.Context, id uint) (*model.Organization, error) {
	var org model.Organization
	err := getDB(ctx).First(&org, id).Error
	if err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *Repository) CreateOrganization(ctx context.Context, org *model.Organization) error {
	return getDB(ctx).Create(org).Error
}

func (r *Repository) UpdateOrganization(ctx context.Context, org *model.Organization) error {
	return getDB(ctx).Save(org).Error
}

func (r *Repository) DeleteOrganization(ctx context.Context, id uint) error {
	return getDB(ctx).Delete(&model.Organization{}, id).Error
}

func (r *Repository) ListOrganizations(ctx context.Context) ([]*model.Organization, error) {
	var orgs []*model.Organization
	err := getDB(ctx).Find(&orgs).Error
	return orgs, err
}

// Department methods

func (r *Repository) GetDepartmentByID(ctx context.Context, id uint) (*model.Department, error) {
	var dept model.Department
	err := getDB(ctx).Preload("Organization").Preload("Parent").First(&dept, id).Error
	if err != nil {
		return nil, err
	}
	return &dept, nil
}

func (r *Repository) CreateDepartment(ctx context.Context, dept *model.Department) error {
	return getDB(ctx).Create(dept).Error
}

func (r *Repository) UpdateDepartment(ctx context.Context, dept *model.Department) error {
	return getDB(ctx).Save(dept).Error
}

func (r *Repository) DeleteDepartment(ctx context.Context, id uint) error {
	return getDB(ctx).Delete(&model.Department{}, id).Error
}

func (r *Repository) ListDepartments(ctx context.Context, orgID string) ([]*model.Department, error) {
	var depts []*model.Department
	db := getDB(ctx).Preload("Organization").Preload("Parent")
	if orgID != "" {
		db = db.Where("organization_id = ?", orgID)
	}
	err := db.Order("sort_order").Find(&depts).Error
	return depts, err
}

// Role methods

func (r *Repository) GetRoleByID(ctx context.Context, id uint) (*model.Role, error) {
	var role model.Role
	err := getDB(ctx).Preload("Permissions").First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *Repository) CreateRole(ctx context.Context, role *model.Role) error {
	return getDB(ctx).Create(role).Error
}

func (r *Repository) UpdateRole(ctx context.Context, role *model.Role) error {
	return getDB(ctx).Save(role).Error
}

func (r *Repository) DeleteRole(ctx context.Context, id uint) error {
	return getDB(ctx).Delete(&model.Role{}, id).Error
}

func (r *Repository) ListRoles(ctx context.Context) ([]*model.Role, error) {
	var roles []*model.Role
	err := getDB(ctx).Preload("Permissions").Find(&roles).Error
	return roles, err
}

func (r *Repository) AssignPermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	// 先清除旧权限
	if err := getDB(ctx).Exec("DELETE FROM role_permissions WHERE role_id = ?", roleID).Error; err != nil {
		return err
	}

	// 分配新权限
	if len(permissionIDs) > 0 {
		rolePerms := make([]map[string]interface{}, len(permissionIDs))
		for i, permID := range permissionIDs {
			rolePerms[i] = map[string]interface{}{
				"role_id":       roleID,
				"permission_id": permID,
			}
		}
		return getDB(ctx).Table("role_permissions").Create(rolePerms).Error
	}
	return nil
}

// getDB 获取数据库实例
func getDB(ctx context.Context) *gorm.DB {
	// TODO: 从context获取DB或使用全局DB
	return globalDB
}

var globalDB *gorm.DB

func SetDB(db *gorm.DB) {
	globalDB = db
}
