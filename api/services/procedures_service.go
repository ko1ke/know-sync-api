package services

import "know-sync-api/domain/procedures"

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

func GetProcedures() (*[]procedures.Procedure, error) {
	ts, err := procedures.GetAll()
	if err != nil {
		return nil, err
	}
	return ts, nil
}

func DeleteProcedure(procedureID uint) error {
	procedure := &procedures.Procedure{ID: procedureID}
	return procedure.Delete()
}
