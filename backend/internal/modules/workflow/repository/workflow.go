package repository

import (
	"context"

	"plm/internal/modules/workflow/model"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// ==================== 流程定义相关 ====================

// CreateDefinition 创建流程定义
func (r *Repository) CreateDefinition(ctx context.Context, def *model.WorkflowDefinition) error {
	return r.db.WithContext(ctx).Create(def).Error
}

// GetDefinitionByID 根据ID获取流程定义
func (r *Repository) GetDefinitionByID(ctx context.Context, id uint) (*model.WorkflowDefinition, error) {
	var def model.WorkflowDefinition
	err := r.db.WithContext(ctx).First(&def, id).Error
	if err != nil {
		return nil, err
	}
	return &def, nil
}

// GetDefinitionByCode 根据编码获取流程定义
func (r *Repository) GetDefinitionByCode(ctx context.Context, code string) (*model.WorkflowDefinition, error) {
	var def model.WorkflowDefinition
	err := r.db.WithContext(ctx).
		Where("code = ? AND status = ?", code, "RELEASED").
		Order("version DESC").
		First(&def).Error
	if err != nil {
		return nil, err
	}
	return &def, nil
}

// GetActiveDefinitionByType 根据类型获取激活的流程定义
func (r *Repository) GetActiveDefinitionByType(ctx context.Context, workflowType string) (*model.WorkflowDefinition, error) {
	var def model.WorkflowDefinition
	err := r.db.WithContext(ctx).
		Where("type = ? AND status = ?", workflowType, "RELEASED").
		Order("version DESC").
		First(&def).Error
	if err != nil {
		return nil, err
	}
	return &def, nil
}

// UpdateDefinition 更新流程定义
func (r *Repository) UpdateDefinition(ctx context.Context, def *model.WorkflowDefinition) error {
	return r.db.WithContext(ctx).Save(def).Error
}

// DeleteDefinition 删除流程定义
func (r *Repository) DeleteDefinition(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.WorkflowDefinition{}, id).Error
}

// ListDefinitions 流程定义列表
func (r *Repository) ListDefinitions(ctx context.Context, query *model.WorkflowDefinitionListQuery) ([]*model.WorkflowDefinition, int64, error) {
	var list []*model.WorkflowDefinition
	var total int64

	db := r.db.WithContext(ctx).Model(&model.WorkflowDefinition{})

	if query.Type != "" {
		db = db.Where("type = ?", query.Type)
	}
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}
	if query.Keyword != "" {
		// 扩展关键字查询范围：名称、编码、描述、配置
		keyword := "%" + query.Keyword + "%"
		db = db.Where("name LIKE ? OR code LIKE ? OR description LIKE ? OR config LIKE ?", keyword, keyword, keyword, keyword)
	}
	if query.NodeName != "" {
		// 按节点名称筛选（JSON字段模糊匹配）
		db = db.Where("config LIKE ?", "%\"name\":\""+query.NodeName+"\"%")
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

// ExistsDefinitionByCode 检查编码是否已存在
func (r *Repository) ExistsDefinitionByCode(ctx context.Context, code string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.WorkflowDefinition{}).
		Where("code = ?", code).Count(&count).Error
	return count > 0, err
}

// ==================== 流程实例相关 ====================

// CreateInstance 创建流程实例
func (r *Repository) CreateInstance(ctx context.Context, instance *model.WorkflowInstance) error {
	return r.db.WithContext(ctx).Create(instance).Error
}

// GetInstanceByID 根据ID获取流程实例
func (r *Repository) GetInstanceByID(ctx context.Context, id uint) (*model.WorkflowInstance, error) {
	var instance model.WorkflowInstance
	err := r.db.WithContext(ctx).
		Preload("Definition").
		Preload("Initiator").
		First(&instance, id).Error
	if err != nil {
		return nil, err
	}
	return &instance, nil
}

