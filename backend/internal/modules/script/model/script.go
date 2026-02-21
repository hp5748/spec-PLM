package model

import (
	"time"

	"gorm.io/gorm"
)

// Script 脚本模型
type Script struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:200;not null" json:"name"`
	Code        string         `gorm:"size:50;not null;uniqueIndex" json:"code"`
	Type        string         `gorm:"size:20;not null;index" json:"type"` // JAVASCRIPT, SQL
	TriggerType string         `gorm:"size:50;not null;index" json:"trigger_type"` // BEFORE_SUBMIT, AFTER_SUBMIT, BEFORE_APPROVE, AFTER_APPROVE, BEFORE_REJECT, AFTER_REJECT
	BusinessType string        `gorm:"size:50;not null;index" json:"business_type"` // MATERIAL, DOCUMENT, BOM, ALL
	Content     string         `gorm:"type:text;not null" json:"content"` // 脚本内容
	Description string         `json:"description"`
	Status      string         `gorm:"size:20;default:DRAFT;index" json:"status"` // DRAFT, RELEASED, OBSOLETE
	Timeout     int            `gorm:"default:5000" json:"timeout"` // 执行超时时间（毫秒）
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	CreatedBy   uint           `gorm:"not null" json:"created_by"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// ScriptExecutionLog 脚本执行日志
type ScriptExecutionLog struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	ScriptID     uint      `gorm:"not null;index" json:"script_id"`
	ScriptName   string    `gorm:"size:200" json:"script_name"`
	InstanceID   uint      `gorm:"index" json:"instance_id"` // 流程实例ID
	BusinessType string    `gorm:"size:50;index" json:"business_type"`
	BusinessID   uint      `gorm:"index" json:"business_id"`
	TriggerType  string    `gorm:"size:50" json:"trigger_type"`
	Input        string    `gorm:"type:text" json:"input"` // 输入参数JSON
	Output       string    `gorm:"type:text" json:"output"` // 输出结果JSON
	Success      bool      `json:"success"`
	ErrorMsg     string    `gorm:"type:text" json:"error_msg"` // 错误信息
	Duration     int       `json:"duration"` // 执行耗时（毫秒）
	ExecutedAt   time.Time `json:"executed_at"`
	ExecutedBy   uint      `json:"executed_by"`
}

// CreateScriptRequest 创建脚本请求
type CreateScriptRequest struct {
	Name         string `json:"name" binding:"required,max=200"`
	Code         string `json:"code" binding:"required,max=50"`
	Type         string `json:"type" binding:"required,oneof=JAVASCRIPT SQL"`
	TriggerType  string `json:"trigger_type" binding:"required,oneof=BEFORE_SUBMIT AFTER_SUBMIT BEFORE_APPROVE AFTER_APPROVE BEFORE_REJECT AFTER_REJECT"`
	BusinessType string `json:"business_type" binding:"required,oneof=MATERIAL DOCUMENT BOM ALL"`
	Content      string `json:"content" binding:"required"`
	Description  string `json:"description"`
	Timeout      int    `json:"timeout"` // 超时时间（毫秒）
}

// UpdateScriptRequest 更新脚本请求
type UpdateScriptRequest struct {
	Name         string `json:"name" binding:"max=200"`
	Content      string `json:"content"`
	Description  string `json:"description"`
	Timeout      int    `json:"timeout"`
}

// ExecuteScriptRequest 执行脚本请求
type ExecuteScriptRequest struct {
	ScriptID     uint                   `json:"script_id"`
	BusinessType string                 `json:"business_type"`
	BusinessID   uint                   `json:"business_id"`
	Context      map[string]interface{} `json:"context"` // 执行上下文
}

// ExecuteScriptResult 执行脚本结果
type ExecuteScriptResult struct {
	Success bool                   `json:"success"`
	Output  map[string]interface{} `json:"output"`
	Error   string                 `json:"error"`
	LogID   uint                   `json:"log_id"`
}

// ScriptListQuery 脚本列表查询参数
type ScriptListQuery struct {
	Page         int    `form:"page"`
	PageSize     int    `form:"page_size"`
	Type         string `form:"type"`
	TriggerType  string `form:"trigger_type"`
	BusinessType string `form:"business_type"`
	Status       string `form:"status"`
	Keyword      string `form:"keyword"`
}

// ScriptTestRequest 脚本测试请求
type ScriptTestRequest struct {
	Type         string                 `json:"type" binding:"required,oneof=JAVASCRIPT SQL"`
	Content      string                 `json:"content" binding:"required"`
	BusinessType string                 `json:"business_type"`
	BusinessID   uint                   `json:"business_id"`
	Context      map[string]interface{} `json:"context"`
}

// ScriptTestResult 脚本测试结果
type ScriptTestResult struct {
	Success  bool                   `json:"success"`
	Output   map[string]interface{} `json:"output"`
	Error    string                 `json:"error"`
	Duration int                    `json:"duration"`
	Logs     []string               `json:"logs"` // 执行日志
}

// TableName 指定表名
func (Script) TableName() string {
	return "scripts"
}

func (ScriptExecutionLog) TableName() string {
	return "script_execution_logs"
}
