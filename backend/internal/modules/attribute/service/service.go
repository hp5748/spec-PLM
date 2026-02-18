package service

import (
	"errors"
	"fmt"
	"strings"

	"plm/internal/modules/attribute/model"
	"plm/internal/modules/attribute/repository"

	"gorm.io/gorm"
)

type Service struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

// VersionRules 版本规则（跳过I、O、Z）
var invalidChars = []string{"I", "O", "Z"}

func isValidVersionChar(char string) bool {
	for _, invalid := range invalidChars {
		if char == invalid {
			return false
		}
	}
	return true
}

func incrementVersionChar(char string) string {
	chars := "ABCDEFGHJKLMNPQRSTUVWXY" // 跳过I、O、Z
	idx := strings.Index(chars, char)
	if idx == -1 || idx == len(chars)-1 {
		return "A"
	}
	return string(chars[idx+1])
}

// incrementVersion 版本号递增
func incrementVersion(version string) string {
	if len(version) != 2 {
		return "AA"
	}
	first := string(version[0])
	second := string(version[1])

	newSecond := incrementVersionChar(second)
	var newFirst string
	if newSecond == "A" {
		newFirst = incrementVersionChar(first)
	} else {
		newFirst = first
	}

	return newFirst + newSecond
}

// CreateAttributeSchema 创建属性Schema
func (s *Service) CreateAttributeSchema(req *model.CreateAttributeSchemaRequest, createdBy uint) (*model.AttributeSchema, error) {
	// 验证实体类型
	if req.EntityType != "MATERIAL" && req.EntityType != "DOCUMENT" {
		return nil, errors.New("无效的实体类型，只能是MATERIAL或DOCUMENT")
	}

	// 验证属性类型
	validAttrTypes := map[string]bool{"main": true, "description": true, "specification": true, "custom": true}
	if !validAttrTypes[req.AttrType] {
		return nil, errors.New("无效的属性类型")
	}

	// 将SchemaConfig转为JSON字符串
	configJSON, err := repository.SchemaConfigToJSON(req.SchemaConfig)
	if err != nil {
		return nil, fmt.Errorf("Schema配置解析失败: %w", err)
	}

	// 获取最新版本号
	latestVersion, err := s.repo.GetLatestVersion(req.EntityType, req.TypeCode, req.AttrType)
	if err != nil {
		return nil, err
	}

	// 计算新版本号
	newVersion := "AA"
	if latestVersion != "" {
		newVersion = incrementVersion(latestVersion)
	}

	// 将同类型的其他版本设为非激活
	s.repo.DeactivateSchemas(req.EntityType, req.TypeCode, req.AttrType)

	schema := &model.AttributeSchema{
		EntityType:   req.EntityType,
		TypeCode:     req.TypeCode,
		AttrType:     req.AttrType,
		SchemaName:   req.SchemaName,
		SchemaConfig: configJSON,
		Version:      newVersion,
		Status:       "DRAFT",
		IsActive:     true,
		CreatedBy:    createdBy,
	}

	if err := s.repo.Create(schema); err != nil {
		return nil, err
	}

	return schema, nil
}

// UpdateAttributeSchema 更新属性Schema（创建新版本）
func (s *Service) UpdateAttributeSchema(id uint, req *model.UpdateAttributeSchemaRequest, updatedBy uint) (*model.AttributeSchema, error) {
	// 获取现有Schema
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 将SchemaConfig转为JSON字符串
	configJSON, err := repository.SchemaConfigToJSON(req.SchemaConfig)
	if err != nil {
		return nil, fmt.Errorf("Schema配置解析失败: %w", err)
	}

	// 如果已发布，则创建新版本
	if existing.Status == "RELEASED" {
		// 计算新版本号
		newVersion := incrementVersion(existing.Version)

		// 将同类型的其他版本设为非激活
		s.repo.DeactivateSchemas(existing.EntityType, existing.TypeCode, existing.AttrType)

		newSchema := &model.AttributeSchema{
			EntityType:   existing.EntityType,
			TypeCode:     existing.TypeCode,
			AttrType:     existing.AttrType,
			SchemaName:   req.SchemaName,
			SchemaConfig: configJSON,
			Version:      newVersion,
			Status:       "DRAFT",
			IsActive:     true,
			CreatedBy:    updatedBy,
		}

		if err := s.repo.Create(newSchema); err != nil {
			return nil, err
		}
		return newSchema, nil
	}

	// 未发布状态，直接更新
	existing.SchemaName = req.SchemaName
	existing.SchemaConfig = configJSON

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}

	return existing, nil
}

