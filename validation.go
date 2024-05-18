package main

import (
	"errors"
	"fmt"

	"changeme/types"
)

var (
	ErrApprovedSelectionPresent = errors.New("Um arquivo com alunos aprovados já foi importado")
	ErrApprovedSelectionMissing = errors.New("Alunos aprovados devem ser importados antes da fila de espera")
	ErrSelectionsImported       = errors.New("Os alunos aprovados e em fila de espera já foram importados")
)

func (s *SISU) ValidateSelection(selection *types.Selection) (bool, error) {
	selections, err := s.repo.FetchSelections()
	if err != nil {
		return false, fmt.Errorf("fetch selections: %w", err)
	}

	if selection.Kind == types.ApprovedSelection && len(selections) > 0 {
		return false, ErrApprovedSelectionPresent
	}

	if selection.Kind == types.InterestedSelection && len(selections) < 1 {
		return false, ErrApprovedSelectionMissing
	}

	if len(selections) >= 2 {
		return false, ErrSelectionsImported
	}

	return true, nil
}
