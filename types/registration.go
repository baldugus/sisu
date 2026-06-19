package types

// ENUM(approved, waitlisted, absent, enrolled, declined_promotion)
type RegistrationStatus string

type Registration struct {
	ID                   int32              `csv:"-"`
	EnrollmentID         string             `csv:"INSCRICAO_ENEM"`
	Option               int32              `csv:"OPCAO"`
	LanguagesScore       *Score             `csv:"NOTA_LINGUAGENS" ts_type:"string"`
	HumanitiesScore      *Score             `csv:"NOTA_HUMANAS" ts_type:"string"`
	NaturalSciencesScore *Score             `csv:"NOTA_NATUREZA" ts_type:"string"`
	MathematicsScore     *Score             `csv:"NOTA_MATEMATICA" ts_type:"string"`
	EssayScore           *Score             `csv:"NOTA_REDACAO" ts_type:"string"`
	CompositeScore       *Score             `csv:"NOTA_FINAL" ts_type:"string"`
	Ranking              int32              `csv:"CLASSIFICACAO"`
	Status               RegistrationStatus `csv:"SITUACAO"`
	Candidate            *Candidate
	SemesterID           *int32
}

type Registrations []*Registration

func (v Registrations) Len() int           { return len(v) }
func (v Registrations) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v Registrations) Less(i, j int) bool { return v[i].Ranking < v[j].Ranking }
