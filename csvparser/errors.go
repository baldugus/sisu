package csvparser

import "fmt"

type FileNotFoundError struct {
	Path string
	Err  error
}

func (e FileNotFoundError) Error() string {
	return fmt.Sprintf("csv file not found: %s", e.Path)
}

func (e FileNotFoundError) Unwrap() error {
	return e.Err
}

type PermissionDeniedError struct {
	Path string
	Err  error
}

func (e PermissionDeniedError) Error() string {
	return fmt.Sprintf("permission denied to open csv file: %s", e.Path)
}

func (e PermissionDeniedError) Unwrap() error {
	return e.Err
}

type FileError struct {
	Path string
	Err  error
}

func (e FileError) Error() string {
	return fmt.Sprintf("opening csv file (%s): %v", e.Path, e.Err)
}

func (e FileError) Unwrap() error {
	return e.Err
}

type ReadError struct {
	Err error
}

func (e ReadError) Error() string {
	return fmt.Sprintf("reading csv file: %v", e.Err)
}

func (e ReadError) Unwrap() error {
	return e.Err
}

type ParseError struct {
	Err error
}

func (e ParseError) Error() string {
	return fmt.Sprintf("parsing csv file: %v", e.Err)
}

func (e ParseError) Unwrap() error {
	return e.Err
}

type EmptyError struct {
	Err error
}

func (e EmptyError) Error() string {
	return fmt.Sprintf("csv file is empty: %v", e.Err)
}

func (e EmptyError) Unwrap() error {
	return e.Err
}

// TODO: temporary type before validation is done in the domain type
type ValidationError struct {
	Field string
	Err   error
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validating csv field %s: %v", e.Field, e.Err)
}

func (e ValidationError) Unwrap() error {
	return e.Err
}

type MapperError struct {
	Line int
	Err  error
}

func (e MapperError) Error() string {
	return fmt.Sprintf("mapping csv line %d: %v", e.Line, e.Err)
}

func (e MapperError) Unwrap() error {
	return e.Err
}
