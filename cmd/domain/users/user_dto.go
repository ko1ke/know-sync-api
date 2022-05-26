package users

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	// gorm.Model
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Username  string         `json:"username" binding:"required,max=100"`
	Email     string         `json:"email" binding:"required,max=255,email"`
	Password  string         `json:"password" binding:"required,max=255"`
}

type Authentication struct {
	Email    string `json:"email" binding:"required,max=255,email"`
	Password string `json:"password" binding:"required,max=255"`
}
