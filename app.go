// nolint: godox, nolintlint, staticcheck
// TODO: make some helper functions to avoid the repetition.
package main

import (
	"context"
	"errors"
	"os"
	"regexp"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.uber.org/zap"

	"github.com/baldugus/sisu/commands"
	"github.com/baldugus/sisu/csvparser"
	"github.com/baldugus/sisu/database"
	"github.com/baldugus/sisu/types"
)

// SISU holds the application dependencies.
type SISU struct {
	database *database.Database
}

// translateError maps known typed errors to a user-facing Portuguese message.
// Wails forwards the returned error's message to the rejected JS promise.
func translateError(err error) error {
	switch {
	case errors.As(err, &commands.ErrApprovedSelectionAlreadyExists{}):
		return errors.New("Não é possível importar a lista de aprovados quando já existe uma seleção.")
	case errors.As(err, &commands.ErrWaitlistSelectionRequiresApproved{}):
		return errors.New("Não é possível importar a lista de espera sem a lista de aprovados.")
	case errors.As(err, &csvparser.ErrFileNotFound{}):
		return errors.New("Arquivo selecionado não existe.")
	case errors.As(err, &csvparser.ErrPermissionDenied{}):
		return errors.New("Permissão negada para acessar o arquivo.")
	case errors.As(err, &csvparser.ErrFileOpen{}):
		return errors.New("Erro ao abrir o arquivo CSV. Contate o desenvolvedor.")
	case errors.As(err, &csvparser.ErrFileRead{}):
		return errors.New("Erro ao ler o arquivo CSV. Contate o desenvolvedor.")
	case errors.As(err, &csvparser.ErrFileEmpty{}):
		return errors.New("Arquivo CSV vazio.")
	case errors.As(err, &csvparser.ErrFileParse{}):
		return errors.New("Erro ao analisar o arquivo CSV. Contate o desenvolvedor.")
	case errors.As(err, &csvparser.ErrOddSeatsCount{}):
		return errors.New("O número total de vagas deve ser par para divisão entre semestres.")
	case errors.As(err, &csvparser.ErrLineMapping{}):
		return errors.New("Erro de validação em um campo do arquivo CSV. Contate o desenvolvedor.")
	case errors.As(err, &commands.ErrSelectionNotFound{}):
		return errors.New("Seleção não encontrada.")
	case errors.As(err, &commands.ErrCannotDeleteApprovedWithWaitlist{}):
		return errors.New("Não é possível excluir a lista de aprovados enquanto a lista de espera existir.")
	case errors.As(err, &commands.ErrSelectionHasModifiedRegistrations{}):
		return errors.New("Não é possível excluir a seleção com inscrições matriculadas ou ausentes.")
	case errors.As(err, &commands.ErrCallHasPendingRegistrations{}):
		return errors.New("Não é possível fechar a chamada com inscrições pendentes.")
	case errors.As(err, &commands.ErrRegistrationNotFound{}):
		return errors.New("Inscrição não encontrada.")
	case errors.As(err, &commands.ErrRegistrationNotInCall{}):
		return errors.New("Inscrição não está em uma chamada.")
	case errors.As(err, &commands.ErrCallNotOpen{}):
		return errors.New("A chamada não está aberta.")
	case errors.As(err, &commands.ErrInvalidStatusTransition{}):
		return errors.New("Transição de status inválida.")
	case errors.As(err, &commands.ErrNoSeatsAvailable{}):
		return errors.New("Não há vagas disponíveis no curso.")
	case errors.As(err, &commands.ErrOpenCallExists{}):
		return errors.New("Não é possível criar nova chamada enquanto outra está aberta.")
	case errors.As(err, &commands.ErrNoWaitlistedRegistrations{}):
		return errors.New("Não há inscrições na lista de espera.")
	case errors.As(err, &commands.ErrAllCoursesFull{}):
		return errors.New("Todas as vagas de todos os cursos estão ocupadas.")
	case errors.As(err, &commands.ErrCannotReopenCallWithLaterCalls{}):
		return errors.New("Não é possível reabrir a chamada enquanto existem chamadas posteriores.")
	case errors.As(err, &commands.ErrCannotDeleteWithMultipleCalls{}):
		return errors.New("Não é possível excluir a seleção quando existem múltiplas chamadas.")
	case errors.As(err, &commands.ErrCannotDeleteWithClosedFirstCall{}):
		return errors.New("Não é possível excluir a seleção quando a primeira chamada está fechada.")
	case errors.As(err, &commands.ErrCannotDeleteCallWithLaterCalls{}):
		return errors.New("Não é possível excluir a chamada enquanto existem chamadas posteriores.")
	case errors.As(err, &commands.ErrCannotDeleteClosedCall{}):
		return errors.New("Não é possível excluir uma chamada fechada.")
	default:
		zap.L().Error("unknown error", zap.Error(err))
		return errors.New("Erro desconhecido. Contate o desenvolvedor.")
	}
}

