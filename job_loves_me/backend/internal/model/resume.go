package model

import (
	"time"

	"gorm.io/gorm"
)

type Resume struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	FilePath  string         `gorm:"not null" json:"file_path"`
	RawText   string         `gorm:"type:longtext" json:"raw_text"`
	Version   int            `gorm:"default:1" json:"version"`
	IsLatest  bool           `gorm:"default:true" json:"is_latest"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
