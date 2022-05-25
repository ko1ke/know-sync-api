package seeds

import (
	"github.com/ko1ke/know-sync-api/domain/procedures"
	"github.com/ko1ke/know-sync-api/domain/users"

	"gorm.io/gorm"
)

func CreateProcedure(db *gorm.DB, title string, content string) error {
	seedUser := users.User{Email: SeedUserEmail}
	seedUser.GetByEmail(db)

	return db.Create(&procedures.Procedure{Title: title, Content: content, UserID: seedUser.ID}).Error
}
