package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"plm/internal/modules/script/model"
	"plm/internal/modules/script/repository"

	"github.com/dop251/goja"
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

// ==================== 脚本管理 ====================

// CreateScript 创建脚本
func (s *Service) CreateScript(ctx context.Context, req *model.CreateScriptRequest, createdBy uint) (*model.Script, error) {
	// 检查编码是否已存在
	exists, err := s.repo.ExistsScriptByCode(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("脚本编码已存在")
	}

	// 验证脚本内容
	if err := s.ValidateScript(req.Type, req.Content); err != nil {
		return nil, fmt.Errorf("脚本验证失败: %v", err)
	}

	timeout := req.Timeout
	if timeout <= 0 {
		timeout = 5000 // 默认5秒超时
	}

	script := &model.Script{
		Name:         req.Name,
		Code:         req.Code,
		Type:         req.Type,
		TriggerType:  req.TriggerType,
		BusinessType: req.BusinessType,
		Content:      req.Content,
		Description:  req.Description,
		Status:       "DRAFT",
		Timeout:      timeout,
		CreatedBy:    createdBy,
	}

	if err := s.repo.CreateScript(ctx, script); err != nil {
		return nil, err
	}

	return s.repo.GetScriptByID(ctx, script.ID)
}

// GetScript 获取脚本详情
func (s *Service) GetScript(ctx context.Context, id uint) (*model.Script, error) {
	return s.repo.GetScriptByID(ctx, id)
}

// UpdateScript 更新脚本
func (s *Service) UpdateScript(ctx context.Context, id uint, req *model.UpdateScriptRequest) (*model.Script, error) {
	script, err := s.repo.GetScriptByID(ctx, id)
	if err != nil {
		return nil, errors.New("脚本不存在")
	}

	// 已发布的脚本不能修改
	if script.Status == "RELEASED" {
		return nil, errors.New("已发布的脚本不能修改")
	}

	if req.Name != "" {
		script.Name = req.Name
	}
	if req.Content != "" {
		// 验证脚本内容
		if err := s.ValidateScript(script.Type, req.Content); err != nil {
			return nil, fmt.Errorf("脚本验证失败: %v", err)
		}
		script.Content = req.Content
	}
	if req.Description != "" {
		script.Description = req.Description
	}
	if req.Timeout > 0 {
		script.Timeout = req.Timeout
	}

	if err := s.repo.UpdateScript(ctx, script); err != nil {
		return nil, err
	}

	return s.repo.GetScriptByID(ctx, id)
}

// ReleaseScript 发布脚本
func (s *Service) ReleaseScript(ctx context.Context, id uint) (*model.Script, error) {
	script, err := s.repo.GetScriptByID(ctx, id)
	if err != nil {
		return nil, errors.New("脚本不存在")
	}

	if script.Status == "RELEASED" {
		return nil, errors.New("脚本已发布")
	}

	script.Status = "RELEASED"
	if err := s.repo.UpdateScript(ctx, script); err != nil {
		return nil, err
	}

	return s.repo.GetScriptByID(ctx, id)
}

// ObsoleteScript 废弃脚本
func (s *Service) ObsoleteScript(ctx context.Context, id uint) (*model.Script, error) {
	script, err := s.repo.GetScriptByID(ctx, id)
	if err != nil {
		return nil, errors.New("脚本不存在")
	}

	script.Status = "OBSOLETE"
	if err := s.repo.UpdateScript(ctx, script); err != nil {
		return nil, err
	}

	return s.repo.GetScriptByID(ctx, id)
}

// DeleteScript 删除脚本
func (s *Service) DeleteScript(ctx context.Context, id uint) error {
	script, err := s.repo.GetScriptByID(ctx, id)
	if err != nil {
		return errors.New("脚本不存在")
	}

	if script.Status == "RELEASED" {
		return errors.New("已发布的脚本不能删除")
	}

	return s.repo.DeleteScript(ctx, id)
}

// ListScripts 脚本列表
func (s *Service) ListScripts(ctx context.Context, query *model.ScriptListQuery) ([]*model.Script, int64, error) {
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 20
	}

	return s.repo.ListScripts(ctx, query)
}

// ==================== 脚本执行 ====================

