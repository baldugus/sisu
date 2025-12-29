// nolint: godox, nolintlint
// TODO: make some helper functions to avoid the repetition.
package main

import (
	"context"
	"fmt"
	"regexp"

	"changeme/types"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Response struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Data   any    `json:"data"`
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
}

func (a *App) LoadApprovedSelection() Response {
	var filter runtime.FileFilter
	filter.DisplayName = "CSV File"
	filter.Pattern = "*.csv"

	var options runtime.OpenDialogOptions
	options.Filters = []runtime.FileFilter{filter}

	file, err := runtime.OpenFileDialog(a.ctx, options)
	if err != nil {
		return Response{500, err.Error(), ""}
	}

	if err := a.sisu.LoadSelection(file, types.ApprovedSelection); err != nil {
		return Response{500, err.Error(), ""}
	}

	return Response{200, "OK", ""}
}

func (a *App) LoadInterestedSelection() Response {
	var filter runtime.FileFilter
	filter.DisplayName = "CSV File"
	filter.Pattern = "*.csv"

	var options runtime.OpenDialogOptions
	options.Filters = []runtime.FileFilter{filter}

	file, err := runtime.OpenFileDialog(a.ctx, options)
	if err != nil {
		return Response{500, err.Error(), ""}
	}

	if err := a.sisu.LoadSelection(file, types.InterestedSelection); err != nil {
		return Response{500, err.Error(), ""}
	}

	return Response{200, "OK", ""}
}

func (a *App) CloseRollCall(id int64) Response {
	if err := a.sisu.CloseRollCall(id); err != nil {
		return Response{500, err.Error(), ""}
	}

	return Response{200, "OK", ""}
}

func (a *App) OpenRollCall(id int64) Response {
	if err := a.sisu.OpenRollCall(id); err != nil {
		return Response{500, err.Error(), ""}
	}

	return Response{200, "OK", ""}
}

func (a *App) FetchApprovedSelection() Response {
	selection, err := a.sisu.FetchApprovedSelection()
	if err != nil {
		return Response{500, err.Error(), ""}
	}

	return Response{200, "OK", selection}
}

func (a *App) FetchInterestedSelection() Response {
	selection, err := a.sisu.FetchInterestedSelection()
	if err != nil {
		return Response{500, err.Error(), ""}
	}

	return Response{200, "OK", selection}
}

func (a *App) FetchRollCalls() Response {
	rollcalls, err := a.sisu.FetchRollcalls()
	if err != nil {
		return Response{500, err.Error(), ""}
	}

	return Response{200, "OK", rollcalls}
}

func (a *App) FetchPeriods() Response {
	periods, err := a.sisu.FetchPeriods()
	if err != nil {
		return Response{500, err.Error(), ""}
	}

	return Response{200, "OK", periods}
}

func (a *App) FetchApplicationsByRollCall(id int64) Response {
	applications, err := a.sisu.FetchApplicationsByRollCall(id)
	if err != nil {
		return Response{500, err.Error(), ""}
	}

	return Response{200, "OK", applications}
}

func (a *App) EnrollApplication(id int64) Response {
	if err := a.sisu.EnrollApplication(id); err != nil {
		return Response{500, err.Error(), ""}
	}

	return Response{200, "OK", ""}
}

func (a *App) ClearApplicationStatus(id int64) Response {
	if err := a.sisu.ClearApplicationStatus(id); err != nil {
		return Response{500, err.Error(), ""}
	}

	return Response{200, "OK", ""}
}

func (a *App) AbsentApplication(id int64) Response {
	if err := a.sisu.AbsentApplication(id); err != nil {
		return Response{500, err.Error(), ""}
	}

	return Response{200, "OK", ""}
}

func (a *App) CreateRollCall() Response {
	if err := a.sisu.CreateRollCall(); err != nil {
		return Response{500, err.Error(), ""}
	}

	return Response{200, "OK", ""}
}

func (a *App) WebsitePDF(rollcallID int64, periodID int64) Response {
	period, err := a.sisu.FetchPeriodName(periodID)
	if err != nil {
		return Response{500, err.Error(), ""}
	}

	rollcall, err := a.sisu.FetchRollcallNumber(rollcallID)
	if err != nil {
		return Response{500, err.Error(), ""}
	}

	var options runtime.SaveDialogOptions
	options.DefaultFilename = fmt.Sprintf("%s-site-chamada-%d.pdf", period, rollcall)

	file, err := runtime.SaveFileDialog(a.ctx, options)
	if err != nil {
		return Response{500, err.Error(), ""}
	}

	if err := a.sisu.WebsitePDF(rollcallID, periodID, file); err != nil {
		return Response{500, err.Error(), ""}
	}

	return Response{200, "OK", ""}
}

