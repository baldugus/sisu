package types

import "encoding/json"

// ENUM(approved, waitlisted, absent, enrolled, declined_promotion)
type RegistrationStatus int

func (rs RegistrationStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(rs.String())
}

type Registration struct {
	ID                   int32              `csv:"-"`
	EnrollmentID         string             `csv:"INSCRICAO_ENEM"`
	Option               int32              `csv:"OPCAO"`
	LanguagesScore       *Score             `csv:"NOTA_LINGUAGENS"`
	HumanitiesScore      *Score             `csv:"NOTA_HUMANAS"`
	NaturalSciencesScore *Score             `csv:"NOTA_NATUREZA"`
	MathematicsScore     *Score             `csv:"NOTA_MATEMATICA"`
	EssayScore           *Score             `csv:"NOTA_REDACAO"`
	CompositeScore       *Score             `csv:"NOTA_FINAL"`
	Ranking              int32              `csv:"CLASSIFICACAO"`
	Status               RegistrationStatus `csv:"SITUACAO"`
	Candidate            *Candidate
	SemesterID           *int32
}

type Registrations []*Registration

func (v Registrations) Len() int           { return len(v) }
func (v Registrations) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v Registrations) Less(i, j int) bool { return v[i].Ranking < v[j].Ranking }
