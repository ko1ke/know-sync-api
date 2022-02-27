package procedures

import (
	"time"

	"gorm.io/gorm"
)

type Procedure struct {
	// gorm.Model
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Title     string         `json:"title"`
	Content   string         `json:"content"`
	UserId    uint           `json:"userId"`
}
