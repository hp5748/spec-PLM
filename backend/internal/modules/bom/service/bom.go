package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"

	"plm/internal/modules/bom/model"
	"plm/internal/modules/bom/repository"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type Service struct {
	repo *repository.Repository
	db   *gorm.DB
}

func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

// CreateBOMView 创建BOM视图
func (s *Service) CreateBOMView(ctx context.Context, req *model.CreateBOMViewRequest, createdBy uint) (*model.BOMView, error) {
	// 检查根物料是否已有BOM视图
	exists, err := s.repo.ExistsRootMaterial(ctx, req.RootMaterialID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("该物料已存在BOM视图")
	}

	// 精确BOM必须指定版本
	if req.IsExact && req.RootVersion == "" {
		return nil, errors.New("精确BOM必须指定根物料版本")
	}

	bom := &model.BOMView{
		Name:           req.Name,
		RootMaterialID: req.RootMaterialID,
		RootVersion:    req.RootVersion,
		IsExact:        req.IsExact,
		Status:         "DRAFT",
		Description:    req.Description,
		CreatedBy:      createdBy,
	}

	if err := s.repo.CreateBOMView(ctx, bom); err != nil {
		return nil, err
	}

	return s.repo.GetBOMViewByID(ctx, bom.ID)
}

// GetBOMView 获取BOM视图详情
func (s *Service) GetBOMView(ctx context.Context, id uint) (*model.BOMViewResponse, error) {
	bom, err := s.repo.GetBOMViewByID(ctx, id)
	if err != nil {
		return nil, errors.New("BOM视图不存在")
	}

	count, _ := s.repo.CountBOMItemsByViewID(ctx, id)

	return &model.BOMViewResponse{
		BOMView:   *bom,
		ItemCount: int(count),
	}, nil
}

// UpdateBOMView 更新BOM视图
func (s *Service) UpdateBOMView(ctx context.Context, id uint, req *model.UpdateBOMViewRequest) (*model.BOMView, error) {
	bom, err := s.repo.GetBOMViewByID(ctx, id)
	if err != nil {
		return nil, errors.New("BOM视图不存在")
	}

	// 检查状态是否允许编辑
	if bom.Status == "RELEASED" {
		return nil, errors.New("已发布的BOM不能编辑")
	}
	if bom.Status == "REVIEWING" {
		return nil, errors.New("审核中的BOM不能编辑")
	}

	if req.Name != "" {
		bom.Name = req.Name
	}
	if req.Status != "" {
		bom.Status = req.Status
	}
	if req.Description != "" {
		bom.Description = req.Description
	}

	if err := s.repo.UpdateBOMView(ctx, bom); err != nil {
		return nil, err
	}

	return s.repo.GetBOMViewByID(ctx, id)
}

// DeleteBOMView 删除BOM视图
func (s *Service) DeleteBOMView(ctx context.Context, id uint) error {
	bom, err := s.repo.GetBOMViewByID(ctx, id)
	if err != nil {
		return errors.New("BOM视图不存在")
	}

	// 检查状态是否允许删除
	if bom.Status == "RELEASED" {
		return errors.New("已发布的BOM不能删除")
	}
	if bom.Status == "REVIEWING" {
		return errors.New("审核中的BOM不能删除")
	}

	// 先删除所有BOM项
	if err := s.repo.DeleteBOMItemsByViewID(ctx, id); err != nil {
		return err
	}

	return s.repo.DeleteBOMView(ctx, id)
}

// ListBOMViews BOM视图列表
func (s *Service) ListBOMViews(ctx context.Context, query *model.BOMListQuery) ([]*model.BOMViewResponse, int64, error) {
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 20
	}

	list, total, err := s.repo.ListBOMViews(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	var result []*model.BOMViewResponse
	for _, bom := range list {
		count, _ := s.repo.CountBOMItemsByViewID(ctx, bom.ID)
		result = append(result, &model.BOMViewResponse{
			BOMView:   *bom,
			ItemCount: int(count),
		})
	}

	return result, total, nil
}

