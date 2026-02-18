package repository

import (
	"encoding/json"

	"plm/internal/modules/attribute/model"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) SetDB(db *gorm.DB) {
	r.db = db
}

// Create 创建属性Schema
func (r *Repository) Create(schema *model.AttributeSchema) error {
	return r.db.Create(schema).Error
}

// Update 更新属性Schema
func (r *Repository) Update(schema *model.AttributeSchema) error {
	return r.db.Save(schema).Error
}

// Delete 删除属性Schema（软删除）
func (r *Repository) Delete(id uint) error {
	return r.db.Delete(&model.AttributeSchema{}, id).Error
}

// GetByID 根据ID获取属性Schema
func (r *Repository) GetByID(id uint) (*model.AttributeSchema, error) {
	var schema model.AttributeSchema
	err := r.db.First(&schema, id).Error
	if err != nil {
		return nil, err
	}
	return &schema, nil
}

// GetActiveSchema 获取当前激活的Schema
func (r *Repository) GetActiveSchema(entityType, typeCode, attrType string) (*model.AttributeSchema, error) {
	var schema model.AttributeSchema
	err := r.db.Where("entity_type = ? AND type_code = ? AND attr_type = ? AND is_active = ?", entityType, typeCode, attrType, true).
		First(&schema).Error
	if err != nil {
		return nil, err
	}
	return &schema, nil
}

// GetAllActiveSchemasByEntity 获取指定实体类型的所有激活Schema
func (r *Repository) GetAllActiveSchemasByEntity(entityType, typeCode string) ([]model.AttributeSchema, error) {
	var schemas []model.AttributeSchema
	err := r.db.Where("entity_type = ? AND type_code = ? AND is_active = ?", entityType, typeCode, true).
		Order("attr_type").
		Find(&schemas).Error
	return schemas, err
}

// List 分页查询属性Schema列表
func (r *Repository) List(query *model.AttributeSchemaQuery) ([]model.AttributeSchema, int64, error) {
	var schemas []model.AttributeSchema
	var total int64

	db := r.db.Model(&model.AttributeSchema{})

	if query.EntityType != "" {
		db = db.Where("entity_type = ?", query.EntityType)
	}
	if query.TypeCode != "" {
		db = db.Where("type_code = ?", query.TypeCode)
	}
	if query.AttrType != "" {
		db = db.Where("attr_type = ?", query.AttrType)
	}
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	// 统计总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	page := query.Page
	pageSize := query.PageSize
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	err := db.Order("entity_type, type_code, attr_type, version DESC").
		Offset(offset).Limit(pageSize).
		Find(&schemas).Error

	return schemas, total, err
}

// DeactivateSchemas 将指定类型的所有Schema设为非激活
func (r *Repository) DeactivateSchemas(entityType, typeCode, attrType string) error {
	return r.db.Model(&model.AttributeSchema{}).
		Where("entity_type = ? AND type_code = ? AND attr_type = ?", entityType, typeCode, attrType).
		Update("is_active", false).Error
}

// GetLatestVersion 获取最新版本号
func (r *Repository) GetLatestVersion(entityType, typeCode, attrType string) (string, error) {
	var schema model.AttributeSchema
	err := r.db.Where("entity_type = ? AND type_code = ? AND attr_type = ?", entityType, typeCode, attrType).
		Order("version DESC").
		First(&schema).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil
		}
		return "", err
	}
	return schema.Version, nil
}

// ParseSchemaConfig 解析Schema配置
func ParseSchemaConfig(configStr string) (model.SchemaConfig, error) {
	var config model.SchemaConfig
	if configStr == "" {
		return config, nil
	}
	err := json.Unmarshal([]byte(configStr), &config)
	return config, err
}

// SchemaConfigToJSON 将Schema配置转为JSON字符串
func SchemaConfigToJSON(config model.SchemaConfig) (string, error) {
	bytes, err := json.Marshal(config)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