// ExecuteScript 执行脚本
func (s *Service) ExecuteScript(ctx context.Context, scriptID uint, businessType string, businessID uint, execContext map[string]interface{}, executedBy uint) (*model.ExecuteScriptResult, error) {
	script, err := s.repo.GetScriptByID(ctx, scriptID)
	if err != nil {
		return nil, errors.New("脚本不存在")
	}

	if script.Status != "RELEASED" {
		return nil, errors.New("脚本未发布，不能执行")
	}

	result, log := s.doExecute(script, businessType, businessID, execContext, executedBy)

	// 保存执行日志
	if err := s.repo.CreateExecutionLog(ctx, log); err != nil {
		// 记录失败不影响返回
	}

	return result, nil
}

// ExecuteScriptsByTrigger 根据触发类型执行脚本
func (s *Service) ExecuteScriptsByTrigger(ctx context.Context, triggerType, businessType string, businessID uint, execContext map[string]interface{}, executedBy uint) ([]*model.ExecuteScriptResult, error) {
	scripts, err := s.repo.GetScriptsByTrigger(ctx, triggerType, businessType)
	if err != nil {
		return nil, err
	}

	var results []*model.ExecuteScriptResult
	for _, script := range scripts {
		result, log := s.doExecute(script, businessType, businessID, execContext, executedBy)

		// 保存执行日志
		s.repo.CreateExecutionLog(ctx, log)

		results = append(results, result)

		// 如果脚本执行失败且是BEFORE类型，中断执行
		if !result.Success && strings.HasPrefix(triggerType, "BEFORE_") {
			return results, fmt.Errorf("脚本[%s]执行失败: %s", script.Name, result.Error)
		}
	}

	return results, nil
}

// doExecute 实际执行脚本
func (s *Service) doExecute(script *model.Script, businessType string, businessID uint, execContext map[string]interface{}, executedBy uint) (*model.ExecuteScriptResult, *model.ScriptExecutionLog) {
	startTime := time.Now()

	result := &model.ExecuteScriptResult{
		Success: false,
		Output:  make(map[string]interface{}),
	}

	log := &model.ScriptExecutionLog{
		ScriptID:     script.ID,
		ScriptName:   script.Name,
		BusinessType: businessType,
		BusinessID:   businessID,
		TriggerType:  script.TriggerType,
		ExecutedBy:   executedBy,
		ExecutedAt:   startTime,
	}

	// 序列化输入
	if inputBytes, err := json.Marshal(execContext); err == nil {
		log.Input = string(inputBytes)
	}

	var output map[string]interface{}
	var execErr error

	switch script.Type {
	case "JAVASCRIPT":
		output, execErr = s.executeJavaScript(script, businessType, businessID, execContext)
	case "SQL":
		output, execErr = s.executeSQL(script, businessType, businessID, execContext)
	default:
		execErr = fmt.Errorf("不支持的脚本类型: %s", script.Type)
	}

	log.Duration = int(time.Since(startTime).Milliseconds())

	if execErr != nil {
		result.Error = execErr.Error()
		log.ErrorMsg = execErr.Error()
		log.Success = false
	} else {
		result.Success = true
		result.Output = output
		log.Success = true
		if outputBytes, err := json.Marshal(output); err == nil {
			log.Output = string(outputBytes)
		}
	}

	result.LogID = log.ID
	return result, log
}

// ==================== JavaScript执行引擎 ====================

// executeJavaScript 执行JavaScript脚本
func (s *Service) executeJavaScript(script *model.Script, businessType string, businessID uint, context map[string]interface{}) (map[string]interface{}, error) {
	vm := goja.New()

	// 设置超时
	timeout := time.Duration(script.Timeout) * time.Millisecond
	done := make(chan struct{})
	var execErr error
	var result map[string]interface{}

	go func() {
		defer close(done)
		defer func() {
			if r := recover(); r != nil {
				execErr = fmt.Errorf("JavaScript执行panic: %v", r)
			}
		}()

		// 注入上下文变量
		vm.Set("businessType", businessType)
		vm.Set("businessID", businessID)
		vm.Set("context", context)

		// 注入数据库查询API
		s.injectDatabaseAPI(vm)

		// 注入工具函数
		s.injectUtilityAPI(vm)

		// 注入日志API
		var logs []string
		vm.Set("log", func(call goja.FunctionCall) goja.Value {
			var args []string
			for _, arg := range call.Arguments {
				args = append(args, arg.String())
			}
			logs = append(logs, strings.Join(args, " "))
			return vm.ToValue(nil)
		})

		// 包装脚本，确保返回结果
		wrappedScript := fmt.Sprintf(`
			(function() {
				%s
				if (typeof main === 'function') {
					return main();
				}
				return { success: true, message: '脚本执行完成' };
			})()
		`, script.Content)

		value, err := vm.RunString(wrappedScript)
		if err != nil {
			execErr = fmt.Errorf("JavaScript执行错误: %v", err)
			return
		}

		// 解析返回值
		if value.Export() != nil {
			switch v := value.Export().(type) {
			case map[string]interface{}:
				result = v
			case string:
				result = map[string]interface{}{"message": v}
			default:
				result = map[string]interface{}{"result": v}
			}
		} else {
			result = map[string]interface{}{"success": true}
		}
	}()

	select {
	case <-done:
		return result, execErr
	case <-time.After(timeout):
		vm.Interrupt("执行超时")
		return nil, fmt.Errorf("脚本执行超时（%dms）", script.Timeout)
	}
}

