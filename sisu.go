// nolint: godox,nolintlint
// TODO: make some helper functions to avoid the repetition.
package main

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	repository2 "changeme/internal/repository"
	"changeme/repository"

	"changeme/types"

	"github.com/dimchansky/utfbom"
	"github.com/gocarina/gocsv"
	"go.uber.org/zap"
)

// FIXME: import is no longer called import.
// TODO: check if all of them are used.
var (
	ErrTxDB                          = errors.New("Erro na transação, contate o desenvolvedor")
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
	ErrNoFile                        = errors.New("Arquivo selecionado não existe")
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
	repo            *repository2.Repository
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
// FIXME: why is it so big and complex?
func (s *SISU) LoadSelection(path string, kind types.SelectionKind) error { //nolint: funlen, cyclop
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrNoFile
		}

		if os.IsPermission(err) {
			return ErrPermFile
		}

		s.l.Errorw("os open",
			"error", err,
		)

		return ErrOpenFile
	}

	status := "WAITING"
	if kind == types.ApprovedSelection {
		status = "APPROVED"
	}

	rawApplications, err := parseCSVApplicants(file, status)
	if err != nil {
		if errors.Is(err, gocsv.ErrEmptyCSVFile) {
			return ErrEmptyCSV
		}

		if errors.Is(err, gocsv.ErrUnmatchedStructTags) ||
			errors.Is(err, gocsv.ErrDoubleHeaderNames) ||
			errors.Is(err, gocsv.ErrNoStructTags) {
			return ErrMalformedCSV
		}

		s.l.Errorw("parse csv applications",
			"file", file,
			"status", status,
			"error", err,
		)

		return ErrParseFile
	}

	applications := make([]*types.Application, len(rawApplications))

	for i, csvApplicant := range rawApplications { //nolint: varnamelen
		application, err := csvApplicant.ToApplication(status)
		if err != nil {
			return fmt.Errorf("csv applicant to application: %w", err)
		}

		applications[i] = application
	}

	var selection types.Selection
	selection.Name = file.Name()
	selection.Kind = kind

	if len(rawApplications) > 0 {
		selection.Date = rawApplications[0].Date
		selection.Institution = rawApplications[0].Institution
		selection.Course = rawApplications[0].Course
	}

	if ok, err := s.ValidateSelection(&selection); !ok {
		return err
	}

	tx, err := s.repo.Begin()
	if err != nil {
		s.l.Errorw("db begin transaction", "err", err)

		return ErrTxDB
	}
	defer tx.Rollback()

	id, err := tx.SaveSelection(&selection)
	if err != nil {
		s.l.Errorw("db save selection", "err", err)

		return ErrTxDB
	}
	selection.ID = id

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

	if err := tx.Commit(); err != nil {
		s.l.Errorw("db commit transaction", "err", err)

		return ErrTxDB
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

func (s *SISU) WebsitePDF(rollcallID int64, periodID int64, file string) error {
	selection, err := s.selectionRepo.FindSelectionByKind(types.ApprovedSelection)
	if err != nil {
		return fmt.Errorf("find selection by kind: %w", err)
	}

	classes, err := s.classRepo.FindClassesByPeriodID(periodID)
	if err != nil {
		return fmt.Errorf("find class by period id: %w", err)
	}

	rollcall, err := s.rollcallRepo.FindRollcallByID(&rollcallID)
	if err != nil {
		return fmt.Errorf("find rollcall by id: %w", err)
	}

	pdf := PDF{
		selection:       selection,
		classes:         classes,
		rollcall:        rollcall,
		applicationRepo: &s.applicationRepo,
	}

	if err := pdf.ToWebsite(file); err != nil {
		return fmt.Errorf("pdf to website: %w", err)
	}

	return nil
}

func (s *SISU) EnrollmentPDF(rollcallID int64, periodID int64, file string) error {
	selection, err := s.selectionRepo.FindSelectionByKind(types.ApprovedSelection)
	if err != nil {
		return fmt.Errorf("find selection by kind: %w", err)
	}

	classes, err := s.classRepo.FindClassesByPeriodID(periodID)
	if err != nil {
		return fmt.Errorf("find class by period id: %w", err)
	}

	rollcall, err := s.rollcallRepo.FindRollcallByID(&rollcallID)
	if err != nil {
		return fmt.Errorf("find rollcall by id: %w", err)
	}

	pdf := PDF{
		selection:       selection,
		classes:         classes,
		rollcall:        rollcall,
		applicationRepo: &s.applicationRepo,
	}

	if err := pdf.ToEnrollment(file); err != nil {
		return fmt.Errorf("pdf to website: %w", err)
	}

	return nil
}

func (s *SISU) EmailPDF(rollcallID int64, periodID int64, file string) error {
	selection, err := s.selectionRepo.FindSelectionByKind(types.ApprovedSelection)
	if err != nil {
		return fmt.Errorf("find selection by kind: %w", err)
	}

	classes, err := s.classRepo.FindClassesByPeriodID(periodID)
	if err != nil {
		return fmt.Errorf("find class by period id: %w", err)
	}

	rollcall, err := s.rollcallRepo.FindRollcallByID(&rollcallID)
	if err != nil {
		return fmt.Errorf("find rollcall by id: %w", err)
	}

	pdf := PDF{
		selection:       selection,
		classes:         classes,
		rollcall:        rollcall,
		applicationRepo: &s.applicationRepo,
	}

	if err := pdf.ToEmail(file); err != nil {
		return fmt.Errorf("pdf to website: %w", err)
	}

	return nil
}

func (s *SISU) TeacherPDF(periodID int64, file string) error {
	selection, err := s.selectionRepo.FindSelectionByKind(types.ApprovedSelection)
	if err != nil {
		return fmt.Errorf("find selection by kind: %w", err)
	}

	classes, err := s.classRepo.FindClassesByPeriodID(periodID)
	if err != nil {
		return fmt.Errorf("find class by period id: %w", err)
	}

	pdf := PDF{
		selection:       selection,
		classes:         classes,
		rollcall:        nil,
		applicationRepo: &s.applicationRepo,
	}

	if err := pdf.ToTeacher(file); err != nil {
		return fmt.Errorf("pdf to website: %w", err)
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

type csvApplication struct {
	Stage                string `csv:"NU_ETAPA"`
	SchedulePeriod       string `csv:"DS_TURNO"`
	Seats                string `csv:"QT_VAGAS_CONCORRENCIA"`
	EnrollmentID         string `csv:"CO_INSCRICAO_ENEM"`
	Name                 string `csv:"NO_INSCRITO"`
	CPF                  string `csv:"NU_CPF_INSCRITO"`
	Date                 string `csv:"DT_OPERACAO"`
	SocialName           string `csv:"NO_SOCIAL"`
	BirthDate            string `csv:"DT_NASCIMENTO"`
	Sex                  string `csv:"TP_SEXO"`
	MotherName           string `csv:"NO_MAE"`
	AddressLine          string `csv:"DS_LOGRADOURO"`
	HouseNumber          string `csv:"NU_ENDERECO"`
	AddressLine2         string `csv:"DS_COMPLEMENTO"`
	State                string `csv:"SG_UF_INSCRITO"`
	Municipality         string `csv:"NO_MUNICIPIO"`
	Neighborhood         string `csv:"NO_BAIRRO"`
	CEP                  string `csv:"NU_CEP"`
	Phone1               string `csv:"NU_FONE1"`
	Phone2               string `csv:"NU_FONE2"`
	Email                string `csv:"DS_EMAIL"`
	LanguagesScore       string `csv:"NU_NOTA_L"`
	HumanitiesScore      string `csv:"NU_NOTA_CH"`
	NaturalSciencesScore string `csv:"NU_NOTA_CN"`
	MathematicsScore     string `csv:"NU_NOTA_M"`
	EssayScore           string `csv:"NU_NOTA_R"`
	Option               string `csv:"ST_OPCAO"`
	Quota                string `csv:"NO_MODALIDADE_CONCORRENCIA"`
	CompositeScore       string `csv:"NU_NOTA_CANDIDATO"`
	MinimumScore         string `csv:"NU_NOTACORTE_CONCORRIDA"`
	Ranking              string `csv:"NU_CLASSIFICACAO"`
	Institution          string `csv:"SG_IES"`
	Course               string `csv:"NO_CURSO"`
}

// TODO: write test for this func.
func parseCSVApplicants(
	csvApplicationReader io.Reader,
	status string,
) ([]*csvApplication, error) {
	preparedCSVReader, err := prepareCSV(csvApplicationReader)
	if err != nil {
		return nil, fmt.Errorf("prepare CSV: %w", err)
	}

	// We use LazyQuotes because some fields are surrounded by quotes while
	// others don't. Same comment from prepareCSV() applies here, we have no clue
	// if this is like this at source.
	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		c := csv.NewReader(in)
		c.Comma = ';'
		c.LazyQuotes = true

		return c
	})

	csvApplicants := []*csvApplication{}

	if err := gocsv.Unmarshal(preparedCSVReader, &csvApplicants); err != nil {
		return nil, fmt.Errorf("CSV unmarshal: %w", err)
	}

	return csvApplicants, nil
}

func (a *csvApplication) ToApplication(status string) (*types.Application, error) {
	minimumScore, err := scoreStringToFloat(a.MinimumScore)
	if err != nil {
		return nil, fmt.Errorf("minimum score string to float: %w", err)
	}

	seats, err := strconv.Atoi(a.Seats)
	if err != nil {
		return nil, fmt.Errorf("seats string to int: %w", err)
	}

	option, err := strconv.Atoi(a.Option)
	if err != nil {
		return nil, fmt.Errorf("option string to int: %w", err)
	}

	languagesScore, err := scoreStringToFloat(a.LanguagesScore)
	if err != nil {
		return nil, fmt.Errorf("languages score string to float: %w", err)
	}

	humanitiesScore, err := scoreStringToFloat(a.HumanitiesScore)
	if err != nil {
		return nil, fmt.Errorf("humanities score string to float: %w", err)
	}

	naturalSciencesScore, err := scoreStringToFloat(a.NaturalSciencesScore)
	if err != nil {
		return nil, fmt.Errorf("natural sciences score string to float: %w", err)
	}

	mathematicsScore, err := scoreStringToFloat(a.MathematicsScore)
	if err != nil {
		return nil, fmt.Errorf("mathematics score string to float: %w", err)
	}

	essayScore, err := scoreStringToFloat(a.EssayScore)
	if err != nil {
		return nil, fmt.Errorf("essay score string to float: %w", err)
	}

	compositeScore, err := scoreStringToFloat(a.CompositeScore)
	if err != nil {
		return nil, fmt.Errorf("composite score string to float: %w", err)
	}

	ranking, err := strconv.Atoi(a.Ranking)
	if err != nil {
		return nil, fmt.Errorf("ranking string to int: %w", err)
	}

	return &types.Application{
		ID:                   0,
		Status:               status,
		EnrollmentID:         a.EnrollmentID,
		Option:               option,
		LanguagesScore:       languagesScore,
		HumanitiesScore:      humanitiesScore,
		NaturalSciencesScore: naturalSciencesScore,
		MathematicsScore:     mathematicsScore,
		EssayScore:           essayScore,
		CompositeScore:       compositeScore,
		Ranking:              ranking,
		Applicant: types.Applicant{
			CPF:          a.CPF,
			Name:         a.Name,
			SocialName:   a.SocialName,
			BirthDate:    a.BirthDate,
			Sex:          a.Sex,
			MotherName:   a.MotherName,
			AddressLine:  a.AddressLine,
			AddressLine2: a.AddressLine2,
			HouseNumber:  a.HouseNumber,
			Neighborhood: a.Neighborhood,
			Municipality: a.Municipality,
			State:        a.State,
			CEP:          a.CEP,
			Email:        a.Email,
			Phone1:       a.Phone1,
			Phone2:       a.Phone2,
		},
		Class: types.Class{
			ID: 0,
			Period: types.Period{
				Name: a.SchedulePeriod,
			},
			Quota: types.Quota{
				Name: a.Quota,
			},
			Seats:        seats,
			MinimumScore: minimumScore,
		},
	}, nil
}

func scoreStringToFloat(s string) (float64, error) {
	score, err := strconv.ParseFloat(strings.ReplaceAll(s, ",", "."), 64)
	if err != nil {
		return 0, fmt.Errorf("parse float: %w", err)
	}

	return score, nil
}

/*
SISU's CSV file has peculiarities that require handling before parsing:

  - First character in some cases is a UTF-8 BOM, hence the SkipOnly call;

  - Line breaks are encoded using CR instead of LF (go expects LF);

  - The first line is separated by commas instead of semicolons (as the rest
    of the file is).

Note that we're not sure if the file comes like this or if all of this is made
by someones excel before the CSV gets here, but we have no direct access to
the source file, so we have to deal with it.
*/
func prepareCSV(reader io.Reader) (io.Reader, error) {
	rawCSV, err := io.ReadAll(utfbom.SkipOnly(reader))
	if err != nil {
		return nil, fmt.Errorf("read all: %w", err)
	}

	lines := strings.Split(string(rawCSV), "\r")
	lines[0] = strings.ReplaceAll(lines[0], ",", ";")
	preparedReader := bytes.NewReader([]byte(strings.Join(lines, "\n")))

	return preparedReader, nil
}
