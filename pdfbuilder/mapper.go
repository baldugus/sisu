package pdfbuilder

import "github.com/baldugus/sisu/types"

type CourseInfo struct {
	Quota         string
	Registrations types.Registrations
}

func NewCourseInfo(quota string, registrations types.Registrations) *CourseInfo {
	return &CourseInfo{
		Quota:         quota,
		Registrations: registrations,
	}
}

type SelectionInfo struct {
	Institution string
	Kind        SelectionKind
	Year        int32
	Semester    int32
	Course      string
	WaitlistNum int64
}

func NewSelectionInfo(selection *types.Selection, waitlistNum int64) *SelectionInfo {
	return &SelectionInfo{
		Institution: selection.Institution,
		Kind:        mapSelectionKind(selection.Kind),
		Year:        selection.Year,
		Semester:    selection.Semester,
		Course:      selection.Degree,
		WaitlistNum: waitlistNum,
	}
}

type SelectionKind string

const (
	RegularCall SelectionKind = "CHAMADA REGULAR"
	WaitingList SelectionKind = "LISTA DE ESPERA"
)

func mapSelectionKind(kind types.SelectionKind) SelectionKind {
	if kind == types.SelectionKindApproved {
		return RegularCall
	}

	return WaitingList
}