// App struct.
type App struct {
	ctx  context.Context
	sisu SISU
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	zap.L().Info("application window ready")
}

// OpenFileDialog opens a file picker dialog and returns the selected file path.
func (a *App) OpenFileDialog(title string, filterPattern string, filterDisplay string) (string, error) {
	homeDir, _ := os.UserHomeDir()
	return runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:            title,
		DefaultDirectory: homeDir,
		Filters: []runtime.FileFilter{
			{DisplayName: filterDisplay, Pattern: filterPattern},
		},
	})
}

// SaveFileDialog opens a save file dialog and returns the selected file path.
func (a *App) SaveFileDialog(title string, defaultFilename string, filterPattern string, filterDisplay string) (string, error) {
	homeDir, _ := os.UserHomeDir()
	return runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:            title,
		DefaultDirectory: homeDir,
		DefaultFilename:  defaultFilename,
		Filters: []runtime.FileFilter{
			{DisplayName: filterDisplay, Pattern: filterPattern},
		},
	})
}

func (a *App) LoadApprovedSelection(year int32, filePath string) error {
	cmd := commands.LoadSelectionCommand{
		Year:     year,
		FilePath: filePath,
		Kind:     types.SelectionKindApproved,
	}

	if err := cmd.Execute(a.sisu.database); err != nil {
		return translateError(err)
	}

	return nil
}

func (a *App) LoadWaitlistSelection(year int32, filePath string) error {
	cmd := commands.LoadSelectionCommand{
		Year:     year,
		FilePath: filePath,
		Kind:     types.SelectionKindWaitlist,
	}

	if err := cmd.Execute(a.sisu.database); err != nil {
		return translateError(err)
	}

	return nil
}

func (a *App) LoadInterestedSelection(year int32, filePath string) error {
	return a.LoadWaitlistSelection(year, filePath)
}

func (a *App) FetchApprovedSelection() (*types.Selection, error) {
	cmd := commands.FetchSelectionCommand{
		Kind: types.SelectionKindApproved,
	}

	selection, err := cmd.Execute(a.sisu.database)
	if err != nil {
		if errors.As(err, &database.ErrNotFound{}) {
			return nil, nil
		}
		return nil, translateError(err)
	}

	return selection, nil
}

func (a *App) FetchInterestedSelection() (*types.Selection, error) {
	cmd := commands.FetchSelectionCommand{
		Kind: types.SelectionKindWaitlist,
	}

	selection, err := cmd.Execute(a.sisu.database)
	if err != nil {
		if errors.As(err, &database.ErrNotFound{}) {
			return nil, nil
		}
		return nil, translateError(err)
	}

	return selection, nil
}

func (a *App) FetchRollCalls() ([]*types.Call, error) {
	cmd := commands.FetchCallsCommand{}

	calls, err := cmd.Execute(a.sisu.database)
	if err != nil {
		return nil, translateError(err)
	}

	return calls, nil
}

func (a *App) FetchSemesters() ([]*types.Semester, error) {
	cmd := commands.FetchSemestersCommand{}

	semesters, err := cmd.Execute(a.sisu.database)
	if err != nil {
		return nil, translateError(err)
	}

	return semesters, nil
}

func (a *App) DeleteApprovedSelection() error {
	cmd := commands.DeleteSelectionCommand{
		Kind: types.SelectionKindApproved,
	}

	if err := cmd.Execute(a.sisu.database); err != nil {
		return translateError(err)
	}

	return nil
}