// GetInstanceWithBusiness 根据ID获取流程实例（含业务对象信息）
func (r *Repository) GetInstanceWithBusiness(ctx context.Context, id uint) (*model.WorkflowInstanceResponse, error) {
	var instance model.WorkflowInstance
	err := r.db.WithContext(ctx).
		Preload("Definition").
		Preload("Initiator").
		First(&instance, id).Error
	if err != nil {
		return nil, err
	}

	resp := &model.WorkflowInstanceResponse{
		WorkflowInstance: instance,
	}

	// 获取业务对象信息
	switch instance.BusinessType {
	case "MATERIAL":
		var material struct {
			ItemID   string `gorm:"column:item_id"`
			ItemName string `gorm:"column:item_name"`
			Status   string `gorm:"column:status"`
		}
		if err := r.db.WithContext(ctx).Table("materials").
			Select("item_id, item_name, status").
			Where("id = ?", instance.BusinessID).
			First(&material).Error; err == nil {
			resp.BusinessCode = material.ItemID
			resp.BusinessName = material.ItemName
			resp.BusinessStatus = material.Status
		}
	case "DOCUMENT":
		var doc struct {
			DocID   string `gorm:"column:doc_id"`
			DocName string `gorm:"column:doc_name"`
			Status  string `gorm:"column:status"`
		}
		if err := r.db.WithContext(ctx).Table("documents").
			Select("doc_id, doc_name, status").
			Where("id = ?", instance.BusinessID).
			First(&doc).Error; err == nil {
			resp.BusinessCode = doc.DocID
			resp.BusinessName = doc.DocName
			resp.BusinessStatus = doc.Status
		}
	case "BOM":
		var bom struct {
			Name   string `gorm:"column:name"`
			Status string `gorm:"column:status"`
		}
		if err := r.db.WithContext(ctx).Table("bom_views").
			Select("name, status").
			Where("id = ?", instance.BusinessID).
			First(&bom).Error; err == nil {
			resp.BusinessCode = "BOM-" + string(rune(instance.BusinessID))
			resp.BusinessName = bom.Name
			resp.BusinessStatus = bom.Status
		}
	}

	return resp, nil
}

// UpdateInstance 更新流程实例
func (r *Repository) UpdateInstance(ctx context.Context, instance *model.WorkflowInstance) error {
	return r.db.WithContext(ctx).Save(instance).Error
}

// DeleteInstance 删除流程实例
func (r *Repository) DeleteInstance(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.WorkflowInstance{}, id).Error
}

