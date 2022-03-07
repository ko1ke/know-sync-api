package services

import (
	"know-sync-api/domain/steps"
)

// func CreateProcedure(procedure procedures.Procedure) (*procedures.Procedure, error) {
// 	if err := procedure.Save(); err != nil {
// 		return nil, err
// 	}

// 	return &procedure, nil
// }

// func UpdateProcedure(isPartial bool, procedure procedures.Procedure) (*procedures.Procedure, error) {
// 	current, err := GetProcedure(procedure.ID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if err = current.Get(); err != nil {
// 		return nil, err
// 	}

// 	if isPartial {
// 		if procedure.Title != "" {
// 			current.Title = procedure.Title
// 		}
// 		if procedure.Content != "" {
// 			current.Content = procedure.Content
// 		}

// 		if err := current.PartialUpdate(); err != nil {
// 			return nil, err
// 		}
// 	} else {
// 		current.Title = procedure.Title
// 		current.Content = procedure.Content
// 		if err := current.Update(); err != nil {
// 			return nil, err
// 		}
// 	}
// 	return current, nil
// }

func GetSteps(procedureId uint) ([]steps.Step, error) {
	steps, err := steps.Index(procedureId)
	if err != nil {
		return nil, err
	}
	return steps, nil
}

// func GetPagination(page int, limit int) (*pagination_utils.Pagination, error) {
// 	itemsCount, err := procedures.CountAll()
// 	if err != nil {
// 		return nil, err
// 	}

// 	pagination := pagination_utils.NewPagination(page, limit, int(itemsCount))
// 	return pagination, nil
// }

// func DeleteProcedure(procedureID uint) error {
// 	procedure := &procedures.Procedure{ID: procedureID}
// 	return procedure.Delete()
// }
