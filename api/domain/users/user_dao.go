package users

import "know-sync-api/datasources/postgres_db"

func (u *User) Get() error {
	if result := postgres_db.Client.Where("id = ?", u.ID).First(&u); result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *User) GetByEmail() error {
	if result := postgres_db.Client.Where("email = ?", u.Email).First(&u); result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *User) Save() error {
	// https://gorm.io/ja_JP/docs/error_handling.html
	if result := postgres_db.Client.Create(&u); result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *User) Delete() error {
	// https://gorm.io/ja_JP/docs/error_handling.html
	if result := postgres_db.Client.Delete(&u); result.Error != nil {
		return result.Error
	}
	return nil
}
