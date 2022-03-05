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
	Title       string         `json:"title"`
	Content     string         `json:"content"`
	ImgName     string         `json:"imgName"`
	ProcedureId uint           `json:"procedureId"`
}
