package services

import (
	"know-sync-api/datasources/postgres_db"
	"know-sync-api/domain/procedures"
	"know-sync-api/domain/steps"
	"know-sync-api/utils/pagination_utils"

	"gorm.io/gorm"
)

func CreateProcedure(procedure procedures.Procedure) (*procedures.Procedure, error) {
	err := postgres_db.Client.Transaction(func(tx *gorm.DB) error {
		if err := procedure.Save(tx); err != nil {
			return err
		}

		ss := procedure.Steps
		for i := 0; i < len(ss); i++ {
			ss[i].ProcedureID = procedure.ID
		}

		newSs, err := steps.BulkCreate(tx, ss)
		if err != nil {
			return err
		}

		procedure.Steps = newSs
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &procedure, nil
}

func UpdateProcedure(isPartial bool, procedure procedures.Procedure) (*procedures.Procedure, error) {
	current, err := GetProcedure(procedure.ID)
	if err != nil {
		return nil, err
	}

	if err := current.Get(); err != nil {
		return nil, err
	}

	txErr := postgres_db.Client.Transaction(func(tx *gorm.DB) error {
		if isPartial {
			if procedure.Title != "" {
				current.Title = procedure.Title
			}
			if procedure.Content != "" {
				current.Content = procedure.Content
			}

			if err := current.PartialUpdate(tx); err != nil {
				return err
			}
		} else {
			current.Title = procedure.Title
			current.Content = procedure.Content
			if err := current.Update(tx); err != nil {
				return err
			}
		}

		ss := procedure.Steps
		for i := 0; i < len(ss); i++ {
			ss[i].ProcedureID = current.ID
		}

		delErr := steps.BulkDeleteByProcedureId(tx, procedure.ID, ss)
		if delErr != nil {
			return delErr
		}

		if ss == nil {
			return nil
		}

		newSs, createErr := steps.BulkCreate(tx, ss)
		if createErr != nil {
			return createErr
		}

		current.Steps = newSs
		return nil
	})

	if txErr != nil {
		return nil, txErr
	}

	return current, nil
}

func GetProcedure(procedureID uint) (*procedures.Procedure, error) {
	p := &procedures.Procedure{ID: procedureID}
	if err := p.Get(); err != nil {
		return nil, err
	}

	ss, err := steps.Index(procedureID)
	if err != nil {
		return nil, err
	}

	p.Steps = ss

	return p, nil
}

func GetProcedures(limit int, offset int) (*[]procedures.Procedure, error) {
	ps, err := procedures.Index(limit, offset)
	if err != nil {
		return nil, err
	}
	return ps, nil
}

func GetPagination(page int, limit int) (*pagination_utils.Pagination, error) {
	itemsCount, err := procedures.CountAll()
	if err != nil {
		return nil, err
	}

	pagination := pagination_utils.NewPagination(page, limit, int(itemsCount))
	return pagination, nil
}

func DeleteProcedure(procedureID uint) error {
	procedure := &procedures.Procedure{ID: procedureID}
	return procedure.Delete()
}
