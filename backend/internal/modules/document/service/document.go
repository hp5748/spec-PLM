package service

import (
	"context"
	"errors"
	"mime/multipart"
	"path/filepath"
	"time"

	"plm/internal/modules/document/model"
	"plm/internal/modules/document/repository"
	sharedModel "plm/internal/shared/model"
	"plm/pkg/storage"

	"gorm.io/gorm"
)

type Service struct {
	repo    *repository.Repository
	storage *storage.MinIOStorage
	db      *gorm.DB
}

func NewService(repo *repository.Repository, storage *storage.MinIOStorage) *Service {
	return &Service{repo: repo, storage: storage}
}

// SetDB 设置数据库连接
func (s *Service) SetDB(db *gorm.DB) {
	s.db = db
}

// UploadDocument 上传文档
func (s *Service) UploadDocument(ctx context.Context, req *model.CreateDocumentRequest, file *multipart.FileHeader, createdBy uint) (*model.Document, error) {
	// 检查文档编码是否已存在
	exists, err := s.repo.ExistsByDocID(ctx, req.DocID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("文档编码已存在")
	}

	// 上传文件到MinIO
	folder := "documents/" + time.Now().Format("2006/01/02")
	filePath, err := s.storage.UploadFile(ctx, file, folder)
	if err != nil {
		return nil, errors.New("文件上传失败: " + err.Error())
	}

	// 处理空attributes
	attributes := req.Attributes
	if attributes == "" {
		attributes = "{}"
	}

	doc := &model.Document{
		DocID:      req.DocID,
		DocName:    req.DocName,
		DocType:    req.DocType,
		FilePath:   filePath,
		FileName:   file.Filename,
		FileSize:   file.Size,
		MimeType:   file.Header.Get("Content-Type"),
		Version:    "AA",
		Status:     "DRAFT",
		Attributes: attributes,
		CreatedBy:  createdBy,
	}

	if err := s.repo.CreateDocument(ctx, doc); err != nil {
		// 删除已上传的文件
		s.storage.DeleteFile(ctx, filePath)
		return nil, err
	}

	return doc, nil
}

// GetDocument 获取文档详情
func (s *Service) GetDocument(ctx context.Context, id uint) (*model.Document, error) {
	return s.repo.GetDocumentByID(ctx, id)
}

// GetDocumentURL 获取文档下载URL
func (s *Service) GetDocumentURL(ctx context.Context, id uint) (string, error) {
	doc, err := s.repo.GetDocumentByID(ctx, id)
	if err != nil {
		return "", errors.New("文档不存在")
	}

	url, err := s.storage.GetFileURL(ctx, doc.FilePath, time.Hour)
	if err != nil {
		return "", errors.New("获取下载链接失败")
	}

	return url, nil
}

// UpdateDocument 更新文档
func (s *Service) UpdateDocument(ctx context.Context, id uint, req *model.UpdateDocumentRequest) (*model.Document, error) {
	doc, err := s.repo.GetDocumentByID(ctx, id)
	if err != nil {
		return nil, errors.New("文档不存在")
	}

	if req.DocName != "" {
		doc.DocName = req.DocName
	}
	if req.DocType != "" {
		doc.DocType = req.DocType
	}
	if req.Status != "" {
		doc.Status = req.Status
	}
	if req.Attributes != "" {
		doc.Attributes = req.Attributes
	}

	if err := s.repo.UpdateDocument(ctx, doc); err != nil {
		return nil, err
	}

	return doc, nil
}

// DeleteDocument 删除文档
func (s *Service) DeleteDocument(ctx context.Context, id uint) error {
	doc, err := s.repo.GetDocumentByID(ctx, id)
	if err != nil {
		return errors.New("文档不存在")
	}

	// 检查状态是否允许删除
	if doc.Status == "RELEASED" {
		return errors.New("已发布的文档不能删除")
	}

	// 删除文件
	if err := s.storage.DeleteFile(ctx, doc.FilePath); err != nil {
		// 记录错误但继续删除数据库记录
	}

	return s.repo.DeleteDocument(ctx, id)
}

