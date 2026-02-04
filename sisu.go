// nolint: godox,nolintlint
// TODO: make some helper functions to avoid the repetition.
package main

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/baldugus/sisu/csvparser"
	"github.com/baldugus/sisu/pdfbuilder"
	"github.com/baldugus/sisu/repository"
	"github.com/baldugus/sisu/types"

	"go.uber.org/zap"
)

// FIXME: import is no longer called import.
// TODO: check if all of them are used.
var (
	ErrAlreadyImported               = errors.New("Arquivo deste tipo já foi importado")
	ErrApplicationNotFoundOrNotAvail = errors.New("Inscrição não encontrada ou não disponível para esta ação")
	ErrCloseFail                     = errors.New("Ainda existem alunos pendentes nesta chamada")
	ErrCloseRollCall                 = errors.New("Erro desconhecido ao atualizar chamada, mais detalhes nos logs")
	ErrCreateRollCall                = errors.New("Existe uma chamada mais antiga que está aberta ou seleção de interessados não foi importada")
	ErrDoneRollCallExists            = errors.New("Existem chamadas finalizada")
	ErrEmptyCSV                      = errors.New("Arquivo csv selecionado está vazio")
	ErrGenericDB                     = errors.New("Erro desconhecido ao consultar o banco de dados, mais detalhes nos logs")
	ErrInterestedSelectionExists     = errors.New("Seleção de interessados existe")
	ErrMalformedCSV                  = errors.New("Arquivo csv mal formatado ou inválido")
	ErrMissingFirstImport            = errors.New("Arquivo de aprovados deve ser importado antes do arquivo de alunos que manifestaram interesse")
	ErrMissingRollCall               = errors.New("Chamada não encontrada")
	ErrNoApplications                = errors.New("Não há inscrições deste tipo")
	ErrNoRollCalls                   = errors.New("Não existem chamadas criadas")
	ErrNoSelection                   = errors.New("Não há seleção deste tipo importada")
	ErrOpenFile                      = errors.New("Erro desconhecido ao abrir arquivo, mais detalhes nos logs")
	ErrOpenRollCall                  = errors.New("Existe uma chamada mais recente que está aberta")
	ErrParseFile                     = errors.New("Erro desconhecido ao processar arquivo CSV, mais detalhes nos logs")
	ErrPermFile                      = errors.New("Sem permissão para abrir o arquivo selecionado")
	ErrRollCallNotFoundOrNotAvail    = errors.New("Chamada não encontrada ou não disponível para esta ação")
	ErrSaveFile                      = errors.New("Erro desconhecido ao salvar inscrições, mais detalhes nos logs")
	ErrUpdatedApplicationExists      = errors.New("Existem inscrições atualizadas")
	ErrNoApplicationsToFill          = errors.New("Todas as vagas já foram preenchidas")
)

type SISU struct {
	applicantRepo   repository.ApplicantRepository
	applicationRepo repository.ApplicationRepository
	classRepo       repository.ClassRepository
	rollcallRepo    repository.RollcallRepository
	selectionRepo   repository.SelectionRepository
	service         *RepositoryService
	l               *zap.SugaredLogger
}

