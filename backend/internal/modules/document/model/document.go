package model

import (
	"time"

	"gorm.io/gorm"
)

// Document 文档模型
type Document struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	DocID       string         `gorm:"size:50;not null;uniqueIndex:uk_doc_version" json:"doc_id"` // 文档编码
	DocName     string         `gorm:"size:200;not null" json:"doc_name"`                          // 文档名称
	DocType     string         `gorm:"size:50;not null;index" json:"doc_type"`                     // 文档类型
	FilePath    string         `gorm:"size:500;not null" json:"file_path"`                         // 文件存储路径
	FileName    string         `gorm:"size:255" json:"file_name"`                                  // 原始文件名
	FileSize    int64          `json:"file_size"`                                                   // 文件大小（字节）
	MimeType    string         `gorm:"size:100" json:"mime_type"`                                   // MIME类型
	Version     string         `gorm:"size:2;not null;default:AA;uniqueIndex:uk_doc_version" json:"version"` // 版本号
	Status      string         `gorm:"size:20;default:DRAFT;index" json:"status"`                  // 状态: DRAFT/CHECKED_IN/RELEASED
	Attributes  string         `gorm:"type:json" json:"attributes"`                                 // 扩展属性JSON
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	CreatedBy   uint           `json:"created_by"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Document) TableName() string {
	return "documents"
}

// CreateDocumentRequest 创建文档请求
type CreateDocumentRequest struct {
	DocID      string `json:"doc_id" binding:"required"`
	DocName    string `json:"doc_name" binding:"required"`
	DocType    string `json:"doc_type" binding:"required"`
	Attributes string `json:"attributes"`
}

// UpdateDocumentRequest 更新文档请求
type UpdateDocumentRequest struct {
	DocName    string `json:"doc_name"`
	DocType    string `json:"doc_type"`
	Status     string `json:"status"`
	Attributes string `json:"attributes"`
}

// DocumentListQuery 文档列表查询参数
type DocumentListQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Sort     string `form:"sort"`
	Order    string `form:"order"`
	DocID    string `form:"doc_id"`
	DocName  string `form:"doc_name"`
	DocType  string `form:"doc_type"`
	Status   string `form:"status"`
	Keyword  string `form:"keyword"` // 关键字搜索
}