// AddBOMItem 添加BOM项
func (s *Service) AddBOMItem(ctx context.Context, bomViewID uint, req *model.AddBOMItemRequest) (*model.BOMItem, error) {
	bom, err := s.repo.GetBOMViewByID(ctx, bomViewID)
	if err != nil {
		return nil, errors.New("BOM视图不存在")
	}

	// 检查状态是否允许编辑
	if bom.Status == "RELEASED" {
		return nil, errors.New("已发布的BOM不能编辑")
	}
	if bom.Status == "REVIEWING" {
		return nil, errors.New("审核中的BOM不能编辑")
	}

	// 不能将根物料作为子项添加
	if req.MaterialID == bom.RootMaterialID {
		return nil, errors.New("不能将根物料作为子项添加")
	}

	// 检查物料是否已存在于BOM中
	existingItem, err := s.repo.GetBOMItemByMaterialInBOM(ctx, bomViewID, req.MaterialID)
	if err == nil && existingItem != nil {
		return nil, errors.New("该物料已存在于BOM中")
	}

	// 精确BOM必须指定版本
	if bom.IsExact && req.Version == "" {
		return nil, errors.New("精确BOM必须指定物料版本")
	}

	// 计算层级
	level := 1
	if req.ParentID != nil {
		parentItem, err := s.repo.GetBOMItemByID(ctx, *req.ParentID)
		if err != nil {
			return nil, errors.New("父节点不存在")
		}
		level = parentItem.Level + 1
	}

	// 检查循环依赖
	if err := s.checkCircularDependency(ctx, bomViewID, req.ParentID, req.MaterialID); err != nil {
		return nil, err
	}

	item := &model.BOMItem{
		BOMViewID:  bomViewID,
		ParentID:   req.ParentID,
		MaterialID: req.MaterialID,
		Version:    req.Version,
		Quantity:   req.Quantity,
		SortOrder:  req.SortOrder,
		Level:      level,
	}

	if err := s.repo.CreateBOMItem(ctx, item); err != nil {
		return nil, err
	}

	return s.repo.GetBOMItemByID(ctx, item.ID)
}

// UpdateBOMItem 更新BOM项
func (s *Service) UpdateBOMItem(ctx context.Context, bomViewID, itemID uint, req *model.UpdateBOMItemRequest) (*model.BOMItem, error) {
	bom, err := s.repo.GetBOMViewByID(ctx, bomViewID)
	if err != nil {
		return nil, errors.New("BOM视图不存在")
	}

	// 检查状态是否允许编辑
	if bom.Status == "RELEASED" {
		return nil, errors.New("已发布的BOM不能编辑")
	}
	if bom.Status == "REVIEWING" {
		return nil, errors.New("审核中的BOM不能编辑")
	}

	item, err := s.repo.GetBOMItemByID(ctx, itemID)
	if err != nil {
		return nil, errors.New("BOM项不存在")
	}

	if item.BOMViewID != bomViewID {
		return nil, errors.New("BOM项不属于该BOM视图")
	}

	item.Quantity = req.Quantity
	item.SortOrder = req.SortOrder

	if err := s.repo.UpdateBOMItem(ctx, item); err != nil {
		return nil, err
	}

	return s.repo.GetBOMItemByID(ctx, itemID)
}

// DeleteBOMItem 删除BOM项
func (s *Service) DeleteBOMItem(ctx context.Context, bomViewID, itemID uint) error {
	bom, err := s.repo.GetBOMViewByID(ctx, bomViewID)
	if err != nil {
		return errors.New("BOM视图不存在")
	}

	// 检查状态是否允许编辑
	if bom.Status == "RELEASED" {
		return errors.New("已发布的BOM不能编辑")
	}
	if bom.Status == "REVIEWING" {
		return errors.New("审核中的BOM不能编辑")
	}

	item, err := s.repo.GetBOMItemByID(ctx, itemID)
	if err != nil {
		return errors.New("BOM项不存在")
	}

	if item.BOMViewID != bomViewID {
		return errors.New("BOM项不属于该BOM视图")
	}

	// 递归删除子项
	if err := s.deleteBOMItemRecursive(ctx, itemID); err != nil {
		return err
	}

	return nil
}

// deleteBOMItemRecursive 递归删除BOM项及其子项
func (s *Service) deleteBOMItemRecursive(ctx context.Context, itemID uint) error {
	// 获取所有子项
	var children []model.BOMItem
	if s.db != nil {
		if err := s.db.WithContext(ctx).Where("parent_id = ?", itemID).Find(&children).Error; err != nil {
			return err
		}
	}

	// 递归删除子项
	for _, child := range children {
		if err := s.deleteBOMItemRecursive(ctx, child.ID); err != nil {
			return err
		}
	}

	// 删除当前项
	return s.repo.DeleteBOMItem(ctx, itemID)
}

