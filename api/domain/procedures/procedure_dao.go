package procedures

import (
	postgres_db "know-sync-api/datasources/postgres_db"
)

func (p *Procedure) Get() error {
	if result := postgres_db.Client.Where("id = ?", p.ID).First(&p); result.Error != nil {
		return result.Error
	}
	return nil
}

func GetAll() (*[]Procedure, error) {
	var procedures *[]Procedure
	if result := postgres_db.Client.Find(&procedures); result.Error != nil {
		return nil, result.Error
	}
	return procedures, nil
}

func (p *Procedure) Update() error {
	if result := postgres_db.Client.Save(&p); result.Error != nil {
		return result.Error
	}
	return nil
}

func (p *Procedure) PartialUpdate() error {
	if result := postgres_db.Client.
		Table("procedures").
		Where("id IN (?)", p.ID).
		Updates(&p); result.Error != nil {
		return result.Error
	}
	return nil
}

func (p *Procedure) Save() error {
	// https://gorm.io/ja_JP/docs/error_handling.html
	if result := postgres_db.Client.Create(&p); result.Error != nil {
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
