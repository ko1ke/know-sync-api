package procedures

import (
	"time"

	"github.com/ko1ke/know-sync-api/domain/steps"
	"gorm.io/gorm"
)

type Procedure struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	Title           string         `json:"title" binding:"required,max=100"`
	Content         string         `json:"content" binding:"max=1000"`
	UserID          uint           `json:"userId"`
	Publish         bool           `json:"publish"`
	EyeCatchImgName string         `json:"eyeCatchImgName" binding:"max=255"`
	Steps           []steps.Step   `json:"steps"`
}

type ProcedureItem struct {
	Procedure
	Username string `json:"username"`
}
