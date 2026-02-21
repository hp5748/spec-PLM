package service

import (
	"context"
	"errors"
	"fmt"

	"plm/internal/modules/workflow/model"
	"plm/internal/modules/workflow/repository"

	"gorm.io/gorm"
)

type Service struct {
	repo *repository.Repository
	db   *gorm.DB
}

func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

// SetDB 设置数据库连接（用于事务）
func (s *Service) SetDB(db *gorm.DB) {
	s.db = db
	s.repo = repository.NewRepository(db)
}

// ==================== 流程定义相关 ====================

// CreateDefinition 创建流程定义
func (s *Service) CreateDefinition(ctx context.Context, req *model.CreateWorkflowDefinitionRequest, createdBy uint) (*model.WorkflowDefinition, error) {
	// 检查编码是否已存在
	exists, err := s.repo.ExistsDefinitionByCode(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("流程编码已存在")
	}

	// 映射业务类型到流程类型
	workflowType := req.Type
	if workflowType == "" {
		switch req.Code {
		case "MATERIAL_APPROVAL":
			workflowType = "MATERIAL_APPROVAL"
		case "DOCUMENT_APPROVAL":
			workflowType = "DOCUMENT_APPROVAL"
		case "BOM_APPROVAL":
			workflowType = "BOM_APPROVAL"
		default:
			workflowType = req.Type
		}
	}

	def := &model.WorkflowDefinition{
		Name:        req.Name,
		Code:        req.Code,
		Type:        workflowType,
		Version:     "AA",
		Config:      req.Config,
		Status:      "DRAFT",
		Description: req.Description,
		CreatedBy:   createdBy,
	}

	if err := s.repo.CreateDefinition(ctx, def); err != nil {
		return nil, err
	}

	return s.repo.GetDefinitionByID(ctx, def.ID)
}

// GetDefinition 获取流程定义详情
func (s *Service) GetDefinition(ctx context.Context, id uint) (*model.WorkflowDefinition, error) {
	return s.repo.GetDefinitionByID(ctx, id)
}

// UpdateDefinition 更新流程定义
func (s *Service) UpdateDefinition(ctx context.Context, id uint, req *model.UpdateWorkflowDefinitionRequest) (*model.WorkflowDefinition, error) {
	def, err := s.repo.GetDefinitionByID(ctx, id)
	if err != nil {
		return nil, errors.New("流程定义不存在")
	}

	// 已发布的流程不能修改
	if def.Status == "RELEASED" {
		return nil, errors.New("已发布的流程定义不能修改")
	}

	if req.Name != "" {
		def.Name = req.Name
	}
	if req.Config != "" {
		def.Config = req.Config
	}
	if req.Description != "" {
		def.Description = req.Description
	}

	if err := s.repo.UpdateDefinition(ctx, def); err != nil {
		return nil, err
	}

	return s.repo.GetDefinitionByID(ctx, id)
}

// ReleaseDefinition 发布流程定义
func (s *Service) ReleaseDefinition(ctx context.Context, id uint) (*model.WorkflowDefinition, error) {
	def, err := s.repo.GetDefinitionByID(ctx, id)
	if err != nil {
		return nil, errors.New("流程定义不存在")
	}

	if def.Status == "RELEASED" {
		return nil, errors.New("流程定义已发布")
	}

	def.Status = "RELEASED"
	if err := s.repo.UpdateDefinition(ctx, def); err != nil {
		return nil, err
	}

	return s.repo.GetDefinitionByID(ctx, id)
}

// DeleteDefinition 删除流程定义
func (s *Service) DeleteDefinition(ctx context.Context, id uint) error {
	def, err := s.repo.GetDefinitionByID(ctx, id)
	if err != nil {
		return errors.New("流程定义不存在")
	}

	if def.Status == "RELEASED" {
		return errors.New("已发布的流程定义不能删除")
	}

	return s.repo.DeleteDefinition(ctx, id)
}

// ListDefinitions 流程定义列表
func (s *Service) ListDefinitions(ctx context.Context, query *model.WorkflowDefinitionListQuery) ([]*model.WorkflowDefinition, int64, error) {
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 20
	}

	return s.repo.ListDefinitions(ctx, query)
}

// GetActiveDefinitionByType 获取指定类型的激活流程定义
func (s *Service) GetActiveDefinitionByType(ctx context.Context, workflowType string) (*model.WorkflowDefinition, error) {
	return s.repo.GetActiveDefinitionByType(ctx, workflowType)
}

// ==================== 流程实例相关 ====================

