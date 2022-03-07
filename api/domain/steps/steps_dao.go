package steps

import (
	"know-sync-api/datasources/postgres_db"
)

func (s *Step) Get() error {
	if result := postgres_db.Client.Where("id = ?", s.ID).First(&s); result.Error != nil {
		return result.Error
	}
	return nil
}

func Index(procedureId uint) ([]Step, error) {
	var steps []Step
	if result := postgres_db.Client.Where("procedure_id = ?", procedureId).Find(&steps); result.Error != nil {
		return nil, result.Error
	}
	return steps, nil
}

func (s *Step) Update() error {
	if result := postgres_db.Client.Save(&s); result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *Step) Save() error {
	if result := postgres_db.Client.Create(&s); result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *Step) Delete() error {
	if result := postgres_db.Client.Delete(&s); result.Error != nil {
		return result.Error
	}
	return nil
}
