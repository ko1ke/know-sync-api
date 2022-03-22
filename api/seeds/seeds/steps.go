package seeds

import (
	"know-sync-api/domain/procedures"
	"know-sync-api/domain/steps"
	"know-sync-api/domain/users"

	"gorm.io/gorm"
)

func CreateStep(db *gorm.DB, title string, content string) error {
	seedUser := users.User{Email: SeedUserEmail}
	seedUser.GetByEmail(db)

	seedProcedure := procedures.Procedure{UserID: seedUser.ID}
	seedProcedure.GetByUserId(db)

	return db.Create(&steps.Step{Title: title, Content: content, ProcedureID: seedProcedure.ID}).Error
}