// InitiateWorkflow 发起流程
func (s *Service) InitiateWorkflow(ctx context.Context, req *model.InitiateWorkflowRequest, initiatorID uint) (*model.WorkflowInstanceResponse, error) {
	// 获取流程定义
	def, err := s.repo.GetDefinitionByID(ctx, req.DefinitionID)
	if err != nil {
		return nil, errors.New("流程定义不存在")
	}

	if def.Status != "RELEASED" {
		return nil, errors.New("流程定义未发布")
	}

	// 检查业务对象是否已有进行中的流程
	hasActive, err := s.repo.HasActiveInstance(ctx, req.BusinessType, req.BusinessID)
	if err != nil {
		return nil, err
	}
	if hasActive {
		return nil, errors.New("该业务对象已有进行中的审批流程")
	}

	// 检查业务对象状态
	if err := s.checkBusinessStatus(ctx, req.BusinessType, req.BusinessID); err != nil {
		return nil, err
	}

	instance := &model.WorkflowInstance{
		DefinitionID: req.DefinitionID,
		BusinessType: req.BusinessType,
		BusinessID:   req.BusinessID,
		Title:        req.Title,
		Status:       "PENDING",
		InitiatorID:  initiatorID,
		CurrentNode:  "审批",
		Comment:      req.Comment,
	}

	if err := s.repo.CreateInstance(ctx, instance); err != nil {
		return nil, err
	}

	// 更新业务对象状态为审核中
	if err := s.updateBusinessStatus(ctx, req.BusinessType, req.BusinessID, "REVIEWING"); err != nil {
		return nil, err
	}

	// 创建流程历史
	history := &model.WorkflowHistory{
		InstanceID:   instance.ID,
		NodeName:     "发起",
		Action:       "SUBMIT",
		OperatorID:   initiatorID,
		OperatorName: s.getUserName(ctx, initiatorID),
		Comment:      req.Comment,
	}
	s.repo.CreateHistory(ctx, history)

	return s.repo.GetInstanceWithBusiness(ctx, instance.ID)
}

// GetInstance 获取流程实例详情
func (s *Service) GetInstance(ctx context.Context, id uint) (*model.WorkflowInstanceResponse, error) {
	return s.repo.GetInstanceWithBusiness(ctx, id)
}

// ListInstances 流程实例列表
func (s *Service) ListInstances(ctx context.Context, query *model.WorkflowListQuery) ([]*model.WorkflowInstanceResponse, int64, error) {
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 20
	}

	return s.repo.ListInstances(ctx, query)
}

// GetMyTodos 获取我的待办
func (s *Service) GetMyTodos(ctx context.Context, userID uint, query *model.TodoListQuery) ([]*model.TodoItemResponse, int64, error) {
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 20
	}

	return s.repo.GetPendingInstancesByUser(ctx, userID, query)
}

// GetTodoCount 获取待办数量
func (s *Service) GetTodoCount(ctx context.Context, userID uint) (int64, error) {
	return s.repo.CountPendingByUser(ctx, userID)
}

// GetHistories 获取流程历史
func (s *Service) GetHistories(ctx context.Context, instanceID uint) ([]model.WorkflowHistory, error) {
	return s.repo.GetHistoriesByInstanceID(ctx, instanceID)
}

// Approve 同意审批
func (s *Service) Approve(ctx context.Context, instanceID uint, operatorID uint, req *model.ApproveRequest) (*model.WorkflowInstanceResponse, error) {
	instance, err := s.repo.GetInstanceByID(ctx, instanceID)
	if err != nil {
		return nil, errors.New("流程实例不存在")
	}

	if instance.Status != "PENDING" {
		return nil, errors.New("当前流程状态不允许审批")
	}

	// 更新流程状态
	instance.Status = "APPROVED"
	instance.CurrentNode = "已通过"
	if err := s.repo.UpdateInstance(ctx, instance); err != nil {
		return nil, err
	}

	// 更新业务对象状态为已发布
	if err := s.updateBusinessStatus(ctx, instance.BusinessType, instance.BusinessID, "RELEASED"); err != nil {
		return nil, err
	}

	// 创建流程历史
	history := &model.WorkflowHistory{
		InstanceID:   instanceID,
		NodeName:     "审批",
		Action:       "APPROVE",
		OperatorID:   operatorID,
		OperatorName: s.getUserName(ctx, operatorID),
		Comment:      req.Comment,
	}
	s.repo.CreateHistory(ctx, history)

	return s.repo.GetInstanceWithBusiness(ctx, instanceID)
}

// Reject 驳回审批
func (s *Service) Reject(ctx context.Context, instanceID uint, operatorID uint, req *model.RejectRequest) (*model.WorkflowInstanceResponse, error) {
	instance, err := s.repo.GetInstanceByID(ctx, instanceID)
	if err != nil {
		return nil, errors.New("流程实例不存在")
	}

	if instance.Status != "PENDING" {
		return nil, errors.New("当前流程状态不允许审批")
	}

	// 更新流程状态
	instance.Status = "REJECTED"
	instance.CurrentNode = "已驳回"
	if err := s.repo.UpdateInstance(ctx, instance); err != nil {
		return nil, err
	}

	// 更新业务对象状态为已驳回
	if err := s.updateBusinessStatus(ctx, instance.BusinessType, instance.BusinessID, "REJECTED"); err != nil {
		return nil, err
	}

	// 创建流程历史
	history := &model.WorkflowHistory{
		InstanceID:   instanceID,
		NodeName:     "审批",
		Action:       "REJECT",
		OperatorID:   operatorID,
		OperatorName: s.getUserName(ctx, operatorID),
		Comment:      req.Comment,
	}
	s.repo.CreateHistory(ctx, history)

	return s.repo.GetInstanceWithBusiness(ctx, instanceID)
}

