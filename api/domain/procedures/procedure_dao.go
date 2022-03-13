package procedures

import (
	postgres_db "know-sync-api/datasources/postgres_db"

	"gorm.io/gorm"
)

func (p *Procedure) Get() error {
	if result := postgres_db.Client.Where("id = ?", p.ID).First(&p); result.Error != nil {
		return result.Error
	}
	return nil
}

func Index(limit int, offset int) (*[]Procedure, error) {
	var procedures *[]Procedure
	if result := postgres_db.Client.Limit(limit).Offset(offset).Find(&procedures); result.Error != nil {
		return nil, result.Error
	}
	return procedures, nil
}

func CountAll() (int64, error) {
	var count int64
	if result := postgres_db.Client.Table("procedures").Count(&count); result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

func (p *Procedure) Update(tx *gorm.DB) error {
	if result := tx.Save(&p); result.Error != nil {
		return result.Error
	}
	return nil
}

func (p *Procedure) PartialUpdate(tx *gorm.DB) error {
	if result := tx.Table("procedures").
		Where("id IN (?)", p.ID).
		Updates(&p); result.Error != nil {
		return result.Error
	}
	return nil
}

func (p *Procedure) Save(tx *gorm.DB) error {
	// https://gorm.io/ja_JP/docs/error_handling.html
	if result := tx.Create(&p); result.Error != nil {
		return result.Error
	}
	return nil
}

func (p *Procedure) Delete() error {
	// https://gorm.io/ja_JP/docs/error_handling.html
	if result := postgres_db.Client.Delete(&p); result.Error != nil {
		return result.Error
	}
	return nil
}