// TODO: write test for this func
// TODO: validate file name and if its the same semester
func (s *SISU) LoadSelection(path string, kind types.SelectionKind) error { //nolint: funlen, cyclop
	s.l.Infow("loading selection", "path", path, "kind", kind)

	// TODO: types/selection.go L12-15
	status := "WAITING"
	if kind == types.ApprovedSelection {
		status = "APPROVED"
	}

	// Should it be current timestamp or user input?
	// Maybe ask in a pop-up for the year and semester
	date := time.Now().Format(time.DateTime)

	parsedCsv, err := csvparser.ParseFile(path)
	if err != nil {
		s.l.Errorw("failed to parse CSV file",
			"path", path,
			"kind", kind,
			"error", err)
		return fmt.Errorf("parse csv: %w", err)
	}

	applications, err := parsedCsv.ToApplications(status)
	if err != nil {
		s.l.Errorw("failed to convert parsed csv file",
			"path", path,
			"kind", kind,
			"error", err)
		return fmt.Errorf("convert csv to applications: %w", err)
	}

	selection := parsedCsv.ToSelection(kind, date)

	// UNTOUCHED AREA BELOW NEEDS REFACTORING AND PROPER LOGGING

	if err := s.selectionRepo.Begin(); err != nil {
		return fmt.Errorf("begin: %w", err)
	}
	defer s.selectionRepo.Rollback()

	if err := s.selectionRepo.CreateSelection(&selection); err != nil {
		return fmt.Errorf("create selection: %w", err)
	}

	var rollcall types.Rollcall
	if selection.Kind == types.ApprovedSelection {
		r, err := s.rollcallRepo.CreateRollcall()
		if err != nil {
			return fmt.Errorf("create rollcall: %w", err)
		}

		rollcall = *r
	}

	for _, application := range applications {
		application.Rollcall = rollcall
		if err := s.classRepo.CreateClass(&application.Class); err != nil {
			return fmt.Errorf("create class: %w", err)
		}
		if err := s.applicantRepo.CreateApplicant(&application.Applicant); err != nil {
			return fmt.Errorf("create applicant: %w", err)
		}
		if err := s.applicationRepo.CreateApplication(application, selection.ID); err != nil {
			return fmt.Errorf("create application: %w", err)
		}
	}

	if err := s.selectionRepo.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}

// FIX: this with andre.
func (s *SISU) CreateRollCall() error {
	if err := s.rollcallRepo.Begin(); err != nil {
		return err
	}
	defer s.rollcallRepo.Rollback()

	rollcall, err := s.rollcallRepo.CreateRollcall()
	if err != nil {
		return err
	}

	if err := s.allocApplications(rollcall); err != nil {
		return err
	}

	if err := s.rollcallRepo.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *SISU) allocApplications(rollcall *types.Rollcall) error {
	var updated bool
	classes, err := s.classRepo.FindClasses()
	if err != nil {
		return fmt.Errorf("find classes :%w", err)
	}

	for _, class := range classes {
		status := "ENROLLED"
		filter := types.ApplicationsFilter{
			ClassID: &class.ID,
			Status:  &status,
		}

		applications, err := s.applicationRepo.FindApplications(filter)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("enrolled applications: %w", err)
		}

		seats := class.Seats - len(applications)

		ok, err := s.fillClass(seats, class, rollcall)
		if err != nil {
			return fmt.Errorf("fill class: %w", err)
		}

		if ok {
			updated = true
		}

	}

	if !updated {
		return ErrNoApplicationsToFill
	}
	return nil
}

type applicationsRanked []*types.Application

func (v applicationsRanked) Len() int           { return len(v) }
func (v applicationsRanked) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v applicationsRanked) Less(i, j int) bool { return v[i].Ranking < v[j].Ranking }

func (s *SISU) fillClass(seats int, class *types.Class, rollcall *types.Rollcall) (bool, error) {
	status := "WAITING"
	filter := types.ApplicationsFilter{
		ClassID: &class.ID,
		Status:  &status,
	}

	applications, err := s.applicationRepo.FindApplications(filter)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, fmt.Errorf("find applications: %w", err)
	}

	sort.Sort(applicationsRanked(applications))

	if len(applications) < seats {
		seats = len(applications)
	}

	applications = applications[:seats]

	if len(applications) == 0 {
		return false, nil
	}

	for _, application := range applications {
		status := "APPROVED"
		update := types.ApplicationUpdate{
			Status:     &status,
			RollcallID: rollcall.ID,
		}
		if _, err := s.applicationRepo.UpdateApplication(application.ID, update); err != nil {
			return false, fmt.Errorf("update application: %w", err)
		}
	}

	return true, nil
}

func (s *SISU) CloseRollCall(id int64) error {
	status := "DONE"
	update := types.RollcallUpdate{Status: &status}
	if _, err := s.rollcallRepo.UpdateRollcall(id, update); err != nil {
		return err
	}

	return nil
}

func (s *SISU) OpenRollCall(id int64) error {
	status := "CALLING"
	update := types.RollcallUpdate{Status: &status}
	if _, err := s.rollcallRepo.UpdateRollcall(id, update); err != nil {
		return err
	}
	return nil
}