// Transfer 转交审批
func (s *Service) Transfer(ctx context.Context, instanceID uint, operatorID uint, req *model.TransferRequest) (*model.WorkflowInstanceResponse, error) {
	instance, err := s.repo.GetInstanceByID(ctx, instanceID)
	if err != nil {
		return nil, errors.New("流程实例不存在")
	}

	if instance.Status != "PENDING" {
		return nil, errors.New("当前流程状态不允许转交")
	}

	// 创建流程历史
	history := &model.WorkflowHistory{
		InstanceID:   instanceID,
		NodeName:     "审批",
		Action:       "TRANSFER",
		OperatorID:   operatorID,
		OperatorName: s.getUserName(ctx, operatorID),
		Comment:      fmt.Sprintf("转交给用户ID: %d. %s", req.TargetUserID, req.Comment),
	}
	s.repo.CreateHistory(ctx, history)

	return s.repo.GetInstanceWithBusiness(ctx, instanceID)
}

// Withdraw 撤回流程
func (s *Service) Withdraw(ctx context.Context, instanceID uint, operatorID uint) (*model.WorkflowInstanceResponse, error) {
	instance, err := s.repo.GetInstanceByID(ctx, instanceID)
	if err != nil {
		return nil, errors.New("流程实例不存在")
	}

	if instance.Status != "PENDING" {
		return nil, errors.New("当前流程状态不允许撤回")
	}

	// 只有发起人可以撤回
	if instance.InitiatorID != operatorID {
		return nil, errors.New("只有发起人可以撤回流程")
	}

	// 更新流程状态
	instance.Status = "CANCELLED"
	instance.CurrentNode = "已撤回"
	if err := s.repo.UpdateInstance(ctx, instance); err != nil {
		return nil, err
	}

	// 更新业务对象状态为草稿
	if err := s.updateBusinessStatus(ctx, instance.BusinessType, instance.BusinessID, "DRAFT"); err != nil {
		return nil, err
	}

	// 创建流程历史
	history := &model.WorkflowHistory{
		InstanceID:   instanceID,
		NodeName:     "审批",
		Action:       "WITHDRAW",
		OperatorID:   operatorID,
		OperatorName: s.getUserName(ctx, operatorID),
		Comment:      "发起人撤回",
	}
	s.repo.CreateHistory(ctx, history)

	return s.repo.GetInstanceWithBusiness(ctx, instanceID)
}

// ==================== 辅助方法 ====================

// checkBusinessStatus 检查业务对象状态是否允许发起流程
func (s *Service) checkBusinessStatus(ctx context.Context, businessType string, businessID uint) error {
	switch businessType {
	case "MATERIAL":
		var status string
		if err := s.db.WithContext(ctx).Table("materials").
			Select("status").Where("id = ?", businessID).Scan(&status).Error; err != nil {
			return errors.New("物料不存在")
		}
		if status != "DRAFT" && status != "REJECTED" {
			return errors.New("物料状态不允许发起审批")
		}
	case "DOCUMENT":
		var status string
		if err := s.db.WithContext(ctx).Table("documents").
			Select("status").Where("id = ?", businessID).Scan(&status).Error; err != nil {
			return errors.New("文档不存在")
		}
		if status != "DRAFT" && status != "REJECTED" {
			return errors.New("文档状态不允许发起审批")
		}
	case "BOM":
		var status string
		if err := s.db.WithContext(ctx).Table("bom_views").
			Select("status").Where("id = ?", businessID).Scan(&status).Error; err != nil {
			return errors.New("BOM不存在")
		}
		if status != "DRAFT" && status != "REJECTED" {
			return errors.New("BOM状态不允许发起审批")
		}
	}
	return nil
}

// updateBusinessStatus 更新业务对象状态
func (s *Service) updateBusinessStatus(ctx context.Context, businessType string, businessID uint, status string) error {
	switch businessType {
	case "MATERIAL":
		return s.db.WithContext(ctx).Table("materials").
			Where("id = ?", businessID).Update("status", status).Error
	case "DOCUMENT":
		return s.db.WithContext(ctx).Table("documents").
			Where("id = ?", businessID).Update("status", status).Error
	case "BOM":
		return s.db.WithContext(ctx).Table("bom_views").
			Where("id = ?", businessID).Update("status", status).Error
	}
	return nil
}

// getUserName 获取用户名
func (s *Service) getUserName(ctx context.Context, userID uint) string {
	var user model.User
	if err := s.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		return ""
	}
	if user.RealName != "" {
		return user.RealName
	}
	return user.Username
}
