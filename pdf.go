package main

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"changeme/types"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/list"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/consts/orientation"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

type PDF struct {
	selection       *types.Selection
	rollcall        *types.Rollcall
	classes         []*types.Class
	applicationRepo types.ApplicationRepository
}

type WebsiteListItem struct {
	Order        int
	Name         string
	Score        float64
	EnrollmentID string
}

func (w WebsiteListItem) GetContent(_ int) core.Row {
	return row.New(5).Add(
		text.NewCol(1, fmt.Sprint(w.Order)),
		text.NewCol(7, w.Name),
		text.NewCol(1, fmt.Sprintf("%.2f", w.Score)),
		text.NewCol(3, w.EnrollmentID),
	)
}

func (w WebsiteListItem) GetHeader() core.Row {
	var cellHeaderText props.Text
	cellHeaderText.Style = fontstyle.Bold

	return row.New(10).Add(
		text.NewCol(1, "ORDEM", cellHeaderText),
		text.NewCol(7, "NOME DO CANDIDATO", cellHeaderText),
		text.NewCol(1, "NOTA", cellHeaderText),
		text.NewCol(3, "INSCRIÇÃO DO ENEM", cellHeaderText),
	)
}

type EnrollmentListItem struct {
	Order        int
	Name         string
	Score        float64
	EnrollmentID string
}

func (e EnrollmentListItem) GetContent(_ int) core.Row {
	return row.New(5).Add(
		text.NewCol(1, fmt.Sprint(e.Order)),
		text.NewCol(5, e.Name),
		text.NewCol(1, fmt.Sprintf("%.2f", e.Score)),
		text.NewCol(2, e.EnrollmentID),
		text.NewCol(5, strings.Repeat("_", 30)),
	)
}

