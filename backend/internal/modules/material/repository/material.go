package repository

import (
	"context"

	"gorm.io/gorm"
	"plm/internal/modules/material/model"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

// GetMaterialByID 根据ID获取物料
func (r *Repository) GetMaterialByID(ctx context.Context, id uint) (*model.Material, error) {
	var material model.Material
	err := getDB(ctx).First(&material, id).Error
	if err != nil {
		return nil, err
	}
	return &material, nil
}

// GetMaterialByItemID 根据物料编码获取物料
func (r *Repository) GetMaterialByItemID(ctx context.Context, itemID string, version string) (*model.Material, error) {
	var material model.Material
	err := getDB(ctx).Where("item_id = ? AND version = ?", itemID, version).First(&material).Error
	if err != nil {
		return nil, err
	}
	return &material, nil
}

// ExistsByItemID 检查物料编码是否存在
func (r *Repository) ExistsByItemID(ctx context.Context, itemID string) (bool, error) {
	var count int64
	err := getDB(ctx).Model(&model.Material{}).Where("item_id = ?", itemID).Count(&count).Error
	return count > 0, err
}

// CreateMaterial 创建物料
func (r *Repository) CreateMaterial(ctx context.Context, material *model.Material) error {
	return getDB(ctx).Create(material).Error
}

// UpdateMaterial 更新物料
func (r *Repository) UpdateMaterial(ctx context.Context, material *model.Material) error {
	return getDB(ctx).Save(material).Error
}

// DeleteMaterial 删除物料
func (r *Repository) DeleteMaterial(ctx context.Context, id uint) error {
	return getDB(ctx).Delete(&model.Material{}, id).Error
}

// ListMaterials 物料列表（分页、筛选）
func (r *Repository) ListMaterials(ctx context.Context, query *model.MaterialListQuery) ([]*model.Material, int64, error) {
	var materials []*model.Material
	var total int64

	db := getDB(ctx).Model(&model.Material{})

	// 筛选条件
	if query.ItemID != "" {
		db = db.Where("item_id LIKE ?", "%"+query.ItemID+"%")
	}
	if query.ItemName != "" {
		db = db.Where("item_name LIKE ?", "%"+query.ItemName+"%")
	}
	if query.ItemType != "" {
		db = db.Where("item_type = ?", query.ItemType)
	}
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}
	if query.Keyword != "" {
		db = db.Where("item_id LIKE ? OR item_name LIKE ?", "%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}

	// 计数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	order := "created_at DESC"
	if query.Sort != "" {
		order = query.Sort
		if query.Order != "" {
			order = query.Sort + " " + query.Order
		}
	}

	// 分页
	offset := (query.Page - 1) * query.PageSize
	if err := db.Order(order).Offset(offset).Limit(query.PageSize).Find(&materials).Error; err != nil {
		return nil, 0, err
	}

	return materials, total, nil
}

// SearchMaterials 搜索物料
func (r *Repository) SearchMaterials(ctx context.Context, keyword string, page, pageSize int) ([]*model.Material, int64, error) {
	var materials []*model.Material
	var total int64

	db := getDB(ctx).Model(&model.Material{}).Where(
		"item_id LIKE ? OR item_name LIKE ? OR description LIKE ?",
		"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%",
	)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&materials).Error; err != nil {
		return nil, 0, err
	}

	return materials, total, nil
}

// MaterialAttribute methods

// GetAttributesByMaterialID 获取物料属性
func (r *Repository) GetAttributesByMaterialID(ctx context.Context, materialID uint) ([]model.MaterialAttribute, error) {
	var attrs []model.MaterialAttribute
	err := getDB(ctx).Where("material_id = ?", materialID).Order("attr_type, sort_order").Find(&attrs).Error
	return attrs, err
}

// CreateAttribute 创建物料属性
func (r *Repository) CreateAttribute(ctx context.Context, attr *model.MaterialAttribute) error {
	return getDB(ctx).Create(attr).Error
}

// CreateAttributes 批量创建物料属性
func (r *Repository) CreateAttributes(ctx context.Context, attrs []model.MaterialAttribute) error {
	if len(attrs) == 0 {
		return nil
	}
	return getDB(ctx).Create(&attrs).Error
}

// DeleteAttributesByMaterialID 删除物料所有属性
func (r *Repository) DeleteAttributesByMaterialID(ctx context.Context, materialID uint) error {
	return getDB(ctx).Where("material_id = ?", materialID).Delete(&model.MaterialAttribute{}).Error
}

// getDB 获取数据库实例
func getDB(ctx context.Context) *gorm.DB {
	return globalDB
}

var globalDB *gorm.DB

func SetDB(db *gorm.DB) {
	globalDB = db
}
