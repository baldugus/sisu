package pdfbuilder

import "github.com/baldugus/sisu/types"

type ClassInfo struct {
	Quota        string
	Applications types.Applications
}

func NewClassInfo(quota string, applications types.Applications) *ClassInfo {
	return &ClassInfo{
		Quota:        quota,
		Applications: applications,
	}
}

type SelectionInfo struct {
	Institution string
	Kind        SelectionKind
	Year        string
	Semester    string
	Course      string
	WaitlistNum int64
}

func NewSelectionInfo(selection *types.Selection, year string, semester string, waitlistNum int64) *SelectionInfo {
	return &SelectionInfo{
		Institution: selection.Institution,
		Kind:        mapSelectionKind(selection.Kind),
		Year:        year,
		Semester:    semester,
		Course:      selection.Course,
		WaitlistNum: waitlistNum,
	}
}

type SelectionKind string

const (
	RegularCall SelectionKind = "CHAMADA REGULAR"
	WaitingList SelectionKind = "LISTA DE ESPERA"
)

func mapSelectionKind(kind types.SelectionKind) SelectionKind {
	if kind == 1 {
		return RegularCall
	}

	return WaitingList
}
