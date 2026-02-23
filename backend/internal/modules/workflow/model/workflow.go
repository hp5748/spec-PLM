package model

import (
	"time"

	"gorm.io/gorm"
)

// WorkflowDefinition 流程定义模型
type WorkflowDefinition struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:200;not null" json:"name"`
	Code        string         `gorm:"size:50;not null" json:"code"`
	Type        string         `gorm:"size:50;not null;index" json:"type"` // BOM_APPROVAL, DOC_APPROVAL, MATERIAL_APPROVAL
	Version     string         `gorm:"size:2;not null;default:AA" json:"version"`
	Config      string         `gorm:"type:json;not null" json:"config"` // 流程配置JSON
	Status      string         `gorm:"size:20;default:DRAFT;index" json:"status"` // DRAFT, RELEASED, OBSOLETE
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	CreatedBy   uint           `gorm:"not null" json:"created_by"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// WorkflowInstance 流程实例模型
type WorkflowInstance struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	DefinitionID  uint           `gorm:"not null;index" json:"definition_id"`
	BusinessType  string         `gorm:"size:50;not null;index" json:"business_type"` // MATERIAL, DOCUMENT, BOM
	BusinessID    uint           `gorm:"not null;index:idx_business" json:"business_id"`
	Title         string         `gorm:"size:200" json:"title"`
	Status        string         `gorm:"size:20;default:DRAFT;index" json:"status"` // DRAFT, PENDING, APPROVED, REJECTED, CANCELLED
	InitiatorID   uint           `gorm:"not null;index" json:"initiator_id"`
	CurrentNode   string         `gorm:"size:100" json:"current_node"`                     // 当前节点名称
	CurrentNodeID string         `gorm:"size:50" json:"current_node_id"`                   // 当前节点ID
	NodePath      string         `gorm:"type:text" json:"node_path"`                       // 经过的节点路径（JSON数组）
	ApprovedNodes string         `gorm:"type:text" json:"approved_nodes"`                  // 已审批节点列表（JSON数组）
	Comment       string         `json:"comment"`                                          // 发起时的意见
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Definition *WorkflowDefinition `gorm:"foreignKey:DefinitionID" json:"definition,omitempty"`
	Initiator  *User               `gorm:"foreignKey:InitiatorID" json:"initiator,omitempty"`
}

// WorkflowHistory 流程历史模型
type WorkflowHistory struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	InstanceID   uint      `gorm:"not null;index" json:"instance_id"`
	NodeName     string    `gorm:"size:100" json:"node_name"`
	FromNodeID   string    `gorm:"size:50" json:"from_node_id"`   // 源节点ID
	ToNodeID     string    `gorm:"size:50" json:"to_node_id"`     // 目标节点ID
	Action       string    `gorm:"size:20;not null" json:"action"` // SUBMIT, APPROVE, REJECT, TRANSFER, CANCEL, WITHDRAW
	OperatorID   uint      `gorm:"not null;index" json:"operator_id"`
	OperatorName string    `gorm:"size:50" json:"operator_name"`
	Comment      string    `json:"comment"`
	CreatedAt    time.Time `json:"created_at"`

	// 关联
	Operator *User `gorm:"foreignKey:OperatorID" json:"operator,omitempty"`
}

// User 简化的用户模型（用于关联查询）
type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"size:50;not null" json:"username"`
	RealName string `gorm:"size:50" json:"real_name"`
}

// CreateWorkflowDefinitionRequest 创建流程定义请求
type CreateWorkflowDefinitionRequest struct {
	Name        string `json:"name" binding:"required,max=200"`
	Code        string `json:"code" binding:"required,max=50"`
	Type        string `json:"type" binding:"required,oneof=MATERIAL_APPROVAL DOCUMENT_APPROVAL BOM_APPROVAL"`
	Config      string `json:"config" binding:"required"`
	Description string `json:"description"`
}

// UpdateWorkflowDefinitionRequest 更新流程定义请求
type UpdateWorkflowDefinitionRequest struct {
	Name        string `json:"name" binding:"max=200"`
	Config      string `json:"config"`
	Description string `json:"description"`
}

// InitiateWorkflowRequest 发起流程请求
type InitiateWorkflowRequest struct {
	DefinitionID uint   `json:"definition_id" binding:"required"`
	BusinessType string `json:"business_type" binding:"required,oneof=MATERIAL DOCUMENT BOM"`
	BusinessID   uint   `json:"business_id" binding:"required"`
	Title        string `json:"title" binding:"required,max=200"`
	Comment      string `json:"comment"`
}

// ApproveRequest 同意审批请求
type ApproveRequest struct {
	Comment string `json:"comment"`
}

// RejectRequest 驳回审批请求
type RejectRequest struct {
	Comment string `json:"comment" binding:"required"`
}

// TransferRequest 转交审批请求
type TransferRequest struct {
	TargetUserID uint   `json:"target_user_id" binding:"required"`
	Comment      string `json:"comment"`
}

// WorkflowListQuery 流程列表查询参数
type WorkflowListQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Type     string `form:"type"`
	Status   string `form:"status"`
	Keyword  string `form:"keyword"`
}

// TodoListQuery 待办列表查询参数
type TodoListQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Type     string `form:"type"`   // 业务类型筛选
	Status   string `form:"status"`
	Keyword  string `form:"keyword"` // 关键字筛选
}

// WorkflowDefinitionListQuery 流程定义列表查询参数
type WorkflowDefinitionListQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Type     string `form:"type"`
	Status   string `form:"status"`
	Keyword  string `form:"keyword"`
	NodeName string `form:"node_name"` // 按节点名称筛选
}

// WorkflowInstanceResponse 流程实例响应
type WorkflowInstanceResponse struct {
	WorkflowInstance
	BusinessName   string   `json:"business_name"`   // 业务对象名称
	BusinessCode   string   `json:"business_code"`   // 业务对象编码
	BusinessStatus string   `json:"business_status"` // 业务对象状态
	NodePathList   []string `json:"node_path_list"`  // 解析后的节点路径
	ApprovedList   []string `json:"approved_list"`   // 解析后的已审批节点
}

// TodoItemResponse 待办事项响应
type TodoItemResponse struct {
	WorkflowInstanceResponse
	PendingAction string `json:"pending_action"` // 待执行操作描述
}

// TableName 指定表名
func (WorkflowDefinition) TableName() string {
	return "workflow_definitions"
}

func (WorkflowInstance) TableName() string {
	return "workflow_instances"
}

func (WorkflowHistory) TableName() string {
	return "workflow_histories"
}

func (User) TableName() string {
	return "users"
}
