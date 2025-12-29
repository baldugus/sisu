package database

import "fmt"

type ErrQuery struct {
	Query string
	Err   error
}

func (e ErrQuery) Error() string {
	return fmt.Sprintf("query exec error %q: %s", e.Query, e.Err)
}

func (e ErrQuery) Unwrap() error {
	return e.Err
}

type ErrUnexpectedNumOfRowsAffected struct {
	Expected int
	Actual   int
}

func (e ErrUnexpectedNumOfRowsAffected) Error() string {
	return fmt.Sprintf("unexpected number of rows affected: expected %d, got %d", e.Expected, e.Actual)
}

type ErrNotFound struct {
	Err error
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("not found: %v", e.Err)
}

func (e ErrNotFound) Unwrap() error {
	return e.Err
}
