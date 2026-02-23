package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"plm/internal/modules/workflow/model"
)

// WorkflowEngine 流程引擎
type WorkflowEngine struct{}

// NewWorkflowEngine 创建流程引擎实例
func NewWorkflowEngine() *WorkflowEngine {
	return &WorkflowEngine{}
}

// ParseConfig 解析流程配置
func (e *WorkflowEngine) ParseConfig(configJSON string) (*model.WorkflowConfig, error) {
	var config model.WorkflowConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return nil, fmt.Errorf("解析流程配置失败: %w", err)
	}

	// 验证配置
	if err := e.validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// validateConfig 验证配置有效性
func (e *WorkflowEngine) validateConfig(config *model.WorkflowConfig) error {
	if len(config.Nodes) == 0 {
		return errors.New("流程配置必须包含节点")
	}

	// 检查必须有start和end节点
	hasStart, hasEnd := false, false
	nodeMap := make(map[string]bool)

	for _, node := range config.Nodes {
		nodeMap[node.ID] = true
		if node.Type == model.NodeTypeStart {
			hasStart = true
		}
		if node.Type == model.NodeTypeEnd {
			hasEnd = true
		}
	}

	if !hasStart {
		return errors.New("流程必须包含开始节点")
	}
	if !hasEnd {
		return errors.New("流程必须包含结束节点")
	}

	// 验证边的引用
	for _, edge := range config.Edges {
		if !nodeMap[edge.Source] {
			return fmt.Errorf("边的源节点 %s 不存在", edge.Source)
		}
		if !nodeMap[edge.Target] {
			return fmt.Errorf("边的目标节点 %s 不存在", edge.Target)
		}
	}

	return nil
}

// GetStartNode 获取开始节点
func (e *WorkflowEngine) GetStartNode(config *model.WorkflowConfig) *model.NodeConfig {
	for i := range config.Nodes {
		if config.Nodes[i].Type == model.NodeTypeStart {
			return &config.Nodes[i]
		}
	}
	return nil
}

// GetNodeByID 根据ID获取节点
func (e *WorkflowEngine) GetNodeByID(config *model.WorkflowConfig, nodeID string) *model.NodeConfig {
	for i := range config.Nodes {
		if config.Nodes[i].ID == nodeID {
			return &config.Nodes[i]
		}
	}
	return nil
}

// GetNextNodes 获取下一个节点列表
func (e *WorkflowEngine) GetNextNodes(config *model.WorkflowConfig, currentNodeID string, context *model.WorkflowContext) ([]model.NodeConfig, error) {
	var nextNodes []model.NodeConfig
	nodeMap := e.buildNodeMap(config)

	// 找到从当前节点出发的所有边
	for _, edge := range config.Edges {
		if edge.Source != currentNodeID {
			continue
		}

		targetNode, exists := nodeMap[edge.Target]
		if !exists {
			continue
		}

		// 如果是条件边，需要评估条件
		if edge.Condition != "" {
			matched, err := e.evaluateCondition(edge.Condition, context)
			if err != nil {
				return nil, err
			}
			if !matched {
				continue
			}
		}

		// 如果目标节点是条件节点，需要继续找下一个节点
		if targetNode.Type == model.NodeTypeCondition {
			conditionNodes, err := e.GetNextNodes(config, targetNode.ID, context)
			if err != nil {
				return nil, err
			}
			nextNodes = append(nextNodes, conditionNodes...)
		} else {
			nextNodes = append(nextNodes, targetNode)
		}
	}

	return nextNodes, nil
}

// GetFirstApprovalNode 获取第一个审批节点
func (e *WorkflowEngine) GetFirstApprovalNode(config *model.WorkflowConfig, context *model.WorkflowContext) (*model.NodeConfig, error) {
	startNode := e.GetStartNode(config)
	if startNode == nil {
		return nil, errors.New("流程配置缺少开始节点")
	}

	nextNodes, err := e.GetNextNodes(config, startNode.ID, context)
	if err != nil {
		return nil, err
	}

	// 获取第一个审批节点
	for _, node := range nextNodes {
		if node.Type == model.NodeTypeApproval {
			return &node, nil
		}
	}

	return nil, errors.New("流程没有审批节点")
}

// buildNodeMap 构建节点映射
func (e *WorkflowEngine) buildNodeMap(config *model.WorkflowConfig) map[string]model.NodeConfig {
	nodeMap := make(map[string]model.NodeConfig)
	for _, node := range config.Nodes {
		nodeMap[node.ID] = node
	}
	return nodeMap
}

