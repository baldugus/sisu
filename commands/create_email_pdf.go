package commands

import (
	"fmt"

	"github.com/baldugus/sisu/database"
	"github.com/baldugus/sisu/pdfbuilder"
	"github.com/baldugus/sisu/types"
)

type CreateEmailPDFCommand struct {
	CallID   int32
	Period   types.CoursePeriod
	FilePath string
}

func (cmd *CreateEmailPDFCommand) Execute(db *database.Database) error {
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

	// Calculate waitlist number (call 1 = regular, call 2 = waitlist 1, etc.)
	waitlistNum := int64(0)
	if call.Number > 1 {
		waitlistNum = int64(call.Number - 1)
	}

	selectionInfo := pdfbuilder.NewSelectionInfo(selection, waitlistNum)

	builder := &pdfbuilder.Builder{
		Period:    coursePeriodToPortuguese(cmd.Period),
		Selection: selectionInfo,
		Courses:   courseInfos,
	}

	if err := builder.BuildEmailPdf(cmd.FilePath); err != nil {
		return fmt.Errorf("build email pdf: %w", err)
	}

	return nil
}
