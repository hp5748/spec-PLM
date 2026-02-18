package repository

import (
	"context"

	"gorm.io/gorm"
	"plm/internal/modules/document/model"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

// GetDocumentByID 根据ID获取文档
func (r *Repository) GetDocumentByID(ctx context.Context, id uint) (*model.Document, error) {
	var doc model.Document
	err := getDB(ctx).First(&doc, id).Error
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

// GetDocumentByDocID 根据文档编码获取文档
func (r *Repository) GetDocumentByDocID(ctx context.Context, docID string, version string) (*model.Document, error) {
	var doc model.Document
	err := getDB(ctx).Where("doc_id = ? AND version = ?", docID, version).First(&doc).Error
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

// ExistsByDocID 检查文档编码是否存在
func (r *Repository) ExistsByDocID(ctx context.Context, docID string) (bool, error) {
	var count int64
	err := getDB(ctx).Model(&model.Document{}).Where("doc_id = ?", docID).Count(&count).Error
	return count > 0, err
}

// CreateDocument 创建文档
func (r *Repository) CreateDocument(ctx context.Context, doc *model.Document) error {
	return getDB(ctx).Create(doc).Error
}

// UpdateDocument 更新文档
func (r *Repository) UpdateDocument(ctx context.Context, doc *model.Document) error {
	return getDB(ctx).Save(doc).Error
}

// DeleteDocument 删除文档
func (r *Repository) DeleteDocument(ctx context.Context, id uint) error {
	return getDB(ctx).Delete(&model.Document{}, id).Error
}

// ListDocuments 文档列表（分页、筛选）
func (r *Repository) ListDocuments(ctx context.Context, query *model.DocumentListQuery) ([]*model.Document, int64, error) {
	var docs []*model.Document
	var total int64

	db := getDB(ctx).Model(&model.Document{})

	// 筛选条件
	if query.DocID != "" {
		db = db.Where("doc_id LIKE ?", "%"+query.DocID+"%")
	}
	if query.DocName != "" {
		db = db.Where("doc_name LIKE ?", "%"+query.DocName+"%")
	}
	if query.DocType != "" {
		db = db.Where("doc_type = ?", query.DocType)
	}
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}
	if query.Keyword != "" {
		db = db.Where("doc_id LIKE ? OR doc_name LIKE ?", "%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}

	// 计数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	order := "created_at DESC"
	if query.Sort != "" {
		order = query.Sort
		if query.Order != "" {
			order = query.Sort + " " + query.Order
		}
	}

	// 分页
	offset := (query.Page - 1) * query.PageSize
	if err := db.Order(order).Offset(offset).Limit(query.PageSize).Find(&docs).Error; err != nil {
		return nil, 0, err
	}

	return docs, total, nil
}

// SearchDocuments 搜索文档
func (r *Repository) SearchDocuments(ctx context.Context, keyword string, page, pageSize int) ([]*model.Document, int64, error) {
	var docs []*model.Document
	var total int64

	db := getDB(ctx).Model(&model.Document{}).Where(
		"doc_id LIKE ? OR doc_name LIKE ?",
		"%"+keyword+"%", "%"+keyword+"%",
	)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&docs).Error; err != nil {
		return nil, 0, err
	}

	return docs, total, nil
}

// getDB 获取数据库实例
func getDB(ctx context.Context) *gorm.DB {
	return globalDB
}

var globalDB *gorm.DB

func SetDB(db *gorm.DB) {
	globalDB = db
}
