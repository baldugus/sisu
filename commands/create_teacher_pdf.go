package commands

import (
	"fmt"

	"github.com/baldugus/sisu/database"
	"github.com/baldugus/sisu/pdfbuilder"
	"github.com/baldugus/sisu/types"
)

// CreateTeacherPDFCommand lists the students currently enrolled in one
// semester and time slot.
type CreateTeacherPDFCommand struct {
	Period   types.CoursePeriod
	Semester int32
	FilePath string
}

func (cmd *CreateTeacherPDFCommand) Execute(db *database.Database) error {
	if cmd.Semester != 1 && cmd.Semester != 2 {
		return ErrInvalidSemester{}
	}

	selection, err := db.FetchSelection(types.SelectionKindApproved)
	if err != nil {
		return fmt.Errorf("fetch selection: %w", err)
	}

	courses, err := db.FetchCoursesByPeriod(cmd.Period)
	if err != nil {
		return fmt.Errorf("fetch courses: %w", err)
	}

	var courseInfos []*pdfbuilder.CourseInfo
	for _, course := range courses {
		registrations, err := db.FetchEnrolledRegistrationsByCourseAndSemester(course.ID, cmd.Semester)
		if err != nil {
			return fmt.Errorf("fetch registrations for course %d: %w", course.ID, err)
		}
		courseInfos = append(courseInfos, pdfbuilder.NewCourseInfo(course.Quota, registrations))
	}

	selectionInfo := pdfbuilder.NewSelectionInfo(selection, cmd.Semester, 0)

	builder := &pdfbuilder.Builder{
		Period:    coursePeriodToPortuguese(cmd.Period),
		Selection: selectionInfo,
		Courses:   courseInfos,
	}

	if err := builder.BuildTeacherPdf(cmd.FilePath); err != nil {
		return fmt.Errorf("build teacher pdf: %w", err)
	}

	return nil
}
