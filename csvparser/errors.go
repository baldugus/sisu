package csvparser

import "fmt"

type ErrFileNotFound struct {
	Path string
	Err  error
}

func (e ErrFileNotFound) Error() string {
	return fmt.Sprintf("csv file not found: %s", e.Path)
}

func (e ErrFileNotFound) Unwrap() error {
	return e.Err
}

type ErrPermissionDenied struct {
	Path string
	Err  error
}

func (e ErrPermissionDenied) Error() string {
	return fmt.Sprintf("permission denied to open csv file: %s", e.Path)
}

func (e ErrPermissionDenied) Unwrap() error {
	return e.Err
}

type ErrFileOpen struct {
	Path string
	Err  error
}

func (e ErrFileOpen) Error() string {
	return fmt.Sprintf("opening csv file (%s): %v", e.Path, e.Err)
}

func (e ErrFileOpen) Unwrap() error {
	return e.Err
}

type ErrFileRead struct {
	Err error
}

func (e ErrFileRead) Error() string {
	return fmt.Sprintf("reading csv file: %v", e.Err)
}

func (e ErrFileRead) Unwrap() error {
	return e.Err
}

type ErrFileParse struct {
	Err error
}

func (e ErrFileParse) Error() string {
	return fmt.Sprintf("parsing csv file: %v", e.Err)
}

func (e ErrFileParse) Unwrap() error {
	return e.Err
}

type ErrFileEmpty struct {
	Err error
}

func (e ErrFileEmpty) Error() string {
	return fmt.Sprintf("csv file is empty: %v", e.Err)
}

func (e ErrFileEmpty) Unwrap() error {
	return e.Err
}

// TODO: temporary type before validation is done in the domain type
type ErrFieldValidation struct {
	Field string
	Err   error
}

func (e ErrFieldValidation) Error() string {
	return fmt.Sprintf("validating csv field %s: %v", e.Field, e.Err)
}

func (e ErrFieldValidation) Unwrap() error {
	return e.Err
}

type ErrLineMapping struct {
	Line int
	Err  error
}

func (e ErrLineMapping) Error() string {
	return fmt.Sprintf("mapping csv line %d: %v", e.Line, e.Err)
}

func (e ErrLineMapping) Unwrap() error {
	return e.Err
}

type ErrInvalidPeriod struct {
	Period string
}

func (e ErrInvalidPeriod) Error() string {
	return fmt.Sprintf("not a valid period: %s", e.Period)
}

type ErrNoRegistrationsFound struct{}

func (e ErrNoRegistrationsFound) Error() string {
	return "no candidates found in file"
}

type ErrOddSeatsCount struct {
	Count int64
}

func (e ErrOddSeatsCount) Error() string {
	return fmt.Sprintf("total seats must be an even number, got %d", e.Count)
}
