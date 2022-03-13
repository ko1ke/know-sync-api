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
	Title     string         `json:"title" binding:"required,max=100"`
	Content   string         `json:"content" binding:"required,max=1000"`
	UserID    uint           `json:"userId"`
	Steps     []steps.Step   `json:"steps"`
}
