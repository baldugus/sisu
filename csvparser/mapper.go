package csvparser

import (
	"strconv"
	"strings"

	"github.com/baldugus/sisu/types"
)

type ParsedCsv struct {
	applicants []*csvApplicant
	name       string
}

func (pc *ParsedCsv) ToApplications(status string) ([]*types.Application, error) {
	applications := make([]*types.Application, len(pc.applicants))

	for i, applicant := range pc.applicants {
		application, err := applicant.Parse(status)
		if err != nil {
			// Line: i + 1 assumes the first line in a csv file is headers
			return nil, &MapperError{Line: i + 1, Err: err}
		}

		applications[i] = application
	}

	return applications, nil
}

func (pc *ParsedCsv) ToSelection(kind types.SelectionKind, date string) types.Selection {
	var selection types.Selection
	selection.Name = pc.name
	selection.Kind = kind

	if len(pc.applicants) > 0 {
		selection.Date = date
		selection.Institution = pc.applicants[0].Institution
		selection.Course = pc.applicants[0].Course
	}

	return selection
}

type csvApplicant struct {
	Stage                string `csv:"NU_ETAPA"`
	SchedulePeriod       string `csv:"DS_TURNO"`
	Seats                string `csv:"QT_VAGAS_CONCORRENCIA"`
	EnrollmentID         string `csv:"CO_INSCRICAO_ENEM"`
	Name                 string `csv:"NO_INSCRITO"`
	CPF                  string `csv:"NU_CPF_INSCRITO"`
	Date                 string `csv:"DT_CURSO_INSCRICAO"`
	SocialName           string `csv:"NO_SOCIAL"`
	BirthDate            string `csv:"DT_NASCIMENTO"`
	Sex                  string `csv:"TP_SEXO"`
	MotherName           string `csv:"NO_MAE"`
	AddressLine          string `csv:"DS_LOGRADOURO"`
	HouseNumber          string `csv:"NU_ENDERECO"`
	AddressLine2         string `csv:"DS_COMPLEMENTO"`
	State                string `csv:"SG_UF_INSCRITO"`
	Municipality         string `csv:"NO_MUNICIPIO"`
	Neighborhood         string `csv:"NO_BAIRRO"`
	CEP                  string `csv:"NU_CEP"`
	Phone1               string `csv:"NU_FONE1"`
	Phone2               string `csv:"NU_FONE2"`
	Email                string `csv:"DS_EMAIL"`
	LanguagesScore       string `csv:"NU_NOTA_L"`
	HumanitiesScore      string `csv:"NU_NOTA_CH"`
	NaturalSciencesScore string `csv:"NU_NOTA_CN"`
	MathematicsScore     string `csv:"NU_NOTA_M"`
	EssayScore           string `csv:"NU_NOTA_R"`
	Option               string `csv:"ST_OPCAO"`
	Quota                string `csv:"NO_MODALIDADE_CONCORRENCIA"`
	CompositeScore       string `csv:"NU_NOTA_CANDIDATO"`
	MinimumScore         string `csv:"NU_NOTACORTE_CONCORRIDA"`
	Ranking              string `csv:"NU_CLASSIFICACAO"`
	Institution          string `csv:"SG_IES"`
	Course               string `csv:"NO_CURSO"`
}

// TODO: move all of these validations to a parser in the types package
// when the database is decoupled from domain types
func (a *csvApplicant) Parse(status string) (*types.Application, error) {
	minimumScore, err := parseLocalizedFloat(a.MinimumScore)
	if err != nil {
		return nil, &ValidationError{Field: "MinimumScore", Err: err}
	}

	seats, err := strconv.Atoi(a.Seats)
	if err != nil {
		return nil, &ValidationError{Field: "Seats", Err: err}
	}

	option, err := strconv.Atoi(a.Option)
	if err != nil {
		return nil, &ValidationError{Field: "Option", Err: err}
	}

	languagesScore, err := parseLocalizedFloat(a.LanguagesScore)
	if err != nil {
		return nil, &ValidationError{Field: "LanguagesScore", Err: err}
	}

	humanitiesScore, err := parseLocalizedFloat(a.HumanitiesScore)
	if err != nil {
		return nil, &ValidationError{Field: "HumanitiesScore", Err: err}
	}

	naturalSciencesScore, err := parseLocalizedFloat(a.NaturalSciencesScore)
	if err != nil {
		return nil, &ValidationError{Field: "NaturalSciencesScore", Err: err}
	}

	mathematicsScore, err := parseLocalizedFloat(a.MathematicsScore)
	if err != nil {
		return nil, &ValidationError{Field: "MathematicsScore", Err: err}
	}

	essayScore, err := parseLocalizedFloat(a.EssayScore)
	if err != nil {
		return nil, &ValidationError{Field: "EssayScore", Err: err}
	}

	compositeScore, err := parseLocalizedFloat(a.CompositeScore)
	if err != nil {
		return nil, &ValidationError{Field: "CompositeScore", Err: err}
	}

	ranking, err := strconv.Atoi(a.Ranking)
	if err != nil {
		return nil, &ValidationError{Field: "Ranking", Err: err}
	}

	return &types.Application{
		ID:                   0,
		Status:               status,
		EnrollmentID:         a.EnrollmentID,
		Option:               option,
		LanguagesScore:       languagesScore,
		HumanitiesScore:      humanitiesScore,
		NaturalSciencesScore: naturalSciencesScore,
		MathematicsScore:     mathematicsScore,
		EssayScore:           essayScore,
		CompositeScore:       compositeScore,
		Ranking:              ranking,
		Applicant: types.Applicant{
			CPF:          a.CPF,
			Name:         a.Name,
			SocialName:   a.SocialName,
			BirthDate:    a.BirthDate,
			Sex:          a.Sex,
			MotherName:   a.MotherName,
			AddressLine:  a.AddressLine,
			AddressLine2: a.AddressLine2,
			HouseNumber:  a.HouseNumber,
			Neighborhood: a.Neighborhood,
			Municipality: a.Municipality,
			State:        a.State,
			CEP:          a.CEP,
			Email:        a.Email,
			Phone1:       a.Phone1,
			Phone2:       a.Phone2,
		},
		Class: types.Class{
			ID: 0,
			Period: types.Period{
				Name: a.SchedulePeriod,
			},
			Quota: types.Quota{
				Name: a.Quota,
			},
			Seats:        seats,
			MinimumScore: minimumScore,
		},
	}, nil

}

func parseLocalizedFloat(s string) (float64, error) {
	score, err := strconv.ParseFloat(strings.ReplaceAll(s, ",", "."), 64)
	if err != nil {
		return 0, err
	}

	return score, nil
}
