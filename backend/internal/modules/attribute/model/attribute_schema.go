package model

import (
	"time"

	"gorm.io/gorm"
)

// AttributeSchema 属性Schema模型
type AttributeSchema struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	EntityType   string         `gorm:"size:50;not null;index;uniqueIndex:uk_entity_type_version" json:"entity_type"`   // 实体类型: MATERIAL/DOCUMENT
	TypeCode     string         `gorm:"size:50;not null;index;uniqueIndex:uk_entity_type_version" json:"type_code"`     // 类型编码(物料类型或文档类型)
	AttrType     string         `gorm:"size:50;not null;index;uniqueIndex:uk_entity_type_version" json:"attr_type"`     // 属性类型: main/description/specification/custom
	SchemaName   string         `gorm:"size:100" json:"schema_name"`                                                     // Schema名称
	SchemaConfig string         `gorm:"type:json;not null" json:"schema_config"`                                          // Schema定义JSON
	Version      string         `gorm:"size:2;not null;default:AA;uniqueIndex:uk_entity_type_version" json:"version"`    // 版本号
	Status       string         `gorm:"size:20;default:DRAFT" json:"status"`                                             // 状态: DRAFT/RELEASED
	IsActive     bool           `gorm:"default:true" json:"is_active"`                                                    // 是否当前激活版本
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	CreatedBy    uint           `json:"created_by"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (AttributeSchema) TableName() string {
	return "attribute_schemas"
}

// FieldSchema 字段Schema定义
type FieldSchema struct {
	Key      string   `json:"key"`                // 字段键名
	Label    string   `json:"label"`              // 字段标签
	Type     string   `json:"type"`               // 字段类型: text/number/select/date/textarea
	Unit     string   `json:"unit,omitempty"`     // 单位
	Required bool     `json:"required,omitempty"` // 是否必填
	Options  []string `json:"options,omitempty"`  // 选项(select类型)
}

// SchemaConfig Schema配置
type SchemaConfig struct {
	Label  string        `json:"label,omitempty"`  // 分组标签
	Fields []FieldSchema `json:"fields,omitempty"` // 字段列表
}

// CreateAttributeSchemaRequest 创建属性Schema请求
type CreateAttributeSchemaRequest struct {
	EntityType   string       `json:"entity_type" binding:"required"`   // MATERIAL/DOCUMENT
	TypeCode     string       `json:"type_code" binding:"required"`     // 物料类型或文档类型
	AttrType     string       `json:"attr_type" binding:"required"`     // main/description/specification/custom
	SchemaName   string       `json:"schema_name"`                      // Schema名称
	SchemaConfig SchemaConfig `json:"schema_config" binding:"required"` // Schema配置
}

// UpdateAttributeSchemaRequest 更新属性Schema请求
type UpdateAttributeSchemaRequest struct {
	SchemaName   string       `json:"schema_name"`
	SchemaConfig SchemaConfig `json:"schema_config"`
}

// AttributeSchemaQuery 属性Schema查询参数
type AttributeSchemaQuery struct {
	Page       int    `form:"page"`
	PageSize   int    `form:"page_size"`
	EntityType string `form:"entity_type"` // MATERIAL/DOCUMENT
	TypeCode   string `form:"type_code"`   // 物料类型或文档类型
	AttrType   string `form:"attr_type"`   // main/description/specification/custom
	Status     string `form:"status"`
}

// AttributeSchemaResponse 属性Schema响应
type AttributeSchemaResponse struct {
	*AttributeSchema
	SchemaConfig SchemaConfig `json:"schema_config"`
}
