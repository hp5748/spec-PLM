package repository

import (
	"context"

	"plm/internal/modules/script/model"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// ==================== 脚本管理 ====================

// CreateScript 创建脚本
func (r *Repository) CreateScript(ctx context.Context, script *model.Script) error {
	return r.db.WithContext(ctx).Create(script).Error
}

// GetScriptByID 根据ID获取脚本
func (r *Repository) GetScriptByID(ctx context.Context, id uint) (*model.Script, error) {
	var script model.Script
	err := r.db.WithContext(ctx).First(&script, id).Error
	if err != nil {
		return nil, err
	}
	return &script, nil
}

// GetScriptByCode 根据编码获取脚本
func (r *Repository) GetScriptByCode(ctx context.Context, code string) (*model.Script, error) {
	var script model.Script
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&script).Error
	if err != nil {
		return nil, err
	}
	return &script, nil
}

// UpdateScript 更新脚本
func (r *Repository) UpdateScript(ctx context.Context, script *model.Script) error {
	return r.db.WithContext(ctx).Save(script).Error
}

// DeleteScript 删除脚本
func (r *Repository) DeleteScript(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Script{}, id).Error
}

// ExistsScriptByCode 检查脚本编码是否存在
func (r *Repository) ExistsScriptByCode(ctx context.Context, code string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Script{}).Where("code = ?", code).Count(&count).Error
	return count > 0, err
}

// ListScripts 脚本列表
func (r *Repository) ListScripts(ctx context.Context, query *model.ScriptListQuery) ([]*model.Script, int64, error) {
	var scripts []*model.Script
	var total int64

	db := r.db.WithContext(ctx).Model(&model.Script{})

	if query.Type != "" {
		db = db.Where("type = ?", query.Type)
	}
	if query.TriggerType != "" {
		db = db.Where("trigger_type = ?", query.TriggerType)
	}
	if query.BusinessType != "" {
		db = db.Where("business_type = ? OR business_type = ?", query.BusinessType, "ALL")
	}
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}
	if query.Keyword != "" {
		db = db.Where("name LIKE ? OR code LIKE ?", "%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (query.Page - 1) * query.PageSize
	if err := db.Order("id DESC").Offset(offset).Limit(query.PageSize).Find(&scripts).Error; err != nil {
		return nil, 0, err
	}

	return scripts, total, nil
}

// GetScriptsByTrigger 根据触发类型和业务类型获取脚本列表
func (r *Repository) GetScriptsByTrigger(ctx context.Context, triggerType, businessType string) ([]*model.Script, error) {
	var scripts []*model.Script
	err := r.db.WithContext(ctx).
		Where("trigger_type = ? AND status = ?", triggerType, "RELEASED").
		Where("business_type = ? OR business_type = ?", businessType, "ALL").
		Order("id ASC").
		Find(&scripts).Error
	return scripts, err
}

// ==================== 执行日志 ====================

// CreateExecutionLog 创建执行日志
func (r *Repository) CreateExecutionLog(ctx context.Context, log *model.ScriptExecutionLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

// GetExecutionLogs 获取执行日志列表
func (r *Repository) GetExecutionLogs(ctx context.Context, scriptID uint, page, pageSize int) ([]*model.ScriptExecutionLog, int64, error) {
	var logs []*model.ScriptExecutionLog
	var total int64

	db := r.db.WithContext(ctx).Model(&model.ScriptExecutionLog{})

	if scriptID > 0 {
		db = db.Where("script_id = ?", scriptID)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Order("id DESC").Offset(offset).Limit(pageSize).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetExecutionLogByID 获取执行日志详情
func (r *Repository) GetExecutionLogByID(ctx context.Context, id uint) (*model.ScriptExecutionLog, error) {
	var log model.ScriptExecutionLog
	err := r.db.WithContext(ctx).First(&log, id).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}