// ListInstances 流程实例列表
func (r *Repository) ListInstances(ctx context.Context, query *model.WorkflowListQuery) ([]*model.WorkflowInstanceResponse, int64, error) {
	var instances []model.WorkflowInstance
	var total int64

	db := r.db.WithContext(ctx).Model(&model.WorkflowInstance{}).
		Preload("Definition").
		Preload("Initiator")

	if query.Type != "" {
		db = db.Joins("JOIN workflow_definitions ON workflow_definitions.id = workflow_instances.definition_id").
			Where("workflow_definitions.type = ?", query.Type)
	}
	if query.Status != "" {
		db = db.Where("workflow_instances.status = ?", query.Status)
	}
	if query.Keyword != "" {
		db = db.Where("workflow_instances.title LIKE ?", "%"+query.Keyword+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (query.Page - 1) * query.PageSize
	if err := db.Order("workflow_instances.created_at DESC").Offset(offset).Limit(query.PageSize).Find(&instances).Error; err != nil {
		return nil, 0, err
	}

	var result []*model.WorkflowInstanceResponse
	for _, inst := range instances {
		resp := &model.WorkflowInstanceResponse{WorkflowInstance: inst}
		// 获取业务对象信息
		r.fillBusinessInfo(ctx, resp)
		result = append(result, resp)
	}

	return result, total, nil
}

// GetPendingInstancesByUser 获取用户的待办实例
func (r *Repository) GetPendingInstancesByUser(ctx context.Context, userID uint, query *model.TodoListQuery) ([]*model.TodoItemResponse, int64, error) {
	var instances []model.WorkflowInstance
	var total int64

	db := r.db.WithContext(ctx).Model(&model.WorkflowInstance{}).
		Preload("Definition").
		Preload("Initiator")

	// 支持按业务类型筛选
	if query.Type != "" {
		db = db.Where("workflow_instances.business_type = ?", query.Type)
	}

	// 支持按状态筛选，默认只显示PENDING
	if query.Status != "" {
		db = db.Where("workflow_instances.status = ?", query.Status)
	} else {
		db = db.Where("workflow_instances.status = ?", "PENDING")
	}

	// 支持关键字筛选（标题）
	if query.Keyword != "" {
		db = db.Where("workflow_instances.title LIKE ?", "%"+query.Keyword+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (query.Page - 1) * query.PageSize
	if err := db.Order("workflow_instances.created_at DESC").Offset(offset).Limit(query.PageSize).Find(&instances).Error; err != nil {
		return nil, 0, err
	}

	var result []*model.TodoItemResponse
	for _, inst := range instances {
		resp := &model.TodoItemResponse{
			WorkflowInstanceResponse: model.WorkflowInstanceResponse{WorkflowInstance: inst},
			PendingAction:            "审批",
		}
		r.fillBusinessInfo(ctx, &resp.WorkflowInstanceResponse)
		result = append(result, resp)
	}

	return result, total, nil
}

// CountPendingByUser 统计用户待办数量
func (r *Repository) CountPendingByUser(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.WorkflowInstance{}).
		Where("status = ?", "PENDING").
		Count(&count).Error
	return count, err
}

// GetInstanceByBusiness 根据业务对象获取流程实例
func (r *Repository) GetInstanceByBusiness(ctx context.Context, businessType string, businessID uint) (*model.WorkflowInstance, error) {
	var instance model.WorkflowInstance
	err := r.db.WithContext(ctx).
		Where("business_type = ? AND business_id = ?", businessType, businessID).
		Order("created_at DESC").
		First(&instance).Error
	if err != nil {
		return nil, err
	}
	return &instance, nil
}

// HasActiveInstance 检查业务对象是否有进行中的流程
func (r *Repository) HasActiveInstance(ctx context.Context, businessType string, businessID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.WorkflowInstance{}).
		Where("business_type = ? AND business_id = ? AND status IN (?)", businessType, businessID, []string{"DRAFT", "PENDING"}).
		Count(&count).Error
	return count > 0, err
}

// ==================== 流程历史相关 ====================

// CreateHistory 创建流程历史
func (r *Repository) CreateHistory(ctx context.Context, history *model.WorkflowHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

// GetHistoriesByInstanceID 根据流程实例ID获取历史
func (r *Repository) GetHistoriesByInstanceID(ctx context.Context, instanceID uint) ([]model.WorkflowHistory, error) {
	var histories []model.WorkflowHistory
	err := r.db.WithContext(ctx).
		Preload("Operator").
		Where("instance_id = ?", instanceID).
		Order("created_at ASC").
		Find(&histories).Error
	return histories, err
}

// ==================== 辅助方法 ====================

// fillBusinessInfo 填充业务对象信息
func (r *Repository) fillBusinessInfo(ctx context.Context, resp *model.WorkflowInstanceResponse) {
	switch resp.BusinessType {
	case "MATERIAL":
		var material struct {
			ItemID   string `gorm:"column:item_id"`
			ItemName string `gorm:"column:item_name"`
			Status   string `gorm:"column:status"`
		}
		if err := r.db.WithContext(ctx).Table("materials").
			Select("item_id, item_name, status").
			Where("id = ?", resp.BusinessID).
			First(&material).Error; err == nil {
			resp.BusinessCode = material.ItemID
			resp.BusinessName = material.ItemName
			resp.BusinessStatus = material.Status
		}
	case "DOCUMENT":
		var doc struct {
			DocID   string `gorm:"column:doc_id"`
			DocName string `gorm:"column:doc_name"`
			Status  string `gorm:"column:status"`
		}
		if err := r.db.WithContext(ctx).Table("documents").
			Select("doc_id, doc_name, status").
			Where("id = ?", resp.BusinessID).
			First(&doc).Error; err == nil {
			resp.BusinessCode = doc.DocID
			resp.BusinessName = doc.DocName
			resp.BusinessStatus = doc.Status
		}
	case "BOM":
		var bom struct {
			Name   string `gorm:"column:name"`
			Status string `gorm:"column:status"`
		}
		if err := r.db.WithContext(ctx).Table("bom_views").
			Select("name, status").
			Where("id = ?", resp.BusinessID).
			First(&bom).Error; err == nil {
			resp.BusinessCode = "BOM-" + string(rune(resp.BusinessID))
			resp.BusinessName = bom.Name
			resp.BusinessStatus = bom.Status
		}
	}
}
