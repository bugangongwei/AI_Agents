package model

import (
	"time"

	"gorm.io/gorm"
)

type JobDescription struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	UserID     uint           `gorm:"not null;index" json:"user_id"`
	Title      string         `json:"title"`
	Company    string         `json:"company"`
	RawText    string         `gorm:"type:longtext" json:"raw_text"`
	SourceType string         `gorm:"default:'text'" json:"source_type"` // text or image
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

type Evaluation struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	TargetID   uint           `gorm:"not null;index" json:"target_id"`
	TargetType string         `gorm:"not null" json:"target_type"` // greeting or answer
	Score      float32        `json:"score"`
	Reason     string         `gorm:"type:text" json:"reason"`
	JudgeModel string         `json:"judge_model"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
