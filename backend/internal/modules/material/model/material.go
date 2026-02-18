package model

import (
	"time"

	"gorm.io/gorm"
)

// Material 物料模型
type Material struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	ItemID      string         `gorm:"size:50;not null;uniqueIndex:uk_item_version" json:"item_id"`       // 物料编码
	ItemName    string         `gorm:"size:200;not null" json:"item_name"`                                 // 物料名称
	Description string         `json:"description"`                                                        // 描述
	ItemType    string         `gorm:"size:50;not null;index" json:"item_type"`                            // 物料类型: PART/ASSEMBLY/RAW_MATERIAL/TOOL
	Version     string         `gorm:"size:2;not null;default:AA;uniqueIndex:uk_item_version" json:"version"` // 版本号
	Unit        string         `gorm:"size:20" json:"unit"`                                                // 计量单位
	Status      string         `gorm:"size:20;default:DRAFT;index" json:"status"`                          // 状态: DRAFT/REVIEWING/RELEASED/REJECTED/OBSOLETE
	Attributes  string         `gorm:"type:json" json:"attributes"`                                        // 主属性JSON
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	CreatedBy   uint           `json:"created_by"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Material) TableName() string {
	return "materials"
}

// MaterialAttribute 物料动态属性模型
type MaterialAttribute struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	MaterialID  uint      `gorm:"not null;index:idx_material_type" json:"material_id"` // 物料ID
	AttrType    string    `gorm:"size:50;not null;index:idx_material_type" json:"attr_type"` // 属性类型
	AttrKey     string    `gorm:"size:100;not null" json:"attr_key"`                          // 属性键
	AttrValue   string    `json:"attr_value"`                                                 // 属性值
	Unit        string    `gorm:"size:20" json:"unit"`                                        // 单位
	SortOrder   int       `gorm:"default:0" json:"sort_order"`                                // 排序
	CreatedAt   time.Time `json:"created_at"`
}

func (MaterialAttribute) TableName() string {
	return "material_attributes"
}

// CreateMaterialRequest 创建物料请求
type CreateMaterialRequest struct {
	ItemID      string `json:"item_id" binding:"required"`
	ItemName    string `json:"item_name" binding:"required"`
	Description string `json:"description"`
	ItemType    string `json:"item_type" binding:"required"`
	Unit        string `json:"unit"`
	Attributes  string `json:"attributes"`
}

// UpdateMaterialRequest 更新物料请求
type UpdateMaterialRequest struct {
	ItemName    string `json:"item_name"`
	Description string `json:"description"`
	ItemType    string `json:"item_type"`
	Unit        string `json:"unit"`
	Status      string `json:"status"`
	Attributes  string `json:"attributes"`
}

// MaterialListQuery 物料列表查询参数
type MaterialListQuery struct {
	Page      int    `form:"page"`
	PageSize  int    `form:"page_size"`
	Sort      string `form:"sort"`
	Order     string `form:"order"`
	ItemID    string `form:"item_id"`
	ItemName  string `form:"item_name"`
	ItemType  string `form:"item_type"`
	Status    string `form:"status"`
	Keyword   string `form:"keyword"` // 关键字搜索（编码或名称）
}

// MaterialResponse 物料响应
type MaterialResponse struct {
	*Material
	Attributes []MaterialAttribute `json:"attributes,omitempty"`
}
