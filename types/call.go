package types

import "encoding/json"

// ENUM(calling, done)
type CallStatus int

func (cs CallStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(cs.String())
}

type Call struct {
	ID     int32      `csv:"-"`
	Status CallStatus `csv:"-"`
	Number int32      `csv:"CHAMADA"`
}