func (s *SISU) EnrollApplication(id int64) error {
	status := "ENROLLED"
	update := types.ApplicationUpdate{Status: &status}

	if _, err := s.applicationRepo.UpdateApplication(id, update); err != nil {
		return err
	}

	return nil
}

func (s *SISU) ClearApplicationStatus(id int64) error {
	status := "APPROVED"
	update := types.ApplicationUpdate{Status: &status}

	if _, err := s.applicationRepo.UpdateApplication(id, update); err != nil {
		return err
	}

	return nil
}

func (s *SISU) DeleteRollcall(id int64) error {
	if err := s.rollcallRepo.Begin(); err != nil {
		return fmt.Errorf("rollcall repo begin: %w", err)
	}
	defer s.rollcallRepo.Rollback()

	applicationFilter := types.ApplicationsFilter{
		RollcallID: &id,
	}

	applications, err := s.applicationRepo.FindApplications(applicationFilter)
	if err != nil {
		return fmt.Errorf("find applications: %w", err)
	}

	status := "WAITING"
	var rollcallID int64 = 0
	applicationUpdate := types.ApplicationUpdate{
		Status:     &status,
		RollcallID: &rollcallID,
	}

	for _, application := range applications {
		if _, err := s.applicationRepo.UpdateApplication(application.ID, applicationUpdate); err != nil {
			return fmt.Errorf("update application: %w", err)
		}
	}

	if err := s.rollcallRepo.DeleteRollcall(id); err != nil {
		return fmt.Errorf("delete rollcall: %w", err)
	}

	if err := s.rollcallRepo.Commit(); err != nil {
		return fmt.Errorf("rollcall repo commit: %w", err)
	}

	return nil
}

func (s *SISU) AbsentApplication(id int64) error {
	status := "ABSENT"
	update := types.ApplicationUpdate{Status: &status}

	if _, err := s.applicationRepo.UpdateApplication(id, update); err != nil {
		return err
	}

	return nil
}

func (s *SISU) FetchApprovedSelection() (*types.Selection, error) {
	selection, err := s.selectionRepo.FindSelectionByKind(types.ApprovedSelection)
	if err != nil {
		return nil, err
	}

	filter := types.ApplicationsFilter{
		SelectionID: &selection.ID,
	}

	applications, err := s.applicationRepo.FindApplications(filter)
	if err != nil {
		return nil, err
	}

	selection.Applications = applications

	return selection, nil
}

func (s *SISU) FetchInterestedSelection() (*types.Selection, error) {
	selection, err := s.selectionRepo.FindSelectionByKind(types.InterestedSelection)
	if err != nil {
		return nil, fmt.Errorf("find selection by kind: %w", err)
	}

	filter := types.ApplicationsFilter{
		SelectionID: &selection.ID,
	}

	applications, err := s.applicationRepo.FindApplications(filter)
	if err != nil {
		return nil, fmt.Errorf("find application by selection: %w", err)
	}

	selection.Applications = applications

	return selection, nil
}

func (s *SISU) FetchRollcalls() ([]*types.Rollcall, error) {
	rollcalls, err := s.rollcallRepo.FindRollcalls(types.RollcallsFilter{})
	if err != nil {
		return nil, err
	}

	return rollcalls, nil
}

func (s *SISU) FetchPeriods() ([]*types.Period, error) {
	periods, err := s.classRepo.FindPeriods()
	if err != nil {
		return nil, err
	}

	return periods, nil
}

func (s *SISU) FetchRollcallNumber(id int64) (int64, error) {
	rollcall, err := s.rollcallRepo.FindRollcallByID(&id)
	if err != nil {
		return 0, err
	}

	return rollcall.Number, nil
}

func (s *SISU) FetchApplicationsByRollCall(id int64) ([]*types.Application, error) {
	filter := types.ApplicationsFilter{
		RollcallID: &id,
	}

	applications, err := s.applicationRepo.FindApplications(filter)
	if err != nil {
		return nil, err
	}

	return applications, nil
}

func (s *SISU) FetchSelection(selection types.SelectionKind) (*types.Selection, error) {
	switch selection {
	case types.ApprovedSelection:
		return s.FetchApprovedSelection()
	case types.InterestedSelection:
		return s.FetchInterestedSelection()
	}

	// TODO
	return nil, errors.New("TODO")
}

