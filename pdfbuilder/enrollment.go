package pdfbuilder

import (
	"fmt"
	"sort"
	"strings"

	"github.com/johnfercher/maroto/v2/pkg/components/list"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

type EnrollmentItem struct {
	order        int
	name         string
	score        float64
	enrollmentID string
}

func newEnrollmentItem(order int, name string, score float64, enrollmentID string) EnrollmentItem {
	return EnrollmentItem{
		order:        order,
		name:         name,
		score:        score,
		enrollmentID: enrollmentID,
	}
}

func (e EnrollmentItem) GetContent(_ int) core.Row {
	return row.New(5).Add(
		text.NewCol(1, fmt.Sprint(e.order)),
		text.NewCol(5, e.name),
		text.NewCol(1, fmt.Sprintf("%.2f", e.score)),
		text.NewCol(2, e.enrollmentID),
		text.NewCol(5, strings.Repeat("_", 30)),
	)
}

func (e EnrollmentItem) GetHeader() core.Row {
	var cellHeaderText props.Text
	cellHeaderText.Style = fontstyle.Bold

	return row.New(10).Add(
		text.NewCol(1, "ORDEM", cellHeaderText),
		text.NewCol(5, "NOME DO CANDIDATO", cellHeaderText),
		text.NewCol(1, "NOTA", cellHeaderText),
		text.NewCol(2, "INSCRIÇÃO DO ENEM", cellHeaderText),
		text.NewCol(5, "ASSINATURA", cellHeaderText),
	)
}

type EnrollmentRenderer struct{}

func (e *EnrollmentRenderer) Render(classes []*ClassInfo) ([]core.Row, error) {
	var titleText props.Text
	titleText.Top = 10
	titleText.Style = fontstyle.Bold
	titleText.Align = align.Center

	var rows []core.Row

	for _, class := range classes {
		sort.Sort(class.Applications)

		var items []EnrollmentItem

		rows = append(rows, text.NewRow(20, class.Quota, titleText))

		for i, application := range class.Applications {
			item := newEnrollmentItem(i+1, application.Applicant.Name, application.CompositeScore, application.EnrollmentID)
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

func (e *EnrollmentRenderer) header(selection *SelectionInfo, period string) []core.Row {
	var call strings.Builder

	call.WriteString(string(selection.Kind))
	if selection.WaitlistNum > 0 {
		call.WriteString(fmt.Sprintf(" %d", selection.WaitlistNum))
	}

	call.WriteString(fmt.Sprintf(" - %s - ", selection.Year))
	call.WriteString(fmt.Sprintf("%so. Semestre", selection.Semester))
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