// ListDocuments 文档列表
func (s *Service) ListDocuments(ctx context.Context, query *model.DocumentListQuery) ([]*model.Document, int64, error) {
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 20
	}

	return s.repo.ListDocuments(ctx, query)
}

// SearchDocuments 搜索文档
func (s *Service) SearchDocuments(ctx context.Context, keyword string, page, pageSize int) ([]*model.Document, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	return s.repo.SearchDocuments(ctx, keyword, page, pageSize)
}

// GetDocTypeFromExt 根据文件扩展名获取文档类型
func GetDocTypeFromExt(filename string) string {
	ext := filepath.Ext(filename)
	switch ext {
	case ".pdf":
		return "PDF"
	case ".doc", ".docx":
		return "DOC"
	case ".xls", ".xlsx":
		return "XLS"
	case ".ppt", ".pptx":
		return "PPT"
	case ".dwg", ".dxf":
		return "CAD"
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp":
		return "IMAGE"
	case ".zip", ".rar", ".7z":
		return "ARCHIVE"
	default:
		return "OTHER"
	}
}

// IncrementVersion 文档版本升级
func (s *Service) IncrementVersion(ctx context.Context, id uint, reason string, createdBy uint) (*model.Document, error) {
	doc, err := s.repo.GetDocumentByID(ctx, id)
	if err != nil {
		return nil, errors.New("文档不存在")
	}

	// 检查状态是否允许版本升级
	if doc.Status == "REVIEWING" {
		return nil, errors.New("审核中的文档不能升级版本")
	}

	oldVersion := doc.Version
	newVersion := s.calculateNextVersion(doc.Version)

	// 使用事务
	if s.db != nil {
		err = s.db.Transaction(func(tx *gorm.DB) error {
			// 更新文档版本
			doc.Version = newVersion
			if err := s.repo.UpdateDocument(ctx, doc); err != nil {
				return err
			}

			// 记录版本历史
			history := &sharedModel.VersionHistory{
				EntityType:   "DOCUMENT",
				EntityID:     id,
				OldVersion:   oldVersion,
				NewVersion:   newVersion,
				ChangeReason: reason,
				CreatedBy:    createdBy,
			}
			return tx.Create(history).Error
		})
		if err != nil {
			return nil, err
		}
	} else {
		// 没有数据库连接时只更新版本
		doc.Version = newVersion
		if err := s.repo.UpdateDocument(ctx, doc); err != nil {
			return nil, err
		}
	}

	return s.repo.GetDocumentByID(ctx, id)
}

// GetVersionHistory 获取文档版本历史
func (s *Service) GetVersionHistory(ctx context.Context, id uint) ([]sharedModel.VersionHistory, error) {
	var histories []sharedModel.VersionHistory
	if s.db == nil {
		return histories, nil
	}

	err := s.db.WithContext(ctx).
		Where("entity_type = ? AND entity_id = ?", "DOCUMENT", id).
		Order("created_at DESC").
		Find(&histories).Error
	return histories, err
}

// GetPreviewURL 获取文档预览URL
func (s *Service) GetPreviewURL(ctx context.Context, id uint) (string, error) {
	doc, err := s.repo.GetDocumentByID(ctx, id)
	if err != nil {
		return "", errors.New("文档不存在")
	}

	// 生成24小时有效的预览URL
	url, err := s.storage.GetFileURL(ctx, doc.FilePath, 24*time.Hour)
	if err != nil {
		return "", errors.New("获取预览链接失败")
	}

	return url, nil
}

// calculateNextVersion 计算下一个版本号
// 版本规则：AA-YY，跳过I/O/Z
func (s *Service) calculateNextVersion(current string) string {
	if len(current) != 2 {
		return "AA"
	}

	// 跳过的字母
	skipChars := map[byte]bool{'I': true, 'O': true, 'Z': true}

	first := current[0]
	second := current[1]

	// 递增第二位
	second++
	for skipChars[second] {
		second++
	}

	// 如果第二位超过'Y'，进位
	if second > 'Y' {
		second = 'A'
		first++
		for skipChars[first] {
			first++
		}
	}

	// 如果第一位超过'Y'，回到AA
	if first > 'Y' {
		return "AA"
	}

	return string([]byte{first, second})
}