func (a *App) DeleteInterestedSelection() error {
	cmd := commands.DeleteSelectionCommand{
		Kind: types.SelectionKindWaitlist,
	}

	if err := cmd.Execute(a.sisu.database); err != nil {
		return translateError(err)
	}

	return nil
}

func (a *App) FetchRegistrations() ([]*types.Registration, error) {
	cmd := commands.FetchRegistrationsCommand{}

	registrations, err := cmd.Execute(a.sisu.database)
	if err != nil {
		return nil, translateError(err)
	}

	return registrations, nil
}

func (a *App) FetchRegistrationsBySelectionID(selectionID int32) ([]*types.Registration, error) {
	cmd := commands.FetchRegistrationsCommand{
		SelectionID: &selectionID,
	}

	registrations, err := cmd.Execute(a.sisu.database)
	if err != nil {
		return nil, translateError(err)
	}

	return registrations, nil
}

func (a *App) FetchRegistrationsByCourseID(courseID int32) ([]*types.Registration, error) {
	cmd := commands.FetchRegistrationsCommand{
		CourseID: &courseID,
	}

	registrations, err := cmd.Execute(a.sisu.database)
	if err != nil {
		return nil, translateError(err)
	}

	return registrations, nil
}

func (a *App) FetchRegistrationsByCallID(callID int32) ([]*types.Registration, error) {
	cmd := commands.FetchRegistrationsCommand{
		CallID: &callID,
	}

	registrations, err := cmd.Execute(a.sisu.database)
	if err != nil {
		return nil, translateError(err)
	}

	return registrations, nil
}

func (a *App) FetchRegistration(registrationID int32) (*types.RegistrationDetail, error) {
	cmd := commands.FetchRegistrationCommand{
		RegistrationID: registrationID,
	}

	registration, err := cmd.Execute(a.sisu.database)
	if err != nil {
		return nil, translateError(err)
	}

	return registration, nil
}

func (a *App) CloseRollCall(id int32) error {
	cmd := commands.CloseCallCommand{
		ID: id,
	}

	if err := cmd.Execute(a.sisu.database); err != nil {
		return translateError(err)
	}

	return nil
}

func (a *App) EnrollRegistration(id int32) error {
	cmd := commands.UpdateRegistrationStatusCommand{
		RegistrationID: id,
		NewStatus:      types.RegistrationStatusEnrolled,
	}

	if err := cmd.Execute(a.sisu.database); err != nil {
		return translateError(err)
	}

	return nil
}

func (a *App) AbsentRegistration(id int32) error {
	cmd := commands.UpdateRegistrationStatusCommand{
		RegistrationID: id,
		NewStatus:      types.RegistrationStatusAbsent,
	}

	if err := cmd.Execute(a.sisu.database); err != nil {
		return translateError(err)
	}

	return nil
}

func (a *App) ClearRegistrationStatus(id int32) error {
	cmd := commands.UpdateRegistrationStatusCommand{
		RegistrationID: id,
		NewStatus:      types.RegistrationStatusApproved,
	}

	if err := cmd.Execute(a.sisu.database); err != nil {
		return translateError(err)
	}

	return nil
}

func (a *App) EnrollApplication(id int32) error {
	return a.EnrollRegistration(id)
}

func (a *App) AbsentApplication(id int32) error {
	return a.AbsentRegistration(id)
}

func (a *App) ClearApplicationStatus(id int32) error {
	return a.ClearRegistrationStatus(id)
}

func (a *App) FetchApplicationsByRollCall(callID int32) ([]*types.Registration, error) {
	return a.FetchRegistrationsByCallID(callID)
}

func (a *App) CreateRollCall(semesterID int32) error {
	cmd := commands.CreateCallCommand{SemesterID: semesterID}

	if err := cmd.Execute(a.sisu.database); err != nil {
		return translateError(err)
	}

	return nil
}

func (a *App) OpenCall(id int32) error {
	cmd := commands.OpenCallCommand{
		ID: id,
	}

	if err := cmd.Execute(a.sisu.database); err != nil {
		return translateError(err)
	}

	return nil
}

