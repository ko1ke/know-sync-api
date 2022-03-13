package steps

import (
	"time"

	"gorm.io/gorm"
)

type Step struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Title       string         `json:"title" binding:"max=100"`
	Content     string         `json:"content" binding:"max=255"`
	ImgName     string         `json:"imgName" binding:"max=255"`
	ProcedureID uint           `json:"procedureId"`
}
