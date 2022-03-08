package services

import (
	"know-sync-api/domain/steps"
)

func CreateSteps(ss []steps.Step) ([]steps.Step, error) {
	ss, err := steps.BulkCreate(ss)
	if err != nil {
		return nil, err
	}
	return ss, nil
}

func DeleteAndCreateSteps(procedureId uint, ss []steps.Step) ([]steps.Step, error) {
	delErr := steps.BulkDeleteByProcedureId(procedureId, ss)
	if delErr != nil {
		return nil, delErr
	}
	if ss == nil {
		return nil, nil
	}

	ss, createErr := steps.BulkCreate(ss)
	if createErr != nil {
		return nil, createErr
	}
	return ss, nil

}

func GetSteps(procedureId uint) ([]steps.Step, error) {
	steps, err := steps.Index(procedureId)
	if err != nil {
		return nil, err
	}
	return steps, nil
}