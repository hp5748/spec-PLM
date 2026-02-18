package model

import (
	"time"
)

// VersionHistory 版本历史模型
type VersionHistory struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	EntityType   string    `gorm:"size:50;not null;index:idx_entity" json:"entity_type"`
	EntityID     uint      `gorm:"not null;index:idx_entity" json:"entity_id"`
	OldVersion   string    `gorm:"size:2" json:"old_version"`
	NewVersion   string    `gorm:"size:2;not null" json:"new_version"`
	ChangeReason string    `json:"change_reason"`
	CreatedAt    time.Time `json:"created_at"`
	CreatedBy    uint      `gorm:"not null" json:"created_by"`
}

// TableName 指定表名
func (VersionHistory) TableName() string {
	return "version_histories"
}