func (e EnrollmentListItem) GetHeader() core.Row {
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

type EmailsListItem struct {
	Order        int
	Name         string
	Score        float64
	EnrollmentID string
	Email        string
}

func (e EmailsListItem) GetContent(_ int) core.Row {
	return row.New(5).Add(
		text.NewCol(1, fmt.Sprint(e.Order)),
		text.NewCol(5, e.Name),
		text.NewCol(1, fmt.Sprintf("%.2f", e.Score)),
		text.NewCol(2, e.EnrollmentID),
		text.NewCol(5, e.Email),
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

type TeacherListItem struct {
	Order int
	Name  string
}

func (t TeacherListItem) GetContent(_ int) core.Row {
	return row.New(5).Add(
		text.NewCol(1, fmt.Sprint(t.Order)),
		text.NewCol(5, t.Name),
		text.NewCol(7, strings.Repeat("_", 50)),
	)
}

func (t TeacherListItem) GetHeader() core.Row {
	var cellHeaderText props.Text
	cellHeaderText.Style = fontstyle.Bold

	return row.New(10).Add(
		text.NewCol(1, "Seq.", cellHeaderText),
		text.NewCol(5, "Nome", cellHeaderText),
		text.NewCol(7, "Assinatura", cellHeaderText),
	)
}

func header(rollcallNumber int64, date time.Time, period string, institution string, course string) []core.Row {
	var timeHeader strings.Builder

	timeHeader.WriteString(fmt.Sprintf("LISTA DE ESPERA %d", rollcallNumber))

	if rollcallNumber == 1 {
		timeHeader.Reset()
		timeHeader.WriteString("CHAMADA REGULAR")
	}

	timeHeader.WriteString(fmt.Sprintf(" - %s - ", date.Format("2006")))

	semester := "1"

	if int(date.Month()) >= 6 {
		semester = "2"
	}

	timeHeader.WriteString(fmt.Sprintf("%so. Semestre", semester))

	periodHeader := fmt.Sprintf("TURNO: %s", period)

	var headerText props.Text
	headerText.Size = 11
	headerText.Style = fontstyle.Bold
	headerText.Align = align.Center

	return []core.Row{
		text.NewRow(7, institution, headerText),
		text.NewRow(7, timeHeader.String(), headerText),
		text.NewRow(7, course, headerText),
		text.NewRow(7, periodHeader, headerText),
	}
}

func applicationsToWebsite(classes []*types.Class, rollcallID *int64, applicationsRepo types.ApplicationRepository) ([]core.Row, error) {
	var titleText props.Text
	titleText.Top = 10
	titleText.Style = fontstyle.Bold
	titleText.Align = align.Center

	var rows []core.Row

	for _, class := range classes {
		filter := types.ApplicationsFilter{
			RollcallID: rollcallID,
			ClassID:    &class.ID,
		}

		applications, err := applicationsRepo.FindApplications(filter)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}

			return nil, fmt.Errorf("find applications: %w", err)
		}

		sort.Sort(applicationsRanked(applications))

		var items []WebsiteListItem

		rows = append(rows, text.NewRow(20, class.Quota.Name, titleText))

		for i, application := range applications {
			var item WebsiteListItem
			item.Order = i + 1
			item.Name = application.Applicant.Name
			item.Score = application.CompositeScore
			item.EnrollmentID = application.EnrollmentID

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

func applicationsToEnrollment(classes []*types.Class, rollcallID *int64, applicationsRepo types.ApplicationRepository) ([]core.Row, error) {
	var titleText props.Text
	titleText.Top = 10
	titleText.Style = fontstyle.Bold
	titleText.Align = align.Center

	var rows []core.Row

	for _, class := range classes {
		filter := types.ApplicationsFilter{
			RollcallID: rollcallID,
			ClassID:    &class.ID,
		}

		applications, err := applicationsRepo.FindApplications(filter)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}

			return nil, fmt.Errorf("find applications: %w", err)
		}

		sort.Sort(applicationsRanked(applications))

		var items []EnrollmentListItem

		rows = append(rows, text.NewRow(20, class.Quota.Name, titleText))

		for i, application := range applications {
			var item EnrollmentListItem
			item.Order = i + 1
			item.Name = application.Applicant.Name
			item.Score = application.CompositeScore
			item.EnrollmentID = application.EnrollmentID

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

func applicationsToEmail(classes []*types.Class, rollcallID *int64, applicationsRepo types.ApplicationRepository) ([]core.Row, error) {
	var titleText props.Text
	titleText.Top = 10
	titleText.Style = fontstyle.Bold
	titleText.Align = align.Center

	var rows []core.Row

	for _, class := range classes {
		filter := types.ApplicationsFilter{
			RollcallID: rollcallID,
			ClassID:    &class.ID,
		}

		applications, err := applicationsRepo.FindApplications(filter)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}

			return nil, fmt.Errorf("find applications: %w", err)
		}

		sort.Sort(applicationsRanked(applications))

		var items []EmailsListItem

		rows = append(rows, text.NewRow(20, class.Quota.Name, titleText))

		for i, application := range applications {
			var item EmailsListItem
			item.Order = i + 1
			item.Name = application.Applicant.Name
			item.Score = application.CompositeScore
			item.EnrollmentID = application.EnrollmentID
			item.Email = application.Email

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

func applicationsToTeacher(classes []*types.Class, applicationsRepo types.ApplicationRepository) ([]core.Row, error) {
	var titleText props.Text
	titleText.Top = 10
	titleText.Style = fontstyle.Bold
	titleText.Align = align.Center

	var items []TeacherListItem

	for _, class := range classes {
		status := "ENROLLED"
		filter := types.ApplicationsFilter{
			ClassID: &class.ID,
			Status:  &status,
		}

		applications, err := applicationsRepo.FindApplications(filter)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}

			return nil, fmt.Errorf("find applications: %w", err)
		}

		for _, application := range applications {
			var item TeacherListItem
			item.Order = len(items) + 1
			item.Name = application.Applicant.Name
			items = append(items, item)
		}
	}

	sort.Sort(itemsToSort(items))

	builtRow, err := list.Build(items)
	if err != nil {
		return nil, fmt.Errorf("list build: %w", err)
	}

	return builtRow, nil
}

type itemsToSort []TeacherListItem

func (t itemsToSort) Len() int { return len(t) }
func (t itemsToSort) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
	t[i].Order, t[j].Order = t[j].Order, t[i].Order
}
func (t itemsToSort) Less(i, j int) bool { return t[i].Name < t[j].Name }

func (p *PDF) ToWebsite(file string) error {
	mrt := maroto.New()

	date, err := p.selection.ParseDate()
	if err != nil {
		return fmt.Errorf("selection parse date: %w", err)
	}

	if err := mrt.RegisterHeader(header(p.rollcall.Number, date, p.classes[0].Period.Name, p.selection.Institution, p.selection.Course)...); err != nil {
		return fmt.Errorf("register header: %w", err)
	}

	rows, err := applicationsToWebsite(p.classes, p.rollcall.ID, p.applicationRepo)
	if err != nil {
		return fmt.Errorf("applications to website: %w", err)
	}

	mrt.AddRows(rows...)

	pdf, err := mrt.Generate()
	if err != nil {
		return fmt.Errorf("maroto generate: %w", err)
	}

	if err := pdf.Save(file); err != nil {
		return fmt.Errorf("maroto save: %w", err)
	}

	return nil
}

func (p *PDF) ToEnrollment(file string) error {
	config := config.NewBuilder().WithOrientation(orientation.Horizontal)
	mrt := maroto.New(config.Build())

	date, err := p.selection.ParseDate()
	if err != nil {
		return fmt.Errorf("selection parse date: %w", err)
	}

	if err := mrt.RegisterHeader(header(p.rollcall.Number, date, p.classes[0].Period.Name, p.selection.Institution, p.selection.Course)...); err != nil {
		return fmt.Errorf("register header: %w", err)
	}

	rows, err := applicationsToEnrollment(p.classes, p.rollcall.ID, p.applicationRepo)
	if err != nil {
		return fmt.Errorf("applications to website: %w", err)
	}

	mrt.AddRows(rows...)

	pdf, err := mrt.Generate()
	if err != nil {
		return fmt.Errorf("maroto generate: %w", err)
	}

	if err := pdf.Save(file); err != nil {
		return fmt.Errorf("maroto save: %w", err)
	}

	return nil
}

func (p *PDF) ToEmail(file string) error {
	config := config.NewBuilder().WithOrientation(orientation.Horizontal)
	mrt := maroto.New(config.Build())

	date, err := p.selection.ParseDate()
	if err != nil {
		return fmt.Errorf("selection parse date: %w", err)
	}

	if err := mrt.RegisterHeader(header(p.rollcall.Number, date, p.classes[0].Period.Name, p.selection.Institution, p.selection.Course)...); err != nil {
		return fmt.Errorf("register header: %w", err)
	}

	rows, err := applicationsToEmail(p.classes, p.rollcall.ID, p.applicationRepo)
	if err != nil {
		return fmt.Errorf("applications to website: %w", err)
	}

	mrt.AddRows(rows...)

	pdf, err := mrt.Generate()
	if err != nil {
		return fmt.Errorf("maroto generate: %w", err)
	}

	if err := pdf.Save(file); err != nil {
		return fmt.Errorf("maroto save: %w", err)
	}

	return nil
}

func (p *PDF) ToTeacher(file string) error {
	config := config.NewBuilder().WithOrientation(orientation.Horizontal)
	mrt := maroto.New(config.Build())

	date, err := p.selection.ParseDate()
	if err != nil {
		return fmt.Errorf("selection parse date: %w", err)
	}

	var timeHeader strings.Builder

	timeHeader.WriteString(fmt.Sprintf("LISTA DE PRESENÇA - %s %s", p.selection.Institution, date.Format("2006")))

	semester := "1"

	if int(date.Month()) >= 6 {
		semester = "2"
	}

	timeHeader.WriteString(fmt.Sprintf(".%s", semester))

	periodHeader := fmt.Sprintf("TURNO: %s", p.classes[0].Period.Name)

	var headerText props.Text
	headerText.Size = 11
	headerText.Style = fontstyle.Bold
	headerText.Align = align.Center

	if err := mrt.RegisterHeader(
		text.NewRow(7, timeHeader.String(), headerText),
		text.NewRow(7, periodHeader, headerText),
	); err != nil {
		return fmt.Errorf("register header: %w", err)
	}

	var fillableText props.Text
	fillableText.Size = 11
	fillableText.Style = fontstyle.Bold
	fillableText.Align = align.Left

	mrt.AddRows(
		text.NewRow(7, "DISCIPLINA:", fillableText),
		text.NewRow(7, "PROFESSOR:", fillableText),
		text.NewRow(7, fmt.Sprintf("DATA: ____/____/%d", date.Year()), fillableText),
	)

	rows, err := applicationsToTeacher(p.classes, p.applicationRepo)
	if err != nil {
		return fmt.Errorf("applications to website: %w", err)
	}

	mrt.AddRows(rows...)

	pdf, err := mrt.Generate()
	if err != nil {
		return fmt.Errorf("maroto generate: %w", err)
	}

	if err := pdf.Save(file); err != nil {
		return fmt.Errorf("maroto save: %w", err)
	}

	return nil
}