// injectDatabaseAPI 注入数据库查询API
func (s *Service) injectDatabaseAPI(vm *goja.Goja) {
	// dbQuery - 执行SQL查询
	vm.Set("dbQuery", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return vm.ToValue(map[string]interface{}{"error": "缺少SQL参数"})
		}

		sql := call.Arguments[0].String()

		// 安全检查：只允许SELECT语句
		if !s.isSafeQuery(sql) {
			return vm.ToValue(map[string]interface{}{"error": "只允许执行SELECT查询"})
		}

		var results []map[string]interface{}
		if err := s.db.Raw(sql).Scan(&results).Error; err != nil {
			return vm.ToValue(map[string]interface{}{"error": err.Error()})
		}

		return vm.ToValue(results)
	})

	// dbQueryOne - 查询单条记录
	vm.Set("dbQueryOne", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return vm.ToValue(map[string]interface{}{"error": "缺少SQL参数"})
		}

		sql := call.Arguments[0].String()

		if !s.isSafeQuery(sql) {
			return vm.ToValue(map[string]interface{}{"error": "只允许执行SELECT查询"})
		}

		var result map[string]interface{}
		if err := s.db.Raw(sql).Scan(&result).Error; err != nil {
			return vm.ToValue(map[string]interface{}{"error": err.Error()})
		}

		return vm.ToValue(result)
	})

	// getMaterial - 获取物料信息
	vm.Set("getMaterial", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return vm.ToValue(nil)
		}

		id := call.Arguments[0].ToInteger()
		var material map[string]interface{}
		if err := s.db.Table("materials").Where("id = ?", id).First(&material).Error; err != nil {
			return vm.ToValue(nil)
		}

		return vm.ToValue(material)
	})

	// getDocument - 获取文档信息
	vm.Set("getDocument", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return vm.ToValue(nil)
		}

		id := call.Arguments[0].ToInteger()
		var document map[string]interface{}
		if err := s.db.Table("documents").Where("id = ?", id).First(&document).Error; err != nil {
			return vm.ToValue(nil)
		}

		return vm.ToValue(document)
	})

	// getBOM - 获取BOM信息
	vm.Set("getBOM", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return vm.ToValue(nil)
		}

		id := call.Arguments[0].ToInteger()
		var bom map[string]interface{}
		if err := s.db.Table("bom_views").Where("id = ?", id).First(&bom).Error; err != nil {
			return vm.ToValue(nil)
		}

		return vm.ToValue(bom)
	})
}

// injectUtilityAPI 注入工具函数API
func (s *Service) injectUtilityAPI(vm *goja.Goja) {
	// JSON解析
	vm.Set("parseJSON", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return vm.ToValue(nil)
		}

		var result interface{}
		if err := json.Unmarshal([]byte(call.Arguments[0].String()), &result); err != nil {
			return vm.ToValue(nil)
		}

		return vm.ToValue(result)
	})

	// JSON序列化
	vm.Set("stringifyJSON", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return vm.ToValue("")
		}

		bytes, err := json.Marshal(call.Arguments[0].Export())
		if err != nil {
			return vm.ToValue("")
		}

		return vm.ToValue(string(bytes))
	})
}

// ==================== SQL执行引擎 ====================

// executeSQL 执行SQL脚本
func (s *Service) executeSQL(script *model.Script, businessType string, businessID uint, context map[string]interface{}) (map[string]interface{}, error) {
	// 安全检查：只允许SELECT语句
	if !s.isSafeQuery(script.Content) {
		return nil, errors.New("只允许执行SELECT查询")
	}

	// 替换脚本中的变量
	sql := s.replaceVariables(script.Content, businessType, businessID, context)

	var results []map[string]interface{}
	if err := s.db.Raw(sql).Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("SQL执行错误: %v", err)
	}

	return map[string]interface{}{
		"success": true,
		"data":    results,
		"count":   len(results),
	}, nil
}

