package services

import (
	"know-sync-api/domain/procedures"
	"know-sync-api/utils/pagination_utils"
)

func CreateProcedure(procedure procedures.Procedure) (*procedures.Procedure, error) {
	if err := procedure.Save(); err != nil {
		return nil, err
	}

	return &procedure, nil
}

func UpdateProcedure(isPartial bool, procedure procedures.Procedure) (*procedures.Procedure, error) {
	current, err := GetProcedure(procedure.ID)
	if err != nil {
		return nil, err
	}

	if err = current.Get(); err != nil {
		return nil, err
	}

	if isPartial {
		if procedure.Title != "" {
			current.Title = procedure.Title
		}
		if procedure.Content != "" {
			current.Content = procedure.Content
		}

		if err := current.PartialUpdate(); err != nil {
			return nil, err
		}
	} else {
		current.Title = procedure.Title
		current.Content = procedure.Content
		if err := current.Update(); err != nil {
			return nil, err
		}
	}
	return current, nil
}

func GetProcedure(procedureID uint) (*procedures.Procedure, error) {
	t := &procedures.Procedure{ID: procedureID}
	if err := t.Get(); err != nil {
		return nil, err
	}
	return t, nil
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