// GetBOMTree 获取BOM树形结构
func (s *Service) GetBOMTree(ctx context.Context, bomViewID uint) (*model.BOMTreeNode, error) {
	bom, err := s.repo.GetBOMViewByID(ctx, bomViewID)
	if err != nil {
		return nil, errors.New("BOM视图不存在")
	}

	items, err := s.repo.GetBOMItemsByViewID(ctx, bomViewID)
	if err != nil {
		return nil, err
	}

	// 构建根节点
	rootNode := &model.BOMTreeNode{
		ID:         0,
		MaterialID: bom.RootMaterialID,
		ItemID:     bom.RootMaterial.ItemID,
		ItemName:   bom.RootMaterial.ItemName,
		Version:    bom.RootVersion,
		Quantity:   1,
		Level:      0,
		Unit:       bom.RootMaterial.Unit,
		Children:   []model.BOMTreeNode{},
	}

	// 构建树形结构
	rootNode.Children = s.buildTree(items, nil)

	return rootNode, nil
}

// buildTree 递归构建树形结构
func (s *Service) buildTree(items []model.BOMItem, parentID *uint) []model.BOMTreeNode {
	var nodes []model.BOMTreeNode

	for _, item := range items {
		isChildOfParent := (parentID == nil && item.ParentID == nil) ||
			(parentID != nil && item.ParentID != nil && *item.ParentID == *parentID)

		if isChildOfParent {
			node := model.BOMTreeNode{
				ID:         item.ID,
				MaterialID: item.MaterialID,
				ItemID:     getItemID(item.Material),
				ItemName:   getItemName(item.Material),
				Version:    item.Version,
				Quantity:   item.Quantity,
				SortOrder:  item.SortOrder,
				Level:      item.Level,
				Unit:       getUnit(item.Material),
				Children:   s.buildTree(items, &item.ID),
			}
			nodes = append(nodes, node)
		}
	}

	// 按sort_order排序
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].SortOrder < nodes[j].SortOrder
	})

	return nodes
}

// checkCircularDependency 检查循环依赖
func (s *Service) checkCircularDependency(ctx context.Context, bomViewID uint, parentID *uint, childMaterialID uint) error {
	// 获取BOM视图
	bom, err := s.repo.GetBOMViewByID(ctx, bomViewID)
	if err != nil {
		return err
	}

	// 如果要添加的物料是根物料，则存在循环
	if childMaterialID == bom.RootMaterialID {
		return errors.New("不能添加根物料作为子项，将形成循环依赖")
	}

	// 如果没有父节点，只需要检查根物料
	if parentID == nil {
		return nil
	}

	// 获取所有BOM项
	items, err := s.repo.GetBOMItemsByViewID(ctx, bomViewID)
	if err != nil {
		return err
	}

	// 构建物料ID到父物料ID的映射
	materialParentMap := make(map[uint]uint)
	for _, item := range items {
		if item.ParentID != nil {
			// 找到父节点的物料ID
			for _, p := range items {
				if p.ID == *item.ParentID {
					materialParentMap[item.MaterialID] = p.MaterialID
					break
				}
			}
		}
	}

	// 获取父节点的物料ID
	parentItem, err := s.repo.GetBOMItemByID(ctx, *parentID)
	if err != nil {
		return nil
	}

	// 检查childMaterialID是否是parentItem.MaterialID的祖先
	visited := make(map[uint]bool)
	currentMaterialID := parentItem.MaterialID

	for currentMaterialID != 0 && !visited[currentMaterialID] {
		if currentMaterialID == childMaterialID {
			return errors.New("添加此物料将形成循环依赖")
		}
		visited[currentMaterialID] = true
		currentMaterialID = materialParentMap[currentMaterialID]
	}

	return nil
}