func (s *SISU) FetchPeriodName(periodID int64) (string, error) {
	period, err := s.classRepo.FindPeriodByID(periodID)
	if err != nil {
		return "", fmt.Errorf("find period by id: %w", err)
	}

	return period.Name, nil
}

func (s *SISU) CreatePDFBuilder(period string, builderClasses []*pdfbuilder.ClassInfo, waitlistNum int64, file string) (*pdfbuilder.Builder, error) {
	selection, err := s.selectionRepo.FindSelectionByKind(types.ApprovedSelection)
	if err != nil {
		return nil, fmt.Errorf("find selection by kind: %w", err)
	}

	date, err := selection.ParseDate()
	if err != nil {
		return nil, fmt.Errorf("parse date: %w", err)
	}

	year := fmt.Sprintf("%d", date.Year())
	semester := "1"
	if int(date.Month()) >= 6 {
		semester = "2"
	}

	builderSelection := pdfbuilder.NewSelectionInfo(selection, year, semester, waitlistNum)

	return &pdfbuilder.Builder{
		Period:    period,
		Classes:   builderClasses,
		Selection: builderSelection,
	}, nil
}

func (s *SISU) WebsitePDF(rollcallID int64, periodID int64, file string) error {
	rollcall, err := s.rollcallRepo.FindRollcallByID(&rollcallID)
	if err != nil {
		return fmt.Errorf("find rollcall by id: %w", err)
	}

	classes, err := s.classRepo.FindClassesByPeriodID(periodID)
	if err != nil {
		return fmt.Errorf("find class by period id: %w", err)
	}

	builderClasses := make([]*pdfbuilder.ClassInfo, len(classes))
	for i, class := range classes {
		filter := types.ApplicationsFilter{
			RollcallID: rollcall.ID,
			ClassID:    &class.ID,
		}

		applications, err := s.applicationRepo.FindApplications(filter)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				applications = []*types.Application{}
			} else {
				return fmt.Errorf("find applications: %w", err)
			}
		}

		builderClasses[i] = pdfbuilder.NewClassInfo(class.Quota.Name, applications)
	}

	if len(classes) == 0 {
		return fmt.Errorf("no classes found")
	}

	period := classes[0].Period.Name

	builder, err := s.CreatePDFBuilder(period, builderClasses, rollcall.Number, file)
	if err != nil {
		return fmt.Errorf("create pdf builder: %w", err)
	}

	if err := builder.BuildWebsitePdf(file); err != nil {
		return fmt.Errorf("pdf to website: %w", err)
	}

	return nil
}

func (s *SISU) EnrollmentPDF(rollcallID int64, periodID int64, file string) error {
	rollcall, err := s.rollcallRepo.FindRollcallByID(&rollcallID)
	if err != nil {
		return fmt.Errorf("find rollcall by id: %w", err)
	}

	classes, err := s.classRepo.FindClassesByPeriodID(periodID)
	if err != nil {
		return fmt.Errorf("find class by period id: %w", err)
	}

	builderClasses := make([]*pdfbuilder.ClassInfo, len(classes))
	for i, class := range classes {
		filter := types.ApplicationsFilter{
			RollcallID: rollcall.ID,
			ClassID:    &class.ID,
		}

		applications, err := s.applicationRepo.FindApplications(filter)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				applications = []*types.Application{}
			} else {
				return fmt.Errorf("find applications: %w", err)
			}
		}

		builderClasses[i] = pdfbuilder.NewClassInfo(class.Quota.Name, applications)
	}

	if len(classes) == 0 {
		return fmt.Errorf("no classes found")
	}

	period := classes[0].Period.Name

	builder, err := s.CreatePDFBuilder(period, builderClasses, rollcall.Number, file)
	if err != nil {
		return fmt.Errorf("create pdf builder: %w", err)
	}

	if err := builder.BuildEnrollmentPdf(file); err != nil {
		return fmt.Errorf("pdf to enrollment: %w", err)
	}

	return nil
}

