// nolint: godox, nolintlint
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

type Response struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Data   any    `json:"data"`
}

func AppErrorToResponse(err error) Response {
	switch {
	case errors.As(err, &commands.ErrApprovedSelectionAlreadyExists{}):
		return Response{500, "Não é possível importar a lista de aprovados quando já existe uma seleção.", ""}
	case errors.As(err, &commands.ErrWaitlistSelectionRequiresApproved{}):
		return Response{500, "Não é possível importar a lista de espera sem a lista de aprovados.", ""}
	case errors.As(err, &csvparser.ErrFileNotFound{}):
		return Response{500, "Arquivo selecionado não existe.", ""}
	case errors.As(err, &csvparser.ErrPermissionDenied{}):
		return Response{500, "Permissão negada para acessar o arquivo.", ""}
	case errors.As(err, &csvparser.ErrFileOpen{}):
		return Response{500, "Erro ao abrir o arquivo CSV. Contate o desenvolvedor.", ""}
	case errors.As(err, &csvparser.ErrFileRead{}):
		return Response{500, "Erro ao ler o arquivo CSV. Contate o desenvolvedor.", ""}
	case errors.As(err, &csvparser.ErrFileEmpty{}):
		return Response{500, "Arquivo CSV vazio.", ""}
	case errors.As(err, &csvparser.ErrFileParse{}):
		return Response{500, "Erro ao analisar o arquivo CSV. Contate o desenvolvedor.", ""}
	case errors.As(err, &csvparser.ErrOddSeatsCount{}):
		return Response{500, "O número total de vagas deve ser par para divisão entre semestres.", ""}
	case errors.As(err, &csvparser.ErrLineMapping{}):
		return Response{500, "Erro de validação em um campo do arquivo CSV. Contate o desenvolvedor.", ""}
	case errors.As(err, &commands.ErrSelectionNotFound{}):
		return Response{500, "Seleção não encontrada.", ""}
	case errors.As(err, &commands.ErrCannotDeleteApprovedWithWaitlist{}):
		return Response{500, "Não é possível excluir a lista de aprovados enquanto a lista de espera existir.", ""}
	case errors.As(err, &commands.ErrSelectionHasModifiedRegistrations{}):
		return Response{500, "Não é possível excluir a seleção com inscrições matriculadas ou ausentes.", ""}
	case errors.As(err, &commands.ErrCallHasPendingRegistrations{}):
		return Response{500, "Não é possível fechar a chamada com inscrições pendentes.", ""}
	case errors.As(err, &commands.ErrRegistrationNotFound{}):
		return Response{500, "Inscrição não encontrada.", ""}
	case errors.As(err, &commands.ErrRegistrationNotInCall{}):
		return Response{500, "Inscrição não está em uma chamada.", ""}
	case errors.As(err, &commands.ErrCallNotOpen{}):
		return Response{500, "A chamada não está aberta.", ""}
	case errors.As(err, &commands.ErrInvalidStatusTransition{}):
		return Response{500, "Transição de status inválida.", ""}
	case errors.As(err, &commands.ErrNoSeatsAvailable{}):
		return Response{500, "Não há vagas disponíveis no curso.", ""}
	case errors.As(err, &commands.ErrOpenCallExists{}):
		return Response{500, "Não é possível criar nova chamada enquanto outra está aberta.", ""}
	case errors.As(err, &commands.ErrNoWaitlistedRegistrations{}):
		return Response{500, "Não há inscrições na lista de espera.", ""}
	case errors.As(err, &commands.ErrAllCoursesFull{}):
		return Response{500, "Todas as vagas de todos os cursos estão ocupadas.", ""}
	case errors.As(err, &commands.ErrCannotReopenCallWithLaterCalls{}):
		return Response{500, "Não é possível reabrir a chamada enquanto existem chamadas posteriores.", ""}
	case errors.As(err, &commands.ErrCannotDeleteWithMultipleCalls{}):
		return Response{500, "Não é possível excluir a seleção quando existem múltiplas chamadas.", ""}
	case errors.As(err, &commands.ErrCannotDeleteWithClosedFirstCall{}):
		return Response{500, "Não é possível excluir a seleção quando a primeira chamada está fechada.", ""}
	case errors.As(err, &commands.ErrCannotDeleteCallWithLaterCalls{}):
		return Response{500, "Não é possível excluir a chamada enquanto existem chamadas posteriores.", ""}
	case errors.As(err, &commands.ErrCannotDeleteClosedCall{}):
		return Response{500, "Não é possível excluir uma chamada fechada.", ""}
	default:
		zap.L().Error("unknown error", zap.Error(err))
		return Response{500, "Erro desconhecido. Contate o desenvolvedor.", ""}
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

func (a *App) LoadApprovedSelection(year int32, filePath string) Response {
	cmd := commands.LoadSelectionCommand{
		Year:     year,
		FilePath: filePath,
		Kind:     types.SelectionKindApproved,
	}

	if err := cmd.Execute(a.sisu.database); err != nil {
		return AppErrorToResponse(err)
	}

	return Response{200, "OK", ""}
}

func (a *App) LoadWaitlistSelection(year int32, filePath string) Response {
	cmd := commands.LoadSelectionCommand{
		Year:     year,
		FilePath: filePath,
		Kind:     types.SelectionKindWaitlist,
	}

	if err := cmd.Execute(a.sisu.database); err != nil {
		return AppErrorToResponse(err)
	}

	return Response{200, "OK", ""}
}

func (a *App) LoadInterestedSelection(year int32, filePath string) Response {
	return a.LoadWaitlistSelection(year, filePath)
}

func (a *App) FetchApprovedSelection() Response {
	cmd := commands.FetchSelectionCommand{
		Kind: types.SelectionKindApproved,
	}

	selection, err := cmd.Execute(a.sisu.database)
	if err != nil {
		if errors.As(err, &database.ErrNotFound{}) {
			return Response{200, "OK", nil}
		}
		return AppErrorToResponse(err)
	}

	return Response{200, "OK", selection}
}

func (a *App) FetchInterestedSelection() Response {
	cmd := commands.FetchSelectionCommand{
		Kind: types.SelectionKindWaitlist,
	}

	selection, err := cmd.Execute(a.sisu.database)
	if err != nil {
		if errors.As(err, &database.ErrNotFound{}) {
			return Response{200, "OK", nil}
		}
		return AppErrorToResponse(err)
	}

	return Response{200, "OK", selection}
}

func (a *App) FetchRollCalls() Response {
	cmd := commands.FetchCallsCommand{}

	calls, err := cmd.Execute(a.sisu.database)
	if err != nil {
		return AppErrorToResponse(err)
	}

	return Response{200, "OK", calls}
}

func (a *App) FetchSemesters() Response {
	cmd := commands.FetchSemestersCommand{}

	semesters, err := cmd.Execute(a.sisu.database)
	if err != nil {
		return AppErrorToResponse(err)
	}

	return Response{200, "OK", semesters}
}

func (a *App) DeleteApprovedSelection() Response {
	cmd := commands.DeleteSelectionCommand{
		Kind: types.SelectionKindApproved,
	}

	if err := cmd.Execute(a.sisu.database); err != nil {
		return AppErrorToResponse(err)
	}

	return Response{200, "OK", ""}
}

func (a *App) DeleteInterestedSelection() Response {
	cmd := commands.DeleteSelectionCommand{
		Kind: types.SelectionKindWaitlist,
	}

	if err := cmd.Execute(a.sisu.database); err != nil {
		return AppErrorToResponse(err)
	}

	return Response{200, "OK", ""}
}

func (a *App) FetchRegistrations() Response {
	cmd := commands.FetchRegistrationsCommand{}

	registrations, err := cmd.Execute(a.sisu.database)
	if err != nil {
		return AppErrorToResponse(err)
	}

	return Response{200, "OK", registrations}
}

func (a *App) FetchRegistrationsBySelectionID(selectionID int32) Response {
	cmd := commands.FetchRegistrationsCommand{
		SelectionID: &selectionID,
	}

	registrations, err := cmd.Execute(a.sisu.database)
	if err != nil {
		return AppErrorToResponse(err)
	}

	return Response{200, "OK", registrations}
}

func (a *App) FetchRegistrationsByCourseID(courseID int32) Response {
	cmd := commands.FetchRegistrationsCommand{
		CourseID: &courseID,
	}

	registrations, err := cmd.Execute(a.sisu.database)
	if err != nil {
		return AppErrorToResponse(err)
	}

	return Response{200, "OK", registrations}
}

func (a *App) FetchRegistrationsByCallID(callID int32) Response {
	cmd := commands.FetchRegistrationsCommand{
		CallID: &callID,
	}

	registrations, err := cmd.Execute(a.sisu.database)
	if err != nil {
		return AppErrorToResponse(err)
	}

	return Response{200, "OK", registrations}
}

func (a *App) FetchRegistration(registrationID int32) Response {
	cmd := commands.FetchRegistrationCommand{
		RegistrationID: registrationID,
	}

	registration, err := cmd.Execute(a.sisu.database)
	if err != nil {
		return AppErrorToResponse(err)
	}

	return Response{200, "OK", registration}
}

func (a *App) CloseRollCall(id int32) Response {
	cmd := commands.CloseCallCommand{
		ID: id,
	}

	if err := cmd.Execute(a.sisu.database); err != nil {
		return AppErrorToResponse(err)
	}

	return Response{200, "OK", ""}
}

func (a *App) EnrollRegistration(id int32) Response {
	cmd := commands.UpdateRegistrationStatusCommand{
		RegistrationID: id,
		NewStatus:      types.RegistrationStatusEnrolled,
	}

	if err := cmd.Execute(a.sisu.database); err != nil {
		return AppErrorToResponse(err)
	}

	return Response{200, "OK", ""}
}

func (a *App) AbsentRegistration(id int32) Response {
	cmd := commands.UpdateRegistrationStatusCommand{
		RegistrationID: id,
		NewStatus:      types.RegistrationStatusAbsent,
	}

	if err := cmd.Execute(a.sisu.database); err != nil {
		return AppErrorToResponse(err)
	}

	return Response{200, "OK", ""}
}

func (a *App) ClearRegistrationStatus(id int32) Response {
	cmd := commands.UpdateRegistrationStatusCommand{
		RegistrationID: id,
		NewStatus:      types.RegistrationStatusApproved,
	}

	if err := cmd.Execute(a.sisu.database); err != nil {
		return AppErrorToResponse(err)
	}

	return Response{200, "OK", ""}
}

func (a *App) EnrollApplication(id int32) Response {
	return a.EnrollRegistration(id)
}

func (a *App) AbsentApplication(id int32) Response {
	return a.AbsentRegistration(id)
}

func (a *App) ClearApplicationStatus(id int32) Response {
	return a.ClearRegistrationStatus(id)
}

func (a *App) FetchApplicationsByRollCall(callID int32) Response {
	return a.FetchRegistrationsByCallID(callID)
}

func (a *App) CreateRollCall(semesterID int32) Response {
	cmd := commands.CreateCallCommand{SemesterID: semesterID}

	if err := cmd.Execute(a.sisu.database); err != nil {
		return AppErrorToResponse(err)
	}

	return Response{200, "OK", ""}
}

func (a *App) OpenCall(id int32) Response {
	cmd := commands.OpenCallCommand{
		ID: id,
	}

	if err := cmd.Execute(a.sisu.database); err != nil {
		return AppErrorToResponse(err)
	}

	return Response{200, "OK", ""}
}

func (a *App) DeleteCall(id int32) Response {
	cmd := commands.DeleteCallCommand{
		ID: id,
	}

	if err := cmd.Execute(a.sisu.database); err != nil {
		return AppErrorToResponse(err)
	}

	return Response{200, "OK", ""}
}

func (a *App) OpenRollCall(id int32) Response {
	return a.OpenCall(id)
}

func (a *App) DeleteRollCall(id int32) Response {
	return a.DeleteCall(id)
}

func (a *App) DeleteRollcall(id int32) Response {
	return a.DeleteCall(id)
}

func (a *App) WebsitePDF(callID int32, period string, filePath string) Response {
	coursePeriod, err := types.ParseCoursePeriod(period)
	if err != nil {
		return Response{500, "Período inválido.", ""}
	}

	cmd := commands.CreateWebsitePDFCommand{
		CallID:   callID,
		Period:   coursePeriod,
		FilePath: filePath,
	}

	if err := cmd.Execute(a.sisu.database); err != nil {
		return AppErrorToResponse(err)
	}

	return Response{200, "OK", ""}
}

func (a *App) EnrollmentPDF(callID int32, period string, filePath string) Response {
	coursePeriod, err := types.ParseCoursePeriod(period)
	if err != nil {
		return Response{500, "Período inválido.", ""}
	}

	cmd := commands.CreateEnrollmentPDFCommand{
		CallID:   callID,
		Period:   coursePeriod,
		FilePath: filePath,
	}

	if err := cmd.Execute(a.sisu.database); err != nil {
		return AppErrorToResponse(err)
	}

	return Response{200, "OK", ""}
}

func (a *App) EmailPDF(callID int32, period string, filePath string) Response {
	coursePeriod, err := types.ParseCoursePeriod(period)
	if err != nil {
		return Response{500, "Período inválido.", ""}
	}

	cmd := commands.CreateEmailPDFCommand{
		CallID:   callID,
		Period:   coursePeriod,
		FilePath: filePath,
	}

	if err := cmd.Execute(a.sisu.database); err != nil {
		return AppErrorToResponse(err)
	}

	return Response{200, "OK", ""}
}

func (a *App) TeacherPDF(period string, filePath string) Response {
	coursePeriod, err := types.ParseCoursePeriod(period)
	if err != nil {
		return Response{500, "Período inválido.", ""}
	}

	cmd := commands.CreateTeacherPDFCommand{
		Period:   coursePeriod,
		FilePath: filePath,
	}

	if err := cmd.Execute(a.sisu.database); err != nil {
		return AppErrorToResponse(err)
	}

	return Response{200, "OK", ""}
}

func (a *App) ExportCSV(filePath string) Response {
	cmd := commands.ExportCSVCommand{
		FilePath: filePath,
	}

	if err := cmd.Execute(a.sisu.database); err != nil {
		return AppErrorToResponse(err)
	}

	return Response{200, "OK", ""}
}

func (a *App) Backup(filePath string) Response {
	cmd := commands.BackupCommand{
		FilePath: filePath,
	}

	if err := cmd.Execute(a.sisu.database); err != nil {
		return AppErrorToResponse(err)
	}

	return Response{200, "OK", ""}
}

func (a *App) Restore(filePath string) Response {
	cmd := commands.RestoreCommand{
		FilePath: filePath,
	}

	if err := cmd.Execute(a.sisu.database); err != nil {
		return AppErrorToResponse(err)
	}

	return Response{200, "OK", ""}
}

func (a *App) Destroy() Response {
	var options runtime.MessageDialogOptions
	options.Type = runtime.QuestionDialog
	options.Title = "Apagar Tudo"
	options.Message = "Essa ação apagará todas as informações contidas aqui e todas as alterações realizadas, antes de confirmar garanta que um backup foi realizado. Deseja continuar?" //nolint: lll

	result, err := runtime.MessageDialog(a.ctx, options)
	if err != nil {
		return Response{500, err.Error(), ""}
	}

	matched, err := regexp.MatchString(`(Ok|Yes)`, result)
	if err != nil || !matched {
		return Response{200, "OK", ""}
	}

	cmd := commands.DestroyCommand{}

	if err := cmd.Execute(a.sisu.database); err != nil {
		return AppErrorToResponse(err)
	}

	runtime.Quit(a.ctx)

	return Response{200, "OK", ""}
}

/*

func (a *App) FetchPeriods() Response {
	periods, err := a.sisu.FetchPeriods()
	if err != nil {
		return Response{500, err.Error(), ""}
	}

	return Response{200, "OK", periods}
}
*/