// evaluateCondition 评估条件表达式
// 支持简单的条件表达式，如: data.item_type == "ASSEMBLY", data.amount > 1000
func (e *WorkflowEngine) evaluateCondition(condition string, context *model.WorkflowContext) (bool, error) {
	if condition == "" {
		return true, nil
	}

	// 解析条件表达式
	// 格式: field operator value
	// 例如: data.item_type == "ASSEMBLY"
	//       data.amount > 1000

	// 支持的操作符
	operators := []string{"==", "!=", ">=", "<=", ">", "<", "contains"}

	for _, op := range operators {
		if strings.Contains(condition, op) {
			parts := strings.SplitN(condition, op, 2)
			if len(parts) != 2 {
				continue
			}

			field := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			// 获取字段值
			fieldValue, err := e.getFieldValue(field, context)
			if err != nil {
				return false, err
			}

			// 比较值
			return e.compareValues(fieldValue, op, value)
		}
	}

	return true, nil
}

// getFieldValue 获取字段值
func (e *WorkflowEngine) getFieldValue(field string, context *model.WorkflowContext) (interface{}, error) {
	// 支持的字段路径:
	// data.xxx - 业务数据字段
	// initiator.xxx - 发起人信息

	parts := strings.Split(field, ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("无效的字段路径: %s", field)
	}

	switch parts[0] {
	case "data":
		if context.Data == nil {
			return nil, nil
		}
		return context.Data[parts[1]], nil
	case "initiator":
		if context.Initiator == nil {
			return nil, nil
		}
		return context.Initiator[parts[1]], nil
	default:
		return nil, fmt.Errorf("不支持的字段前缀: %s", parts[0])
	}
}

// compareValues 比较值
func (e *WorkflowEngine) compareValues(fieldValue interface{}, operator string, compareValue string) (bool, error) {
	// 移除引号
	compareValue = strings.Trim(compareValue, "\"'")

	switch operator {
	case "==":
		return fmt.Sprintf("%v", fieldValue) == compareValue, nil
	case "!=":
		return fmt.Sprintf("%v", fieldValue) != compareValue, nil
	case ">":
		return e.compareNumbers(fieldValue, compareValue, ">")
	case "<":
		return e.compareNumbers(fieldValue, compareValue, "<")
	case ">=":
		return e.compareNumbers(fieldValue, compareValue, ">=")
	case "<=":
		return e.compareNumbers(fieldValue, compareValue, "<=")
	case "contains":
		return strings.Contains(fmt.Sprintf("%v", fieldValue), compareValue), nil
	default:
		return false, fmt.Errorf("不支持的操作符: %s", operator)
	}
}

// compareNumbers 比较数字
func (e *WorkflowEngine) compareNumbers(fieldValue interface{}, compareValue string, operator string) (bool, error) {
	var fieldNum float64
	var err error

	switch v := fieldValue.(type) {
	case int:
		fieldNum = float64(v)
	case int64:
		fieldNum = float64(v)
	case float64:
		fieldNum = v
	case string:
		fieldNum, err = strconv.ParseFloat(v, 64)
		if err != nil {
			return false, fmt.Errorf("字段值不是数字: %s", v)
		}
	default:
		return false, fmt.Errorf("字段值类型不支持数字比较: %T", fieldValue)
	}

	compareNum, err := strconv.ParseFloat(compareValue, 64)
	if err != nil {
		return false, fmt.Errorf("比较值不是数字: %s", compareValue)
	}

	switch operator {
	case ">":
		return fieldNum > compareNum, nil
	case "<":
		return fieldNum < compareNum, nil
	case ">=":
		return fieldNum >= compareNum, nil
	case "<=":
		return fieldNum <= compareNum, nil
	default:
		return false, fmt.Errorf("不支持的操作符: %s", operator)
	}
}

// IsEndNode 检查是否是结束节点
func (e *WorkflowEngine) IsEndNode(config *model.WorkflowConfig, nodeID string) bool {
	node := e.GetNodeByID(config, nodeID)
	return node != nil && node.Type == model.NodeTypeEnd
}

// ValidateCondition 验证条件表达式语法
func (e *WorkflowEngine) ValidateCondition(condition string) error {
	if condition == "" {
		return nil
	}

	// 简单的语法检查
	validPattern := regexp.MustCompile(`^[a-zA-Z_]+\.[a-zA-Z_]+\s*(==|!=|>|<|>=|<=|contains)\s*("[^"]*"|'[^']*'|\d+)$`)
	if !validPattern.MatchString(condition) {
		return fmt.Errorf("条件表达式语法错误: %s", condition)
	}

	return nil
}