func (s *SISU) EmailPDF(rollcallID int64, periodID int64, file string) error {
	rollcall, err := s.rollcallRepo.FindRollcallByID(&rollcallID)
	if err != nil {
		return fmt.Errorf("find rollcall by id: %w", err)
	}

	classes, err := s.classRepo.FindClassesByPeriodID(periodID)
	if err != nil {
		return fmt.Errorf("find class by period id: %w", err)
	}

	builderClasses := make([]*pdfbuilder.ClassInfo, len(classes))
	for i, class := range classes {
		filter := types.ApplicationsFilter{
			RollcallID: rollcall.ID,
			ClassID:    &class.ID,
		}

		applications, err := s.applicationRepo.FindApplications(filter)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				applications = []*types.Application{}
			} else {
				return fmt.Errorf("find applications: %w", err)
			}
		}

		builderClasses[i] = pdfbuilder.NewClassInfo(class.Quota.Name, applications)
	}

	if len(classes) == 0 {
		return fmt.Errorf("no classes found")
	}

	period := classes[0].Period.Name

	builder, err := s.CreatePDFBuilder(period, builderClasses, rollcall.Number, file)
	if err != nil {
		return fmt.Errorf("create pdf builder: %w", err)
	}

	if err := builder.BuildEmailPdf(file); err != nil {
		return fmt.Errorf("pdf to email: %w", err)
	}

	return nil
}

func (s *SISU) TeacherPDF(periodID int64, file string) error {
	classes, err := s.classRepo.FindClassesByPeriodID(periodID)
	if err != nil {
		return fmt.Errorf("find class by period id: %w", err)
	}

	var builderClasses []*pdfbuilder.ClassInfo
	for _, class := range classes {
		status := "ENROLLED"
		filter := types.ApplicationsFilter{
			ClassID: &class.ID,
			Status:  &status,
		}

		applications, err := s.applicationRepo.FindApplications(filter)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue
			} else {
				return fmt.Errorf("find applications: %w", err)
			}
		}

		builderClasses = append(builderClasses, pdfbuilder.NewClassInfo(class.Quota.Name, applications))
	}

	if len(classes) == 0 {
		return fmt.Errorf("no classes found")
	}

	period := classes[0].Period.Name

	builder, err := s.CreatePDFBuilder(period, builderClasses, 0, file)
	if err != nil {
		return fmt.Errorf("create pdf builder: %w", err)
	}

	if err := builder.BuildTeacherPdf(file); err != nil {
		return fmt.Errorf("pdf to teacher: %w", err)
	}

	return nil
}

func (s *SISU) DeleteSelection(kind types.SelectionKind) error {
	if err := s.selectionRepo.Begin(); err != nil {
		return fmt.Errorf("begin: %w", err)
	}
	defer s.selectionRepo.Rollback()

	selection, err := s.FetchSelection(kind)
	if err != nil {
		return fmt.Errorf("fetch selection: %w", err)
	}

	for _, a := range selection.Applications {
		if err := s.applicantRepo.DeleteApplicant(a.Applicant.ID); err != nil {
			return fmt.Errorf("delete applicant: %w", err)
		}
		if err := s.applicationRepo.DeleteApplication(a.ID); err != nil {
			return fmt.Errorf("delete application: %w", err)
		}
	}

	if kind == types.ApprovedSelection {
		if err := s.rollcallRepo.DeleteRollcall(1); err != nil {
			return fmt.Errorf("delete rollcall: %w", err)
		}
	}

	if err := s.selectionRepo.DeleteSelection(selection.ID); err != nil {
		return fmt.Errorf("delete selection: %w", err)
	}

	if err := s.selectionRepo.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}

func (s *SISU) ExportCSV(file string) error {
	status := "ENROLLED"
	filter := types.ApplicationsFilter{
		Status: &status,
	}

	applications, err := s.applicationRepo.FindApplications(filter)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("find applications: %w", err)
	}

	return ApplicationsToCSV(applications, file)
}

func (s *SISU) Backup(file string) error {
	return s.service.Backup(file)
}

func (s *SISU) Restore(file string) error {
	return s.service.Restore(file)
}

func (s *SISU) Destroy() error {
	err := s.rollcallRepo.Close()
	if err != nil {
		return fmt.Errorf("close: %w", err)
	}

	s.service.Destroy()
	return nil
}
