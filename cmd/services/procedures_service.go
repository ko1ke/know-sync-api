package services

import (
	"github.com/ko1ke/know-sync-api/cmd/datasources/postgres_db"
	"github.com/ko1ke/know-sync-api/cmd/domain/procedures"
	"github.com/ko1ke/know-sync-api/cmd/domain/steps"
	"github.com/ko1ke/know-sync-api/cmd/utils/pagination_utils"

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

		if len(ss) == 0 {
			procedure.Steps = ss
			return nil
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

	if err := current.Get(postgres_db.Client); err != nil {
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
			if procedure.Publish || !procedure.Publish {
				current.Publish = procedure.Publish
			}
			if procedure.EyeCatchImgName != "" {
				current.EyeCatchImgName = procedure.EyeCatchImgName
			}
			if err := current.PartialUpdate(tx); err != nil {
				return err
			}
		} else {
			current.Title = procedure.Title
			current.Content = procedure.Content
			current.Publish = procedure.Publish
			current.EyeCatchImgName = procedure.EyeCatchImgName
			if err := current.Update(tx); err != nil {
				return err
			}
		}

		ss := procedure.Steps
		for i := 0; i < len(ss); i++ {
			ss[i].ProcedureID = current.ID
		}

		delErr := steps.BulkDeleteByProcedureId(tx, procedure.ID)
		if delErr != nil {
			return delErr
		}

		if len(ss) == 0 {
			current.Steps = ss
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
	if err := p.Get(postgres_db.Client); err != nil {
		return nil, err
	}

	ss, err := steps.Index(postgres_db.Client, procedureID)
	if err != nil {
		return nil, err
	}

	p.Steps = ss

	return p, nil
}

func GetProcedureItem(procedureID uint) (*procedures.ProcedureItem, error) {
	p, err := procedures.GetItem(postgres_db.Client, procedureID)
	if err != nil {
		return nil, err
	}

	ss, err := steps.Index(postgres_db.Client, procedureID)
	if err != nil {
		return nil, err
	}

	p.Steps = ss
	return p, nil
}

func GetProcedures(limit int, offset int, keyword string, userID uint) (*[]procedures.ProcedureItem, error) {
	ps, err := procedures.Index(postgres_db.Client, limit, offset, keyword, userID)
	if err != nil {
		return nil, err
	}
	return ps, nil
}

func GetPublicProcedures(limit int, offset int, keyword string) (*[]procedures.ProcedureItem, error) {
	ps, err := procedures.PublicIndex(postgres_db.Client, limit, offset, keyword)
	if err != nil {
		return nil, err
	}
	return ps, nil
}

func GetPagination(page int, limit int, pLen int) (*pagination_utils.Pagination, error) {
	pCountAll, err := procedures.CountAll(postgres_db.Client)
	if err != nil {
		return nil, err
	}

	pagination := pagination_utils.NewPagination(page, limit, pLen, int(pCountAll))
	return pagination, nil
}

func DeleteProcedure(procedureID uint) error {
	procedure := &procedures.Procedure{ID: procedureID}
	return procedure.Delete(postgres_db.Client)
}
