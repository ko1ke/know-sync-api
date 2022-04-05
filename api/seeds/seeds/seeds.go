package seeds

import (
	"gorm.io/gorm"
)

func All() []Seed {
	return []Seed{
		{
			Name: "CreateUser",
			Run: func(db *gorm.DB) error {
				return CreateUser(db, "Seed", "$2a$14$0uA5/5EZoI4Lq8bYXZecMuKYMOBL.Z4rkNHTGHwXoRq.aZUkCEoZa", SeedUserEmail)
			},
		},
		{
			Name: "CreateProcedure",
			Run: func(db *gorm.DB) error {
				return CreateProcedure(db, "Seed title", "Seed content")
			},
		},
		{
			Name: "CreateStep",
			Run: func(db *gorm.DB) error {
				return CreateStep(db, "Seed step content")
			},
		},
	}
}
