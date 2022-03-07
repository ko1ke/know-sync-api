package procedures

import (
	"know-sync-api/domain/steps"
	"time"

	"gorm.io/gorm"
)

type Procedure struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Title     string         `json:"title"`
	Content   string         `json:"content"`
	UserID    uint           `json:"userId"`
	Steps     []steps.Step  `json:"steps"`
}
