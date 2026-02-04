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

type WebsiteItem struct {
	order        int
	name         string
	score        *types.Score
	enrollmentID string
}

func newWebsiteItem(order int, name string, score *types.Score, enrollmentID string) WebsiteItem {
	return WebsiteItem{
		order:        order,
		name:         name,
		score:        score,
		enrollmentID: enrollmentID,
	}
}

func (w WebsiteItem) GetHeader() core.Row {
	var cellHeaderText props.Text
	cellHeaderText.Style = fontstyle.Bold

	return row.New(10).Add(
		text.NewCol(1, "ORDEM", cellHeaderText),
		text.NewCol(7, "NOME DO CANDIDATO", cellHeaderText),
		text.NewCol(1, "NOTA", cellHeaderText),
		text.NewCol(3, "INSCRIÇÃO DO ENEM", cellHeaderText),
	)
}

func (w WebsiteItem) GetContent(_ int) core.Row {
	return row.New(5).Add(
		text.NewCol(1, fmt.Sprint(w.order)),
		text.NewCol(7, w.name),
		text.NewCol(1, fmt.Sprint(w.score)),
		text.NewCol(3, w.enrollmentID),
	)
}

type WebsiteRenderer struct{}

func (w *WebsiteRenderer) Render(courses []*CourseInfo) ([]core.Row, error) {
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

		var items []WebsiteItem

		rows = append(rows, text.NewRow(20, course.Quota, titleText))

		for i, registration := range course.Registrations {
			item := newWebsiteItem(i+1, registration.Candidate.Name, registration.CompositeScore, registration.EnrollmentID)
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

func (w *WebsiteRenderer) header(selection *SelectionInfo, period string) []core.Row {
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
