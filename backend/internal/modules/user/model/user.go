package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Username       string         `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Password       string         `gorm:"size:255;not null" json:"-"`
	Email          string         `gorm:"size:100" json:"email"`
	Phone          string         `gorm:"size:20" json:"phone"`
	RealName       string         `gorm:"size:50" json:"real_name"`
	Avatar         string         `gorm:"size:255" json:"avatar"`
	OrganizationID *uint          `json:"organization_id"`
	Organization   *Organization  `json:"organization,omitempty"`
	DepartmentID   *uint          `json:"department_id"`
	Department     *Department    `json:"department,omitempty"`
	Status         string         `gorm:"size:20;default:active" json:"status"` // active, inactive, locked
	LastLoginAt    *time.Time     `json:"last_login_at"`
	LoginFailCount int            `gorm:"default:0" json:"login_fail_count"`
	LockedUntil    *time.Time     `json:"locked_until"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	CreatedBy      uint           `json:"created_by"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	Roles          []Role         `gorm:"many2many:user_roles;" json:"roles,omitempty"`
}

func (User) TableName() string {
	return "users"
}

// Organization 组织模型
type Organization struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Code        string         `gorm:"size:50;uniqueIndex;not null" json:"code"`
	Description string         `json:"description"`
	Logo        string         `gorm:"size:255" json:"logo"`
	Status      string         `gorm:"size:20;default:active" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	CreatedBy   uint           `json:"created_by"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Organization) TableName() string {
	return "organizations"
}

// Department 部门模型
type Department struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	OrganizationID uint           `gorm:"not null;index" json:"organization_id"`
	Organization   *Organization  `json:"organization,omitempty"`
	ParentID       *uint          `gorm:"index" json:"parent_id"`
	Parent         *Department    `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children       []Department   `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Name           string         `gorm:"size:100;not null" json:"name"`
	Code           string         `gorm:"size:50" json:"code"`
	SortOrder      int            `gorm:"default:0" json:"sort_order"`
	ManagerID      *uint          `json:"manager_id"`
	Status         string         `gorm:"size:20;default:active" json:"status"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	CreatedBy      uint           `json:"created_by"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Department) TableName() string {
	return "departments"
}

// Role 角色模型
type Role struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	OrganizationID *uint          `gorm:"index" json:"organization_id"`
	Name           string         `gorm:"size:50;not null" json:"name"`
	Code           string         `gorm:"size:50;not null" json:"code"`
	Description    string         `json:"description"`
	IsSystem       bool           `gorm:"default:false" json:"is_system"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	CreatedBy      uint           `json:"created_by"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	Permissions    []Permission   `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
}

func (Role) TableName() string {
	return "roles"
}

// Permission 权限模型
type Permission struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:50;not null" json:"name"`
	Code        string    `gorm:"size:100;uniqueIndex;not null" json:"code"`
	Module      string    `gorm:"size:50;not null;index" json:"module"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Permission) TableName() string {
	return "permissions"
}
