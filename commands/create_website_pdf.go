package commands

import (
	"fmt"

	"github.com/baldugus/sisu/database"
	"github.com/baldugus/sisu/pdfbuilder"
	"github.com/baldugus/sisu/types"
)

type CreateWebsitePDFCommand struct {
	CallID   int32
	Period   types.CoursePeriod
	FilePath string
}

func (cmd *CreateWebsitePDFCommand) Execute(db *database.Database) error {
	call, err := db.FetchCallByID(cmd.CallID)
	if err != nil {
		return fmt.Errorf("fetch call: %w", err)
	}

	selectionKind := types.SelectionKindApproved
	if call.Number > 1 {
		selectionKind = types.SelectionKindWaitlist
	}

	selection, err := db.FetchSelection(selectionKind)
	if err != nil {
		return fmt.Errorf("fetch selection: %w", err)
	}

	courses, err := db.FetchCoursesByPeriod(cmd.Period)
	if err != nil {
		return fmt.Errorf("fetch courses: %w", err)
	}

	var courseInfos []*pdfbuilder.CourseInfo
	for _, course := range courses {
		registrations, err := db.FetchRegistrationsByCallIDAndCourseID(cmd.CallID, course.ID)
		if err != nil {
			return fmt.Errorf("fetch registrations for course %d: %w", course.ID, err)
		}
		courseInfos = append(courseInfos, pdfbuilder.NewCourseInfo(course.Quota, registrations))
	}

	selectionInfo := pdfbuilder.NewSelectionInfo(selection, int64(call.Number))

	builder := &pdfbuilder.Builder{
		Period:    coursePeriodToPortuguese(cmd.Period),
		Selection: selectionInfo,
		Courses:   courseInfos,
	}

	if err := builder.BuildWebsitePdf(cmd.FilePath); err != nil {
		return fmt.Errorf("build website pdf: %w", err)
	}

	return nil
}

func coursePeriodToPortuguese(period types.CoursePeriod) string {
	switch period {
	case types.CoursePeriodMorning:
		return "MATUTINO"
	case types.CoursePeriodEvening:
		return "NOTURNO"
	default:
		return period.String()
	}
}
