package steps

import (
	"gorm.io/gorm"
)

func (s *Step) Get(db *gorm.DB) error {
	if result := db.Where("id = ?", s.ID).First(&s); result.Error != nil {
		return result.Error
	}
	return nil
}

func Index(db *gorm.DB, procedureId uint) ([]Step, error) {
	var steps []Step
	if result := db.Where("procedure_id = ?", procedureId).Find(&steps); result.Error != nil {
		return nil, result.Error
	}
	return steps, nil
}

func BulkCreate(tx *gorm.DB, steps []Step) ([]Step, error) {
	if result := tx.Save(&steps); result.Error != nil {
		return nil, result.Error
	}
	return steps, nil
}

func BulkDeleteByProcedureId(tx *gorm.DB, procedureId uint) error {
	var step Step
	if result := tx.Where("procedure_id = ?", procedureId).Delete(&step); result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *Step) Update(db *gorm.DB) error {
	if result := db.Save(&s); result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *Step) Save(db *gorm.DB) error {
	if result := db.Create(&s); result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *Step) Delete(db *gorm.DB) error {
	if result := db.Delete(&s); result.Error != nil {
		return result.Error
	}
	return nil
}
