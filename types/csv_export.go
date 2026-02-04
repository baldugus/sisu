package types

// EnrolledRegistrationCSV is a flat structure for CSV export of enrolled registrations.
type EnrolledRegistrationCSV struct {
	// Candidate fields
	CPF          string `csv:"CPF"`
	Name         string `csv:"NOME"`
	SocialName   string `csv:"NOME_SOCIAL"`
	BirthDate    string `csv:"DATA_NASCIMENTO"`
	Sex          string `csv:"SEXO"`
	MotherName   string `csv:"NOME_MAE"`
	AddressLine  string `csv:"ENDERECO"`
	AddressLine2 string `csv:"COMPLEMENTO"`
	HouseNumber  string `csv:"NUMERO"`
	Neighborhood string `csv:"BAIRRO"`
	Municipality string `csv:"MUNICIPIO"`
	State        string `csv:"UF"`
	CEP          string `csv:"CEP"`
	Email        string `csv:"EMAIL"`
	Phone1       string `csv:"TELEFONE1"`
	Phone2       string `csv:"TELEFONE2"`

	// Registration fields
	EnrollmentID         string `csv:"INSCRICAO_ENEM"`
	Option               int32  `csv:"OPCAO"`
	LanguagesScore       string `csv:"NOTA_LINGUAGENS"`
	HumanitiesScore      string `csv:"NOTA_HUMANAS"`
	NaturalSciencesScore string `csv:"NOTA_NATUREZA"`
	MathematicsScore     string `csv:"NOTA_MATEMATICA"`
	EssayScore           string `csv:"NOTA_REDACAO"`
	CompositeScore       string `csv:"NOTA_FINAL"`
	Ranking              int32  `csv:"CLASSIFICACAO"`
	Status               string `csv:"SITUACAO"`

	// Course fields
	Period string `csv:"TURNO"`
	Quota  string `csv:"MODALIDADE"`
}
