package model

import (
	"time"

	"gorm.io/gorm"
)

// BOMView BOM视图模型
type BOMView struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Name           string         `gorm:"size:200;not null" json:"name"`
	RootMaterialID uint           `gorm:"not null;uniqueIndex:uk_root_material" json:"root_material_id"`
	RootVersion    string         `gorm:"size:2" json:"root_version"`
	IsExact        bool           `gorm:"not null;default:false" json:"is_exact"`
	Status         string         `gorm:"size:20;default:DRAFT;index" json:"status"`
	Description    string         `json:"description"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	CreatedBy      uint           `gorm:"not null" json:"created_by"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	RootMaterial *Material    `gorm:"foreignKey:RootMaterialID" json:"root_material,omitempty"`
	Items        []BOMItem    `gorm:"foreignKey:BOMViewID" json:"items,omitempty"`
}

// BOMItem BOM项模型
type BOMItem struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	BOMViewID   uint      `gorm:"not null;index" json:"bom_view_id"`
	ParentID    *uint     `gorm:"index" json:"parent_id"`
	MaterialID  uint      `gorm:"not null;index" json:"material_id"`
	Version     string    `gorm:"size:2" json:"version"`
	Quantity    float64   `gorm:"type:decimal(10,4);not null;default:1" json:"quantity"`
	SortOrder   int       `gorm:"default:0" json:"sort_order"`
	Level       int       `gorm:"default:1" json:"level"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// 关联
	Material *Material  `gorm:"foreignKey:MaterialID" json:"material,omitempty"`
	Children []BOMItem  `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

// Material 简化的物料模型（用于关联查询）
type Material struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	ItemID   string `gorm:"size:50;not null" json:"item_id"`
	ItemName string `gorm:"size:200;not null" json:"item_name"`
	ItemType string `gorm:"size:50" json:"item_type"`
	Version  string `gorm:"size:2" json:"version"`
	Unit     string `gorm:"size:20" json:"unit"`
}

// BOMTreeNode BOM树节点（用于前端展示）
type BOMTreeNode struct {
	ID         uint          `json:"id"`
	MaterialID uint          `json:"material_id"`
	ItemID     string        `json:"item_id"`
	ItemName   string        `json:"item_name"`
	Version    string        `json:"version"`
	Quantity   float64       `json:"quantity"`
	SortOrder  int           `json:"sort_order"`
	Level      int           `json:"level"`
	Unit       string        `json:"unit"`
	Children   []BOMTreeNode `json:"children,omitempty"`
}

// CreateBOMViewRequest 创建BOM视图请求
type CreateBOMViewRequest struct {
	Name           string  `json:"name" binding:"required,max=200"`
	RootMaterialID uint    `json:"root_material_id" binding:"required"`
	RootVersion    string  `json:"root_version"`
	IsExact        bool    `json:"is_exact"`
	Description    string  `json:"description"`
}

// UpdateBOMViewRequest 更新BOM视图请求
type UpdateBOMViewRequest struct {
	Name        string `json:"name" binding:"max=200"`
	Status      string `json:"status" binding:"omitempty,oneof=DRAFT REVIEWING RELEASED REJECTED OBSOLETE"`
	Description string `json:"description"`
}

// AddBOMItemRequest 添加BOM项请求
type AddBOMItemRequest struct {
	ParentID   *uint   `json:"parent_id"`
	MaterialID uint    `json:"material_id" binding:"required"`
	Version    string  `json:"version"`
	Quantity   float64 `json:"quantity" binding:"required,gt=0"`
	SortOrder  int     `json:"sort_order"`
}

// UpdateBOMItemRequest 更新BOM项请求
type UpdateBOMItemRequest struct {
	Quantity  float64 `json:"quantity" binding:"required,gt=0"`
	SortOrder int     `json:"sort_order"`
}

// ConvertBOMRequest 转换BOM类型请求
type ConvertBOMRequest struct {
	ToExact bool `json:"to_exact"`
}

// BOMListQuery BOM列表查询参数
type BOMListQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Name     string `form:"name"`
	Status   string `form:"status"`
	IsExact  *bool  `form:"is_exact"`
}

// BOMViewResponse BOM视图响应
type BOMViewResponse struct {
	BOMView
	ItemCount int `json:"item_count"`
}

// ImportBOMResult 导入BOM结果
type ImportBOMResult struct {
	SuccessCount int      `json:"success_count"`
	FailCount    int      `json:"fail_count"`
	Errors       []string `json:"errors,omitempty"`
}

// TableName 指定表名
func (BOMView) TableName() string {
	return "bom_views"
}

func (BOMItem) TableName() string {
	return "bom_items"
}

func (Material) TableName() string {
	return "materials"
}