// ConvertBOMType 转换BOM类型
func (s *Service) ConvertBOMType(ctx context.Context, bomViewID uint, toExact bool) error {
	bom, err := s.repo.GetBOMViewByID(ctx, bomViewID)
	if err != nil {
		return errors.New("BOM视图不存在")
	}

	// 检查状态是否允许编辑
	if bom.Status == "RELEASED" {
		return errors.New("已发布的BOM不能编辑")
	}
	if bom.Status == "REVIEWING" {
		return errors.New("审核中的BOM不能编辑")
	}

	// 如果类型相同，无需转换
	if bom.IsExact == toExact {
		return nil
	}

	// 转为精确BOM时，检查所有项是否都有版本
	if toExact {
		items, err := s.repo.GetBOMItemsByViewID(ctx, bomViewID)
		if err != nil {
			return err
		}

		for _, item := range items {
			if item.Version == "" {
				return fmt.Errorf("物料 %s 没有指定版本，无法转换为精确BOM", getItemID(item.Material))
			}
		}

		if bom.RootVersion == "" {
			return errors.New("根物料没有指定版本，无法转换为精确BOM")
		}
	}

	// 更新BOM视图类型
	bom.IsExact = toExact
	if err := s.repo.UpdateBOMView(ctx, bom); err != nil {
		return err
	}

	// 如果转为非精确BOM，清空所有版本字段
	if !toExact {
		items, _ := s.repo.GetBOMItemsByViewID(ctx, bomViewID)
		for _, item := range items {
			item.Version = ""
			s.repo.UpdateBOMItem(ctx, &item)
		}
	}

	return nil
}

// 辅助函数
func getItemID(m *model.Material) string {
	if m == nil {
		return ""
	}
	return m.ItemID
}

func getItemName(m *model.Material) string {
	if m == nil {
		return ""
	}
	return m.ItemName
}

func getUnit(m *model.Material) string {
	if m == nil {
		return ""
	}
	return m.Unit
}

// SetDB 设置数据库连接（用于事务）
func (s *Service) SetDB(db *gorm.DB) {
	s.repo = repository.NewRepository(db)
}

// ExportBOM 导出BOM为Excel
func (s *Service) ExportBOM(ctx context.Context, bomViewID uint) (*bytes.Buffer, string, error) {
	// 获取BOM树形结构
	tree, err := s.GetBOMTree(ctx, bomViewID)
	if err != nil {
		return nil, "", err
	}

	// 创建Excel文件
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "BOM结构"
	// 设置表头
	headers := []string{"层级", "物料编码", "物料名称", "版本", "数量", "单位", "排序"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	// 设置表头样式
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#CCCCCC"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	f.SetCellStyle(sheetName, "A1", "G1", headerStyle)

	// 递归写入BOM数据
	rowNum := 2
	var writeNode func(node *model.BOMTreeNode, level int)
	writeNode = func(node *model.BOMTreeNode, level int) {
		// 写入当前节点
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), level)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), node.ItemID)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), node.ItemName)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), node.Version)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), node.Quantity)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowNum), node.Unit)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", rowNum), node.SortOrder)
		rowNum++

		// 递归写入子节点
		for _, child := range node.Children {
			writeNode(&child, level+1)
		}
	}

	// 从根节点开始写入
	writeNode(tree, 0)

	// 设置列宽
	f.SetColWidth(sheetName, "A", "A", 8)
	f.SetColWidth(sheetName, "B", "B", 20)
	f.SetColWidth(sheetName, "C", "C", 30)
	f.SetColWidth(sheetName, "D", "D", 8)
	f.SetColWidth(sheetName, "E", "E", 10)
	f.SetColWidth(sheetName, "F", "F", 10)
	f.SetColWidth(sheetName, "G", "G", 8)

	// 写入buffer
	buf := new(bytes.Buffer)
	if err := f.Write(buf); err != nil {
		return nil, "", err
	}

	filename := fmt.Sprintf("BOM_%s.xlsx", tree.ItemID)
	return buf, filename, nil
}

