package pdfbuilder

import (
	"fmt"
	"sort"
	"strings"

	"github.com/baldugus/sisu/types"
	"github.com/johnfercher/maroto/v2/pkg/components/list"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

type EmailsListItem struct {
	order        int
	name         string
	score        *types.Score
	enrollmentID string
	email        string
}

func newEmailListItem(order int, name string, score *types.Score, enrollmentID string, email string) EmailsListItem {
	return EmailsListItem{
		order:        order,
		name:         name,
		score:        score,
		enrollmentID: enrollmentID,
		email:        email,
	}
}

func (e EmailsListItem) GetContent(_ int) core.Row {
	return row.New(5).Add(
		text.NewCol(1, fmt.Sprint(e.order)),
		text.NewCol(5, e.name),
		text.NewCol(1, fmt.Sprint(e.score)),
		text.NewCol(2, e.enrollmentID),
		text.NewCol(5, e.email),
	)
}

func (e EmailsListItem) GetHeader() core.Row {
	var cellHeaderText props.Text
	cellHeaderText.Style = fontstyle.Bold

	return row.New(10).Add(
		text.NewCol(1, "ORDEM", cellHeaderText),
		text.NewCol(5, "NOME DO CANDIDATO", cellHeaderText),
		text.NewCol(1, "NOTA", cellHeaderText),
		text.NewCol(2, "INSCRIÇÃO DO ENEM", cellHeaderText),
		text.NewCol(5, "EMAIL", cellHeaderText),
	)
}

type EmailRenderer struct{}

func (e *EmailRenderer) Render(courses []*CourseInfo) ([]core.Row, error) {
	var titleText props.Text
	titleText.Top = 10
	titleText.Style = fontstyle.Bold
	titleText.Align = align.Center

	var rows []core.Row

	for _, course := range courses {
		// Skip courses with no registrations
		if len(course.Registrations) == 0 {
			continue
		}

		sort.Sort(course.Registrations)

		var items []EmailsListItem

		rows = append(rows, text.NewRow(20, course.Quota, titleText))

		for i, registration := range course.Registrations {
			item := newEmailListItem(i+1, registration.Candidate.Name, registration.CompositeScore, registration.EnrollmentID, registration.Candidate.Email)
			items = append(items, item)
		}

		builtRows, err := list.Build(items)
		if err != nil {
			return nil, fmt.Errorf("list build: %w", err)
		}

		rows = append(rows, builtRows...)
	}

	return rows, nil
}

func (e *EmailRenderer) header(selection *SelectionInfo, period string) []core.Row {
	var call strings.Builder

	call.WriteString(string(selection.Kind))
	if selection.WaitlistNum > 0 {
		call.WriteString(fmt.Sprintf(" %d", selection.WaitlistNum))
	}

	call.WriteString(fmt.Sprintf(" - %d - ", selection.Year))
	call.WriteString(fmt.Sprintf("%do. Semestre", selection.Semester))
	period = fmt.Sprintf("TURNO: %s", period)

	var headerText props.Text
	headerText.Size = 11
	headerText.Style = fontstyle.Bold
	headerText.Align = align.Center

	return []core.Row{
		text.NewRow(7, selection.Institution, headerText),
		text.NewRow(7, call.String(), headerText),
		text.NewRow(7, selection.Course, headerText),
		text.NewRow(7, period, headerText),
	}
}