func (a *App) DeleteCall(id int32) error {
	cmd := commands.DeleteCallCommand{
		ID: id,
	}

	if err := cmd.Execute(a.sisu.database); err != nil {
		return translateError(err)
	}

	return nil
}

func (a *App) OpenRollCall(id int32) error {
	return a.OpenCall(id)
}

func (a *App) DeleteRollCall(id int32) error {
	return a.DeleteCall(id)
}

func (a *App) DeleteRollcall(id int32) error {
	return a.DeleteCall(id)
}

func (a *App) WebsitePDF(callID int32, period string, filePath string) error {
	coursePeriod, err := types.ParseCoursePeriod(period)
	if err != nil {
		return errors.New("Período inválido.")
	}

	cmd := commands.CreateWebsitePDFCommand{
		CallID:   callID,
		Period:   coursePeriod,
		FilePath: filePath,
	}

	if err := cmd.Execute(a.sisu.database); err != nil {
		return translateError(err)
	}

	return nil
}

func (a *App) EnrollmentPDF(callID int32, period string, filePath string) error {
	coursePeriod, err := types.ParseCoursePeriod(period)
	if err != nil {
		return errors.New("Período inválido.")
	}

	cmd := commands.CreateEnrollmentPDFCommand{
		CallID:   callID,
		Period:   coursePeriod,
		FilePath: filePath,
	}

	if err := cmd.Execute(a.sisu.database); err != nil {
		return translateError(err)
	}

	return nil
}

func (a *App) EmailPDF(callID int32, period string, filePath string) error {
	coursePeriod, err := types.ParseCoursePeriod(period)
	if err != nil {
		return errors.New("Período inválido.")
	}

	cmd := commands.CreateEmailPDFCommand{
		CallID:   callID,
		Period:   coursePeriod,
		FilePath: filePath,
	}

	if err := cmd.Execute(a.sisu.database); err != nil {
		return translateError(err)
	}

	return nil
}

func (a *App) TeacherPDF(period string, filePath string) error {
	coursePeriod, err := types.ParseCoursePeriod(period)
	if err != nil {
		return errors.New("Período inválido.")
	}

	cmd := commands.CreateTeacherPDFCommand{
		Period:   coursePeriod,
		FilePath: filePath,
	}

	if err := cmd.Execute(a.sisu.database); err != nil {
		return translateError(err)
	}

	return nil
}

func (a *App) ExportCSV(filePath string) error {
	cmd := commands.ExportCSVCommand{
		FilePath: filePath,
	}

	if err := cmd.Execute(a.sisu.database); err != nil {
		return translateError(err)
	}

	return nil
}

func (a *App) Backup(filePath string) error {
	cmd := commands.BackupCommand{
		FilePath: filePath,
	}

	if err := cmd.Execute(a.sisu.database); err != nil {
		return translateError(err)
	}

	return nil
}

func (a *App) Restore(filePath string) error {
	cmd := commands.RestoreCommand{
		FilePath: filePath,
	}

	if err := cmd.Execute(a.sisu.database); err != nil {
		return translateError(err)
	}

	return nil
}

func (a *App) Destroy() error {
	var options runtime.MessageDialogOptions
	options.Type = runtime.QuestionDialog
	options.Title = "Apagar Tudo"
	options.Message = "Essa ação apagará todas as informações contidas aqui e todas as alterações realizadas, antes de confirmar garanta que um backup foi realizado. Deseja continuar?" //nolint: lll

	result, err := runtime.MessageDialog(a.ctx, options)
	if err != nil {
		return err
	}

	matched, err := regexp.MatchString(`(Ok|Yes)`, result)
	if err != nil || !matched {
		return nil
	}

	cmd := commands.DestroyCommand{}

	if err := cmd.Execute(a.sisu.database); err != nil {
		return translateError(err)
	}

	runtime.Quit(a.ctx)

	return nil
}

/*

func (a *App) FetchPeriods() ([]types.CoursePeriod, error) {
	periods, err := a.sisu.FetchPeriods()
	if err != nil {
		return nil, err
	}

	return periods, nil
}
*/