func (a *App) ExportCSV() Response {
	var options runtime.SaveDialogOptions
	options.DefaultFilename = fmt.Sprintf("alunos-aprovados.csv")

	file, err := runtime.SaveFileDialog(a.ctx, options)
	if err != nil {
		return Response{500, err.Error(), ""}
	}

	if err := a.sisu.ExportCSV(file); err != nil {
		return Response{500, err.Error(), ""}
	}

	return Response{200, "OK", ""}
}

func (a *App) Backup() Response {
	var options runtime.SaveDialogOptions
	options.DefaultFilename = fmt.Sprintf("backup.sisu")

	file, err := runtime.SaveFileDialog(a.ctx, options)
	if err != nil {
		return Response{500, err.Error(), ""}
	}

	if err := a.sisu.Backup(file); err != nil {
		return Response{500, err.Error(), ""}
	}

	return Response{200, "OK", ""}
}

func (a *App) Restore() Response {
	var filter runtime.FileFilter
	filter.DisplayName = "Restore Backup"
	filter.Pattern = "*.sisu"

	var options runtime.OpenDialogOptions
	options.Filters = []runtime.FileFilter{filter}

	file, err := runtime.OpenFileDialog(a.ctx, options)
	if err != nil {
		return Response{500, err.Error(), ""}
	}

	if err := a.sisu.Restore(file); err != nil {
		return Response{500, err.Error(), ""}
	}

	return Response{200, "OK", ""}
}

func (a *App) EnrollmentPDF(rollcallID int64, periodID int64) Response {
	period, err := a.sisu.FetchPeriodName(periodID)
	if err != nil {
		return Response{500, err.Error(), ""}
	}

	rollcall, err := a.sisu.FetchRollcallNumber(rollcallID)
	if err != nil {
		return Response{500, err.Error(), ""}
	}

	var options runtime.SaveDialogOptions
	options.DefaultFilename = fmt.Sprintf("%s-matricula-chamada-%d.pdf", period, rollcall)

	file, err := runtime.SaveFileDialog(a.ctx, options)
	if err != nil {
		return Response{500, err.Error(), ""}
	}

	if err := a.sisu.EnrollmentPDF(rollcallID, periodID, file); err != nil {
		return Response{500, err.Error(), ""}
	}

	return Response{200, "OK", ""}
}

func (a *App) EmailPDF(rollcallID int64, periodID int64) Response {
	period, err := a.sisu.FetchPeriodName(periodID)
	if err != nil {
		return Response{500, err.Error(), ""}
	}

	rollcall, err := a.sisu.FetchRollcallNumber(rollcallID)
	if err != nil {
		return Response{500, err.Error(), ""}
	}

	var options runtime.SaveDialogOptions
	options.DefaultFilename = fmt.Sprintf("%s-email-chamada-%d.pdf", period, rollcall)

	file, err := runtime.SaveFileDialog(a.ctx, options)
	if err != nil {
		return Response{500, err.Error(), ""}
	}

	if err := a.sisu.EmailPDF(rollcallID, periodID, file); err != nil {
		return Response{500, err.Error(), ""}
	}

	return Response{200, "OK", ""}
}

func (a *App) TeacherPDF(periodID int64) Response {
	period, err := a.sisu.FetchPeriodName(periodID)
	if err != nil {
		return Response{500, err.Error(), ""}
	}

	var options runtime.SaveDialogOptions
	options.DefaultFilename = fmt.Sprintf("presenca-%s.pdf", period)

	file, err := runtime.SaveFileDialog(a.ctx, options)
	if err != nil {
		return Response{500, err.Error(), ""}
	}

	if err := a.sisu.TeacherPDF(periodID, file); err != nil {
		return Response{500, err.Error(), ""}
	}

	return Response{200, "OK", ""}
}

func (a *App) DeleteApprovedSelection() Response {
	if err := a.sisu.DeleteSelection(types.ApprovedSelection); err != nil {
		return Response{500, err.Error(), ""}
	}

	return Response{200, "OK", ""}
}

func (a *App) DeleteInterestedSelection() Response {
	if err := a.sisu.DeleteSelection(types.InterestedSelection); err != nil {
		return Response{500, err.Error(), ""}
	}

	return Response{200, "OK", ""}
}

func (a *App) DeleteRollcall(id int64) Response {
	if err := a.sisu.DeleteRollcall(id); err != nil {
		return Response{500, err.Error(), ""}
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
		return Response{500, err.Error(), ""}
	}

	a.sisu.Destroy()

	runtime.Quit(a.ctx)

	return Response{200, "OK", ""}
}
