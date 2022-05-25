package seeds

import (
	"github.com/ko1ke/know-sync-api/domain/procedures"
	"github.com/ko1ke/know-sync-api/domain/steps"
	"github.com/ko1ke/know-sync-api/domain/users"

	"gorm.io/gorm"
)

func CreateStep(db *gorm.DB, content string) error {
	seedUser := users.User{Email: SeedUserEmail}
	seedUser.GetByEmail(db)

	seedProcedure := procedures.Procedure{UserID: seedUser.ID}
	seedProcedure.GetByUserId(db)

	return db.Create(&steps.Step{Content: content, ProcedureID: seedProcedure.ID}).Error
}
