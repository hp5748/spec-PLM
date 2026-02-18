package repository

import (
	"context"

	"plm/internal/modules/bom/model"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// CreateBOMView 创建BOM视图
func (r *Repository) CreateBOMView(ctx context.Context, bom *model.BOMView) error {
	return r.db.WithContext(ctx).Create(bom).Error
}

// GetBOMViewByID 根据ID获取BOM视图
func (r *Repository) GetBOMViewByID(ctx context.Context, id uint) (*model.BOMView, error) {
	var bom model.BOMView
	err := r.db.WithContext(ctx).
		Preload("RootMaterial").
		First(&bom, id).Error
	if err != nil {
		return nil, err
	}
	return &bom, nil
}

// GetBOMViewByRootMaterialID 根据根物料ID获取BOM视图
func (r *Repository) GetBOMViewByRootMaterialID(ctx context.Context, rootMaterialID uint) (*model.BOMView, error) {
	var bom model.BOMView
	err := r.db.WithContext(ctx).
		Where("root_material_id = ?", rootMaterialID).
		First(&bom).Error
	if err != nil {
		return nil, err
	}
	return &bom, nil
}

// UpdateBOMView 更新BOM视图
func (r *Repository) UpdateBOMView(ctx context.Context, bom *model.BOMView) error {
	return r.db.WithContext(ctx).Save(bom).Error
}

// DeleteBOMView 删除BOM视图
func (r *Repository) DeleteBOMView(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.BOMView{}, id).Error
}

// ListBOMViews BOM视图列表
func (r *Repository) ListBOMViews(ctx context.Context, query *model.BOMListQuery) ([]*model.BOMView, int64, error) {
	var list []*model.BOMView
	var total int64

	db := r.db.WithContext(ctx).Model(&model.BOMView{}).Preload("RootMaterial")

	if query.Name != "" {
		db = db.Where("name LIKE ?", "%"+query.Name+"%")
	}
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}
	if query.IsExact != nil {
		db = db.Where("is_exact = ?", *query.IsExact)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (query.Page - 1) * query.PageSize
	if err := db.Order("created_at DESC").Offset(offset).Limit(query.PageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

// CreateBOMItem 创建BOM项
func (r *Repository) CreateBOMItem(ctx context.Context, item *model.BOMItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

// GetBOMItemByID 根据ID获取BOM项
func (r *Repository) GetBOMItemByID(ctx context.Context, id uint) (*model.BOMItem, error) {
	var item model.BOMItem
	err := r.db.WithContext(ctx).
		Preload("Material").
		First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// UpdateBOMItem 更新BOM项
func (r *Repository) UpdateBOMItem(ctx context.Context, item *model.BOMItem) error {
	return r.db.WithContext(ctx).Save(item).Error
}

// DeleteBOMItem 删除BOM项
func (r *Repository) DeleteBOMItem(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.BOMItem{}, id).Error
}

// GetBOMItemsByViewID 根据BOM视图ID获取所有项
func (r *Repository) GetBOMItemsByViewID(ctx context.Context, bomViewID uint) ([]model.BOMItem, error) {
	var items []model.BOMItem
	err := r.db.WithContext(ctx).
		Where("bom_view_id = ?", bomViewID).
		Preload("Material").
		Order("level ASC, sort_order ASC").
		Find(&items).Error
	return items, err
}

// GetBOMItemsByParentID 根据父ID获取子项
func (r *Repository) GetBOMItemsByParentID(ctx context.Context, bomViewID uint, parentID *uint) ([]model.BOMItem, error) {
	var items []model.BOMItem
	db := r.db.WithContext(ctx).Where("bom_view_id = ?", bomViewID)
	if parentID == nil {
		db = db.Where("parent_id IS NULL")
	} else {
		db = db.Where("parent_id = ?", *parentID)
	}
	err := db.Preload("Material").Order("sort_order ASC").Find(&items).Error
	return items, err
}

// GetChildrenByMaterialID 获取指定物料作为父节点的所有子项
func (r *Repository) GetChildrenByMaterialID(ctx context.Context, bomViewID uint, materialID uint) ([]model.BOMItem, error) {
	var items []model.BOMItem
	err := r.db.WithContext(ctx).
		Joins("JOIN bom_items AS parent ON parent.id = bom_items.parent_id").
		Where("parent.bom_view_id = ? AND parent.material_id = ?", bomViewID, materialID).
		Preload("Material").
		Find(&items).Error
	return items, err
}

// DeleteBOMItemsByViewID 删除BOM视图的所有项
func (r *Repository) DeleteBOMItemsByViewID(ctx context.Context, bomViewID uint) error {
	return r.db.WithContext(ctx).Where("bom_view_id = ?", bomViewID).Delete(&model.BOMItem{}).Error
}

// CountBOMItemsByViewID 统计BOM视图的项数量
func (r *Repository) CountBOMItemsByViewID(ctx context.Context, bomViewID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.BOMItem{}).Where("bom_view_id = ?", bomViewID).Count(&count).Error
	return count, err
}

// GetBOMItemByMaterialInBOM 检查物料是否已在BOM中存在
func (r *Repository) GetBOMItemByMaterialInBOM(ctx context.Context, bomViewID uint, materialID uint) (*model.BOMItem, error) {
	var item model.BOMItem
	err := r.db.WithContext(ctx).
		Where("bom_view_id = ? AND material_id = ?", bomViewID, materialID).
		First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// ExistsRootMaterial 检查根物料是否已有BOM视图
func (r *Repository) ExistsRootMaterial(ctx context.Context, rootMaterialID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.BOMView{}).Where("root_material_id = ?", rootMaterialID).Count(&count).Error
	return count > 0, err
}