// ReleaseAttributeSchema 发布属性Schema
func (s *Service) ReleaseAttributeSchema(id uint) error {
	schema, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	schema.Status = "RELEASED"
	return s.repo.Update(schema)
}

// DeleteAttributeSchema 删除属性Schema
func (s *Service) DeleteAttributeSchema(id uint) error {
	schema, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	// 已发布的不能删除
	if schema.Status == "RELEASED" {
		return errors.New("已发布的Schema不能删除")
	}

	return s.repo.Delete(id)
}

// GetAttributeSchema 获取属性Schema详情
func (s *Service) GetAttributeSchema(id uint) (*model.AttributeSchemaResponse, error) {
	schema, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	config, err := repository.ParseSchemaConfig(schema.SchemaConfig)
	if err != nil {
		return nil, err
	}

	return &model.AttributeSchemaResponse{
		AttributeSchema: schema,
		SchemaConfig:    config,
	}, nil
}

// GetActiveSchema 获取当前激活的Schema
func (s *Service) GetActiveSchema(entityType, typeCode, attrType string) (*model.AttributeSchemaResponse, error) {
	schema, err := s.repo.GetActiveSchema(entityType, typeCode, attrType)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // 没有配置Schema时返回nil
		}
		return nil, err
	}

	config, err := repository.ParseSchemaConfig(schema.SchemaConfig)
	if err != nil {
		return nil, err
	}

	return &model.AttributeSchemaResponse{
		AttributeSchema: schema,
		SchemaConfig:    config,
	}, nil
}

// GetAllActiveSchemas 获取指定类型的所有激活Schema
func (s *Service) GetAllActiveSchemas(entityType, typeCode string) ([]model.AttributeSchemaResponse, error) {
	schemas, err := s.repo.GetAllActiveSchemasByEntity(entityType, typeCode)
	if err != nil {
		return nil, err
	}

	responses := make([]model.AttributeSchemaResponse, 0, len(schemas))
	for _, schema := range schemas {
		config, err := repository.ParseSchemaConfig(schema.SchemaConfig)
		if err != nil {
			continue
		}
		responses = append(responses, model.AttributeSchemaResponse{
			AttributeSchema: &schema,
			SchemaConfig:    config,
		})
	}

	return responses, nil
}

// ListAttributeSchemas 分页查询属性Schema列表
func (s *Service) ListAttributeSchemas(query *model.AttributeSchemaQuery) ([]model.AttributeSchemaResponse, int64, error) {
	schemas, total, err := s.repo.List(query)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]model.AttributeSchemaResponse, 0, len(schemas))
	for _, schema := range schemas {
		config, err := repository.ParseSchemaConfig(schema.SchemaConfig)
		if err != nil {
			config = model.SchemaConfig{}
		}
		responses = append(responses, model.AttributeSchemaResponse{
			AttributeSchema: &schema,
			SchemaConfig:    config,
		})
	}

	return responses, total, nil
}

// ActivateSchema 激活指定版本的Schema
func (s *Service) ActivateSchema(id uint) error {
	schema, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	// 必须是已发布状态才能激活
	if schema.Status != "RELEASED" {
		return errors.New("只有已发布的Schema才能激活")
	}

	// 将同类型的其他版本设为非激活
	if err := s.repo.DeactivateSchemas(schema.EntityType, schema.TypeCode, schema.AttrType); err != nil {
		return err
	}

	// 激活当前版本
	schema.IsActive = true
	return s.repo.Update(schema)
}
