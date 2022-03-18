package procedures

import (
	"gorm.io/gorm"
)

func (p *Procedure) Get(db *gorm.DB) error {
	if result := db.Where("id = ?", p.ID).First(&p); result.Error != nil {
		return result.Error
	}
	return nil
}

func Index(db *gorm.DB, limit int, offset int) (*[]Procedure, error) {
	var procedures *[]Procedure
	if result := db.Limit(limit).Offset(offset).Find(&procedures); result.Error != nil {
		return nil, result.Error
	}
	return procedures, nil
}

func CountAll(db *gorm.DB) (int64, error) {
	var count int64
	if result := db.Table("procedures").Count(&count); result.Error != nil {
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

func (p *Procedure) Delete(db *gorm.DB) error {
	// https://gorm.io/ja_JP/docs/error_handling.html
	if result := db.Delete(&p); result.Error != nil {
		return result.Error
	}
	return nil
}