// isSafeQuery 检查是否是安全的查询（只允许SELECT）
func (s *Service) isSafeQuery(sql string) bool {
	// 移除注释和多余空格
	sql = regexp.MustCompile(`--.*$`).ReplaceAllString(sql, "")
	sql = regexp.MustCompile(`/\\*.*\\*/`).ReplaceAllString(sql, "")
	sql = strings.TrimSpace(strings.ToUpper(sql))

	// 只允许SELECT开头
	return strings.HasPrefix(sql, "SELECT")
}

// replaceVariables 替换SQL中的变量
func (s *Service) replaceVariables(sql string, businessType string, businessID uint, context map[string]interface{}) string {
	// 替换内置变量
	sql = strings.ReplaceAll(sql, "${businessType}", businessType)
	sql = strings.ReplaceAll(sql, "${businessID}", fmt.Sprintf("%d", businessID))

	// 替换上下文变量
	for key, value := range context {
		placeholder := fmt.Sprintf("${%s}", key)
		var strValue string
		switch v := value.(type) {
		case string:
			strValue = fmt.Sprintf("'%s'", v)
		case int, int64, float64:
			strValue = fmt.Sprintf("%v", v)
		default:
			bytes, _ := json.Marshal(v)
			strValue = fmt.Sprintf("'%s'", string(bytes))
		}
		sql = strings.ReplaceAll(sql, placeholder, strValue)
	}

	return sql
}

// ==================== 脚本验证 ====================

// ValidateScript 验证脚本内容
func (s *Service) ValidateScript(scriptType, content string) error {
	switch scriptType {
	case "JAVASCRIPT":
		return s.validateJavaScript(content)
	case "SQL":
		return s.validateSQL(content)
	default:
		return fmt.Errorf("不支持的脚本类型: %s", scriptType)
	}
}

// validateJavaScript 验证JavaScript脚本
func (s *Service) validateJavaScript(content string) error {
	vm := goja.New()

	// 禁止危险操作
	dangerousPatterns := []string{
		"require(", "import ", "eval(", "Function(",
		"process.", "global.", "module.exports",
		"__dirname", "__filename",
	}

	for _, pattern := range dangerousPatterns {
		if strings.Contains(content, pattern) {
			return fmt.Errorf("脚本包含禁止的操作: %s", pattern)
		}
	}

	// 尝试编译脚本
	_, err := vm.RunString(fmt.Sprintf("(function() { %s })()", content))
	if err != nil {
		return fmt.Errorf("脚本语法错误: %v", err)
	}

	return nil
}

// validateSQL 验证SQL脚本
func (s *Service) validateSQL(content string) error {
	// 检查是否是SELECT语句
	if !s.isSafeQuery(content) {
		return errors.New("只允许SELECT查询")
	}

	return nil
}

// ==================== 测试功能 ====================

// TestScript 测试脚本
func (s *Service) TestScript(ctx context.Context, req *model.ScriptTestRequest) (*model.ScriptTestResult, error) {
	result := &model.ScriptTestResult{
		Success: false,
		Output:  make(map[string]interface{}),
		Logs:    []string{},
	}

	// 验证脚本
	if err := s.ValidateScript(req.Type, req.Content); err != nil {
		result.Error = err.Error()
		return result, nil
	}

	startTime := time.Now()

	var output map[string]interface{}
	var execErr error

	switch req.Type {
	case "JAVASCRIPT":
		output, execErr = s.executeJavaScript(&model.Script{
			Content: req.Content,
			Timeout: 5000,
		}, req.BusinessType, req.BusinessID, req.Context)
	case "SQL":
		output, execErr = s.executeSQL(&model.Script{
			Content: req.Content,
		}, req.BusinessType, req.BusinessID, req.Context)
	}

	result.Duration = int(time.Since(startTime).Milliseconds())

	if execErr != nil {
		result.Error = execErr.Error()
	} else {
		result.Success = true
		result.Output = output
	}

	return result, nil
}

// ==================== 执行日志 ====================

// GetExecutionLogs 获取执行日志列表
func (s *Service) GetExecutionLogs(ctx context.Context, scriptID uint, page, pageSize int) ([]*model.ScriptExecutionLog, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	return s.repo.GetExecutionLogs(ctx, scriptID, page, pageSize)
}

// GetExecutionLog 获取执行日志详情
func (s *Service) GetExecutionLog(ctx context.Context, id uint) (*model.ScriptExecutionLog, error) {
	return s.repo.GetExecutionLogByID(ctx, id)
}
