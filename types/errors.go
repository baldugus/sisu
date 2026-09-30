package types

import "fmt"

type ErrOddSeatsCount struct {
	Count int32
}

func (e ErrOddSeatsCount) Error() string {
	return fmt.Sprintf("total seats must be an even number, got %d", e.Count)
}
