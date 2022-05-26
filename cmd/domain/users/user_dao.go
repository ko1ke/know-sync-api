package users

import (
	"gorm.io/gorm"
)

func (u *User) Get(db *gorm.DB) error {
	if result := db.Where("id = ?", u.ID).First(&u); result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *User) GetByEmail(db *gorm.DB) error {
	if result := db.Where("email = ?", u.Email).First(&u); result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *User) Save(db *gorm.DB) error {
	// https://gorm.io/ja_JP/docs/error_handling.html
	if result := db.Create(&u); result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *User) Delete(db *gorm.DB) error {
	// https://gorm.io/ja_JP/docs/error_handling.html
	if result := db.Delete(&u); result.Error != nil {
		return result.Error
	}
	return nil
}
