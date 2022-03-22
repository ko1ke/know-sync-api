package seeds

import "gorm.io/gorm"

type Seed struct {
	Name string
	Run  func(*gorm.DB) error
}

var SeedUserEmail = "test@test.com"
