package commands

import (
	"fmt"

	"github.com/baldugus/sisu/database"
	"github.com/baldugus/sisu/pdfbuilder"
	"github.com/baldugus/sisu/types"
)

// callPDFArgs selects the students of one call, time slot and semester.
// Promotions appear in the semester-1 list like any other called student.
type callPDFArgs struct {
	CallID   int32
	Period   types.CoursePeriod
	Semester int32
	FilePath string
}

type CreateWebsitePDFCommand callPDFArgs

func (cmd *CreateWebsitePDFCommand) Execute(db *database.Database) error {
	builder, err := callPDFBuilder(db, (*callPDFArgs)(cmd))
	if err != nil {
		return err
	}

	if err := builder.BuildWebsitePdf(cmd.FilePath); err != nil {
		return fmt.Errorf("build website pdf: %w", err)
	}

	return nil
}

type CreateEnrollmentPDFCommand callPDFArgs

func (cmd *CreateEnrollmentPDFCommand) Execute(db *database.Database) error {
	builder, err := callPDFBuilder(db, (*callPDFArgs)(cmd))
	if err != nil {
		return err
	}

	if err := builder.BuildEnrollmentPdf(cmd.FilePath); err != nil {
		return fmt.Errorf("build enrollment pdf: %w", err)
	}

	return nil
}

type CreateEmailPDFCommand callPDFArgs

func (cmd *CreateEmailPDFCommand) Execute(db *database.Database) error {
	builder, err := callPDFBuilder(db, (*callPDFArgs)(cmd))
	if err != nil {
		return err
	}

	if err := builder.BuildEmailPdf(cmd.FilePath); err != nil {
		return fmt.Errorf("build email pdf: %w", err)
	}

	return nil
}

func callPDFBuilder(db *database.Database, args *callPDFArgs) (*pdfbuilder.Builder, error) {
	if args.Semester != 1 && args.Semester != 2 {
		return nil, ErrInvalidSemester{}
	}

	call, err := db.FetchCallByID(args.CallID)
	if err != nil {
		return nil, fmt.Errorf("fetch call: %w", err)
	}

	selectionKind := types.SelectionKindApproved
	if call.Number > 1 {
		selectionKind = types.SelectionKindWaitlist
	}

	selection, err := db.FetchSelection(selectionKind)
	if err != nil {
		return nil, fmt.Errorf("fetch selection: %w", err)
	}

	courses, err := db.FetchCoursesByPeriod(args.Period)
	if err != nil {
		return nil, fmt.Errorf("fetch courses: %w", err)
	}

	var courseInfos []*pdfbuilder.CourseInfo

	for _, course := range courses {
		details, err := database.FetchCallEntryDetails(db.DB(), args.CallID, &course.ID, &args.Semester)
		if err != nil {
			return nil, fmt.Errorf("fetch call entries for course %d: %w", course.ID, err)
		}

		registrations := make(types.Registrations, len(details))
		for i, d := range details {
			registrations[i] = d.Registration
		}

		courseInfos = append(courseInfos, pdfbuilder.NewCourseInfo(course.Quota, registrations))
	}

	// Call 1 is the regular call; call N is waitlist call N-1.
	waitlistNum := int64(0)
	if call.Number > 1 {
		waitlistNum = int64(call.Number - 1)
	}

	return &pdfbuilder.Builder{
		Period:    coursePeriodToPortuguese(args.Period),
		Selection: pdfbuilder.NewSelectionInfo(selection, args.Semester, waitlistNum),
		Courses:   courseInfos,
	}, nil
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
