package seeds

import (
	"know-sync-api/domain/users"

	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB, username string, password string, email string) error {
	return db.Create(&users.User{Username: username, Password: password, Email: email}).Error
}
