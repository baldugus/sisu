package types

type Candidate struct {
	ID           int32  `csv:"-"`
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
}
