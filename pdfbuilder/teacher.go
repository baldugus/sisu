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

type TeacherItem struct {
	order int
	name  string
}

func newTeacherItem(order int, name string) TeacherItem {
	return TeacherItem{
		order: order,
		name:  name,
	}
}

func (t TeacherItem) GetContent(_ int) core.Row {
	return row.New(5).Add(
		text.NewCol(1, fmt.Sprint(t.order)),
		text.NewCol(5, t.name),
		text.NewCol(7, strings.Repeat("_", 50)),
	)
}

func (t TeacherItem) GetHeader() core.Row {
	var cellHeaderText props.Text
	cellHeaderText.Style = fontstyle.Bold

	return row.New(10).Add(
		text.NewCol(1, "Seq.", cellHeaderText),
		text.NewCol(5, "Nome", cellHeaderText),
		text.NewCol(7, "Assinatura", cellHeaderText),
	)
}

type TeacherRenderer struct{}

func (t *TeacherRenderer) Render(courses []*CourseInfo) ([]core.Row, error) {
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

		var items []TeacherItem

		rows = append(rows, text.NewRow(20, course.Quota, titleText))

		for i, registration := range course.Registrations {
			item := newTeacherItem(i+1, registration.Candidate.Name)
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

func (t *TeacherRenderer) header(selection *SelectionInfo, period string) []core.Row {
	var header strings.Builder

	header.WriteString(fmt.Sprintf("LISTA DE PRESENÇA - %s %d", selection.Institution, selection.Year))
	header.WriteString(fmt.Sprintf(".%d", selection.Semester))
	period = fmt.Sprintf("TURNO: %s", period)

	var textStyle props.Text
	textStyle.Size = 11
	textStyle.Style = fontstyle.Bold
	textStyle.Align = align.Center

	rows := []core.Row{
		text.NewRow(7, header.String(), textStyle),
		text.NewRow(7, period, textStyle),
	}

	textStyle.Align = align.Left

	rows = append(rows, text.NewRow(7, "DISCIPLINA:", textStyle))
	rows = append(rows, text.NewRow(7, "PROFESSOR:", textStyle))
	rows = append(rows, text.NewRow(7, fmt.Sprintf("DATA: ____/____/%d", selection.Year), textStyle))

	return rows
}
