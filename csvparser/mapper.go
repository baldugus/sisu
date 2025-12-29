package csvparser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/baldugus/sisu/types"
)

type ParsedCsv struct {
	candidates []*csvCandidate
	name       string
}

type ParsedRegistration struct {
	Registration *types.Registration
	Course       *types.Course
	Call         *types.Call
}

type ParsedSelection struct {
	Selection     *types.Selection
	Registrations []*ParsedRegistration
}

func (pc *ParsedCsv) toRegistrationsDomain(status types.RegistrationStatus) ([]*ParsedRegistration, error) {
	registrations := make([]*ParsedRegistration, len(pc.candidates))

	for i, candidate := range pc.candidates {
		parsed, err := candidate.Parse(status)
		if err != nil {
			// Line: i + 1 assumes the first line in a csv file is headers
			return nil, &ErrLineMapping{Line: i + 1, Err: err}
		}

		registrations[i] = parsed
	}

	return registrations, nil
}

func (pc *ParsedCsv) ToSelectionDomain(kind types.SelectionKind, year int32, semester int32) (*ParsedSelection, error) {
	registrationStatus, _ := types.ParseRegistrationStatus(kind.ToRegistrationStatus().String())
	registrations, err := pc.toRegistrationsDomain(registrationStatus)
	if err != nil {
		return nil, err
	}

	if len(pc.candidates) == 0 {
		return nil, ErrNoRegistrationsFound{}
	}

	selection := &types.Selection{
		Name:        pc.name,
		Kind:        kind,
		Year:        year,
		Semester:    semester,
		Institution: pc.candidates[0].Institution,
		Degree:      pc.candidates[0].Course,
	}

	return &ParsedSelection{
		Selection:     selection,
		Registrations: registrations,
	}, nil
}

type csvCandidate struct {
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
func (a *csvCandidate) Parse(status types.RegistrationStatus) (*ParsedRegistration, error) {
	period, err := parsePeriod(a.SchedulePeriod)
	if err != nil {
		return nil, &ErrFieldValidation{Field: "Period", Err: err}
	}

	// The waitlist file has several MinimumScore and Ranking fields empt
	// Should we only allow this in the waitlist file?
	ms := a.MinimumScore
	if a.MinimumScore == "" {
		ms = "0"
	}

	r := a.Ranking
	if a.Ranking == "" {
		r = "0"
	}

	minimumScore, err := parseScoreString(ms)
	if err != nil {
		return nil, &ErrFieldValidation{Field: "MinimumScore", Err: err}
	}

	seats, err := strconv.ParseInt(a.Seats, 10, 32)
	if err != nil {
		return nil, &ErrFieldValidation{Field: "Seats", Err: err}
	}

	option, err := strconv.ParseInt(a.Option, 10, 32)
	if err != nil {
		return nil, &ErrFieldValidation{Field: "Option", Err: err}
	}

	languagesScore, err := parseScoreString(a.LanguagesScore)
	if err != nil {
		return nil, &ErrFieldValidation{Field: "LanguagesScore", Err: err}
	}

	humanitiesScore, err := parseScoreString(a.HumanitiesScore)
	if err != nil {
		return nil, &ErrFieldValidation{Field: "HumanitiesScore", Err: err}
	}

	naturalSciencesScore, err := parseScoreString(a.NaturalSciencesScore)
	if err != nil {
		return nil, &ErrFieldValidation{Field: "NaturalSciencesScore", Err: err}
	}

	mathematicsScore, err := parseScoreString(a.MathematicsScore)
	if err != nil {
		return nil, &ErrFieldValidation{Field: "MathematicsScore", Err: err}
	}

	essayScore, err := parseScoreString(a.EssayScore)
	if err != nil {
		return nil, &ErrFieldValidation{Field: "EssayScore", Err: err}
	}

	compositeScore, err := parseScoreString(a.CompositeScore)
	if err != nil {
		return nil, &ErrFieldValidation{Field: "CompositeScore", Err: err}
	}

	ranking, err := strconv.Atoi(r)
	if err != nil {
		return nil, &ErrFieldValidation{Field: "Ranking", Err: err}
	}

	var call *types.Call
	if status == types.RegistrationStatusApproved {
		call = &types.Call{}
	}

	return &ParsedRegistration{
		Registration: &types.Registration{
			Status:               status,
			EnrollmentID:         a.EnrollmentID,
			Option:               int32(option),
			LanguagesScore:       languagesScore,
			HumanitiesScore:      humanitiesScore,
			NaturalSciencesScore: naturalSciencesScore,
			MathematicsScore:     mathematicsScore,
			EssayScore:           essayScore,
			CompositeScore:       compositeScore,
			Ranking:              int32(ranking),
			Candidate: &types.Candidate{
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
		},
		Course: &types.Course{
			Period:       *period,
			Quota:        a.Quota,
			Seats:        int32(seats),
			MinimumScore: minimumScore,
		},
		Call: call,
	}, nil
}

func parseScoreString(value string) (*types.Score, error) {
	scoreFloat, err := strconv.ParseFloat(strings.ReplaceAll(value, ",", "."), 32)
	if err != nil {
		return nil, err
	}

	score := types.NewScoreFromFloat(float32(scoreFloat))

	return &score, nil
}

func parsePeriod(name string) (*types.CoursePeriod, error) {
	periods := map[string]types.CoursePeriod{
		"Matutino": types.CoursePeriodMorning,
		"Noturno":  types.CoursePeriodEvening,
	}

	if x, ok := periods[name]; ok {
		return &x, nil
	}

	return nil, fmt.Errorf("%s is %w", name, ErrInvalidPeriod{Period: name})
}
