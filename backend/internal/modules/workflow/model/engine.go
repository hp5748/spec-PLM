package model

// WorkflowConfig 流程配置
type WorkflowConfig struct {
	Version string       `json:"version"` // 配置版本
	Nodes   []NodeConfig `json:"nodes"`   // 节点列表
	Edges   []EdgeConfig `json:"edges"`   // 边（连接线）列表
}

// NodeConfig 节点配置
type NodeConfig struct {
	ID         string          `json:"id"`                   // 节点唯一标识
	Name       string          `json:"name"`                 // 节点名称（显示用）
	Type       string          `json:"type"`                 // 节点类型: start, approval, condition, end
	Assignee   *AssigneeConfig `json:"assignee,omitempty"`   // 审批人配置（审批节点）
	Condition  string          `json:"condition,omitempty"`  // 条件表达式（条件节点）
}

// AssigneeConfig 审批人配置
type AssigneeConfig struct {
	Type  string `json:"type"`  // 审批人类型: role, department, user, initiator
	Value string `json:"value"` // 对应值（角色编码/部门ID/用户ID）
}

// EdgeConfig 边配置
type EdgeConfig struct {
	ID        string `json:"id"`                  // 边的唯一标识
	Source    string `json:"source"`              // 源节点ID
	Target    string `json:"target"`              // 目标节点ID
	Label     string `json:"label,omitempty"`     // 边的标签（用于条件分支）
	Condition string `json:"condition,omitempty"` // 条件表达式
}

// WorkflowContext 流程上下文
type WorkflowContext struct {
	InitiatorID   uint                   `json:"initiator_id"`
	BusinessType  string                 `json:"business_type"`
	BusinessID    uint                   `json:"business_id"`
	Data          map[string]interface{} `json:"data,omitempty"`
	Initiator     map[string]interface{} `json:"initiator,omitempty"`
}

// NodeTypes 节点类型常量
const (
	NodeTypeStart     = "start"
	NodeTypeApproval  = "approval"
	NodeTypeCondition = "condition"
	NodeTypeEnd       = "end"
)

// AssigneeTypes 审批人类型常量
const (
	AssigneeTypeRole       = "role"
	AssigneeTypeDepartment = "department"
	AssigneeTypeUser       = "user"
	AssigneeTypeInitiator  = "initiator"
)