// ImportBOM 从Excel导入BOM
func (s *Service) ImportBOM(ctx context.Context, bomViewID uint, fileData []byte, createdBy uint) (*model.ImportBOMResult, error) {
	result := &model.ImportBOMResult{}

	// 获取BOM视图
	bom, err := s.repo.GetBOMViewByID(ctx, bomViewID)
	if err != nil {
		return nil, errors.New("BOM视图不存在")
	}

	// 检查状态是否允许编辑
	if bom.Status == "RELEASED" {
		return nil, errors.New("已发布的BOM不能编辑")
	}
	if bom.Status == "REVIEWING" {
		return nil, errors.New("审核中的BOM不能编辑")
	}

	// 打开Excel文件
	f, err := excelize.OpenReader(bytes.NewReader(fileData))
	if err != nil {
		return nil, fmt.Errorf("无法读取Excel文件: %v", err)
	}
	defer f.Close()

	// 获取所有行
	rows, err := f.GetRows("BOM结构")
	if err != nil {
		// 尝试获取第一个工作表
		sheets := f.GetSheetList()
		if len(sheets) > 0 {
			rows, err = f.GetRows(sheets[0])
			if err != nil {
				return nil, fmt.Errorf("无法读取Excel数据: %v", err)
			}
		} else {
			return nil, errors.New("Excel文件没有工作表")
		}
	}

	// 跳过表头，从第二行开始
	for i, row := range rows {
		if i == 0 {
			continue // 跳过表头
		}
		if len(row) < 5 {
			result.FailCount++
			result.Errors = append(result.Errors, fmt.Sprintf("第%d行数据不完整", i+1))
			continue
		}

		// 解析行数据
		level := 0
		if row[0] != "" {
			level, _ = strconv.Atoi(row[0])
		}
		itemID := row[1]
		itemName := row[2]
		version := row[3]
		quantity := 1.0
		if row[4] != "" {
			quantity, _ = strconv.ParseFloat(row[4], 64)
		}

		if itemID == "" {
			result.FailCount++
			result.Errors = append(result.Errors, fmt.Sprintf("第%d行物料编码为空", i+1))
			continue
		}

		// 查找或创建物料
		var material *model.Material
		if s.db != nil {
			material = &model.Material{}
			if err := s.db.WithContext(ctx).Where("item_id = ?", itemID).First(material).Error; err != nil {
				// 物料不存在，创建新物料
				material = &model.Material{
					ItemID:   itemID,
					ItemName: itemName,
					Version:  version,
				}
				if err := s.db.WithContext(ctx).Create(material).Error; err != nil {
					result.FailCount++
					result.Errors = append(result.Errors, fmt.Sprintf("第%d行创建物料失败: %v", i+1, err))
					continue
				}
			}
		}

		// 添加BOM项（只有level > 0的才是子项）
		if level > 0 && material != nil {
			req := &model.AddBOMItemRequest{
				MaterialID: material.ID,
				Version:    version,
				Quantity:   quantity,
			}
			if _, err := s.AddBOMItem(ctx, bomViewID, req); err != nil {
				result.FailCount++
				result.Errors = append(result.Errors, fmt.Sprintf("第%d行添加BOM项失败: %v", i+1, err))
				continue
			}
		}

		result.SuccessCount++
	}

	return result, nil
}

// GenerateImportTemplate 生成导入模板
func (s *Service) GenerateImportTemplate() (*bytes.Buffer, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "BOM结构"

	// 设置表头
	headers := []string{"层级", "物料编码", "物料名称", "版本", "数量", "单位", "排序"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	// 设置表头样式
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#CCCCCC"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	f.SetCellStyle(sheetName, "A1", "G1", headerStyle)

	// 添加示例数据
	exampleData := [][]interface{}{
		{0, "ROOT-001", "根物料", "AA", 1, "个", 0},
		{1, "PART-001", "零件1", "AA", 2, "个", 1},
		{1, "PART-002", "零件2", "AA", 1, "个", 2},
		{2, "PART-003", "子零件1", "AA", 4, "个", 1},
	}
	for i, data := range exampleData {
		rowNum := i + 2
		for j, value := range data {
			cell, _ := excelize.CoordinatesToCellName(j+1, rowNum)
			f.SetCellValue(sheetName, cell, value)
		}
	}

	// 设置列宽
	f.SetColWidth(sheetName, "A", "A", 8)
	f.SetColWidth(sheetName, "B", "B", 20)
	f.SetColWidth(sheetName, "C", "C", 30)
	f.SetColWidth(sheetName, "D", "D", 8)
	f.SetColWidth(sheetName, "E", "E", 10)
	f.SetColWidth(sheetName, "F", "F", 10)
	f.SetColWidth(sheetName, "G", "G", 8)

	buf := new(bytes.Buffer)
	if err := f.Write(buf); err != nil {
		return nil, err
	}

	return buf, nil
}
