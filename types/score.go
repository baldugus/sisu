package types

import (
	"fmt"
)

type Score struct {
	Value int32
}

func NewScoreFromFloat(value float32) Score {
	// Float is stupid, 655.17*100 equals 65516
	// 0.5 is needed for rounding
	score := value*100 + 0.5

	return Score{Value: int32(score)}
}

func (s Score) String() string {
	value := s.Value / 100
	decimals := s.Value % 100

	if decimals < 0 {
		decimals = -decimals
	}

	return fmt.Sprintf("%d,%02d", value, decimals)
}
