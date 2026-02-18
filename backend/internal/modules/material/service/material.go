package service

import (
	"context"
	"errors"

	"plm/internal/modules/material/model"
	"plm/internal/modules/material/repository"
	sharedModel "plm/internal/shared/model"

	"gorm.io/gorm"
)

type Service struct {
	repo *repository.Repository
	db   *gorm.DB
}

func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

// SetDB 设置数据库连接
func (s *Service) SetDB(db *gorm.DB) {
	s.db = db
}

// CreateMaterial 创建物料
func (s *Service) CreateMaterial(ctx context.Context, req *model.CreateMaterialRequest, createdBy uint) (*model.Material, error) {
	// 检查物料编码是否已存在
	exists, err := s.repo.ExistsByItemID(ctx, req.ItemID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("物料编码已存在")
	}

	// 处理空attributes
	attributes := req.Attributes
	if attributes == "" {
		attributes = "{}"
	}

	material := &model.Material{
		ItemID:      req.ItemID,
		ItemName:    req.ItemName,
		Description: req.Description,
		ItemType:    req.ItemType,
		Version:     "AA", // 初始版本
		Unit:        req.Unit,
		Status:      "DRAFT",
		Attributes:  attributes,
		CreatedBy:   createdBy,
	}

	if err := s.repo.CreateMaterial(ctx, material); err != nil {
		return nil, err
	}

	return material, nil
}

// GetMaterial 获取物料详情
func (s *Service) GetMaterial(ctx context.Context, id uint) (*model.MaterialResponse, error) {
	material, err := s.repo.GetMaterialByID(ctx, id)
	if err != nil {
		return nil, errors.New("物料不存在")
	}

	attrs, _ := s.repo.GetAttributesByMaterialID(ctx, id)

	return &model.MaterialResponse{
		Material:   material,
		Attributes: attrs,
	}, nil
}

// UpdateMaterial 更新物料
func (s *Service) UpdateMaterial(ctx context.Context, id uint, req *model.UpdateMaterialRequest) (*model.Material, error) {
	material, err := s.repo.GetMaterialByID(ctx, id)
	if err != nil {
		return nil, errors.New("物料不存在")
	}

	// 检查状态是否允许编辑
	if material.Status == "RELEASED" {
		return nil, errors.New("已发布的物料不能编辑")
	}
	if material.Status == "REVIEWING" {
		return nil, errors.New("审核中的物料不能编辑")
	}

	if req.ItemName != "" {
		material.ItemName = req.ItemName
	}
	if req.Description != "" {
		material.Description = req.Description
	}
	if req.ItemType != "" {
		material.ItemType = req.ItemType
	}
	if req.Unit != "" {
		material.Unit = req.Unit
	}
	if req.Status != "" {
		material.Status = req.Status
	}
	if req.Attributes != "" {
		material.Attributes = req.Attributes
	}

	if err := s.repo.UpdateMaterial(ctx, material); err != nil {
		return nil, err
	}

	return material, nil
}

// DeleteMaterial 删除物料
func (s *Service) DeleteMaterial(ctx context.Context, id uint) error {
	material, err := s.repo.GetMaterialByID(ctx, id)
	if err != nil {
		return errors.New("物料不存在")
	}

	// 检查状态是否允许删除
	if material.Status == "RELEASED" {
		return errors.New("已发布的物料不能删除")
	}
	if material.Status == "REVIEWING" {
		return errors.New("审核中的物料不能删除")
	}

	// 删除物料属性
	s.repo.DeleteAttributesByMaterialID(ctx, id)

	return s.repo.DeleteMaterial(ctx, id)
}

// ListMaterials 物料列表
func (s *Service) ListMaterials(ctx context.Context, query *model.MaterialListQuery) ([]*model.Material, int64, error) {
	// 设置默认值
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 20
	}

	return s.repo.ListMaterials(ctx, query)
}

// SearchMaterials 搜索物料
func (s *Service) SearchMaterials(ctx context.Context, keyword string, page, pageSize int) ([]*model.Material, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	return s.repo.SearchMaterials(ctx, keyword, page, pageSize)
}

// IncrementVersion 版本升级
func (s *Service) IncrementVersion(ctx context.Context, id uint, reason string, createdBy uint) (*model.Material, error) {
	material, err := s.repo.GetMaterialByID(ctx, id)
	if err != nil {
		return nil, errors.New("物料不存在")
	}

	oldVersion := material.Version
	newVersion := s.calculateNextVersion(material.Version)

	// 使用事务
	if s.db != nil {
		err = s.db.Transaction(func(tx *gorm.DB) error {
			// 更新物料版本
			material.Version = newVersion
			if err := s.repo.UpdateMaterial(ctx, material); err != nil {
				return err
			}

			// 记录版本历史
			history := &sharedModel.VersionHistory{
				EntityType:   "MATERIAL",
				EntityID:     id,
				OldVersion:   oldVersion,
				NewVersion:   newVersion,
				ChangeReason: reason,
				CreatedBy:    createdBy,
			}
			return tx.Create(history).Error
		})
		if err != nil {
			return nil, err
		}
	} else {
		// 没有数据库连接时只更新版本
		material.Version = newVersion
		if err := s.repo.UpdateMaterial(ctx, material); err != nil {
			return nil, err
		}
	}

	return s.repo.GetMaterialByID(ctx, id)
}

// GetVersionHistory 获取物料版本历史
func (s *Service) GetVersionHistory(ctx context.Context, id uint) ([]sharedModel.VersionHistory, error) {
	var histories []sharedModel.VersionHistory
	if s.db == nil {
		return histories, nil
	}

	err := s.db.WithContext(ctx).
		Where("entity_type = ? AND entity_id = ?", "MATERIAL", id).
		Order("created_at DESC").
		Find(&histories).Error
	return histories, err
}

// calculateNextVersion 计算下一个版本号
// 版本规则：AA-YY，跳过I/O/Z
func (s *Service) calculateNextVersion(current string) string {
	if len(current) != 2 {
		return "AA"
	}

	// 跳过的字母
	skipChars := map[byte]bool{'I': true, 'O': true, 'Z': true}

	first := current[0]
	second := current[1]

	// 递增第二位
	second++
	for skipChars[second] {
		second++
	}

	// 如果第二位超过'Y'，进位
	if second > 'Y' {
		second = 'A'
		first++
		for skipChars[first] {
			first++
		}
	}

	// 如果第一位超过'Y'，回到AA
	if first > 'Y' {
		return "AA"
	}

	return string([]byte{first, second})
}
