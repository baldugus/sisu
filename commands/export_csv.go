package commands

import (
	"fmt"
	"os"

	"github.com/gocarina/gocsv"

	"github.com/baldugus/sisu/database"
	"github.com/baldugus/sisu/types"
)

type ExportCSVCommand struct {
	FilePath string
}

func (cmd *ExportCSVCommand) Execute(db *database.Database) error {
	details, err := db.FetchEnrolledRegistrationDetails()
	if err != nil {
		return fmt.Errorf("fetch enrolled registration details: %w", err)
	}

	// Convert to flat CSV structure
	csvRecords := make([]*types.EnrolledRegistrationCSV, len(details))
	for i, detail := range details {
		csvRecords[i] = toEnrolledRegistrationCSV(detail)
	}

	file, err := os.Create(cmd.FilePath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer file.Close()

	if err := gocsv.Marshal(csvRecords, file); err != nil {
		return fmt.Errorf("csv marshal: %w", err)
	}

	return nil
}

func toEnrolledRegistrationCSV(detail *types.RegistrationDetail) *types.EnrolledRegistrationCSV {
	reg := detail.Registration
	candidate := reg.Candidate
	course := detail.Course

	return &types.EnrolledRegistrationCSV{
		// Candidate fields
		CPF:          candidate.CPF,
		Name:         candidate.Name,
		SocialName:   candidate.SocialName,
		BirthDate:    candidate.BirthDate,
		Sex:          candidate.Sex,
		MotherName:   candidate.MotherName,
		AddressLine:  candidate.AddressLine,
		AddressLine2: candidate.AddressLine2,
		HouseNumber:  candidate.HouseNumber,
		Neighborhood: candidate.Neighborhood,
		Municipality: candidate.Municipality,
		State:        candidate.State,
		CEP:          candidate.CEP,
		Email:        candidate.Email,
		Phone1:       candidate.Phone1,
		Phone2:       candidate.Phone2,

		// Registration fields
		EnrollmentID:         reg.EnrollmentID,
		Option:               reg.Option,
		LanguagesScore:       reg.LanguagesScore.String(),
		HumanitiesScore:      reg.HumanitiesScore.String(),
		NaturalSciencesScore: reg.NaturalSciencesScore.String(),
		MathematicsScore:     reg.MathematicsScore.String(),
		EssayScore:           reg.EssayScore.String(),
		CompositeScore:       reg.CompositeScore.String(),
		Ranking:              reg.Ranking,
		Status:               registrationStatusToPortuguese(reg.Status),

		// Course fields
		Period: coursePeriodToPortugueseCSV(course.Period),
		Quota:  course.Quota,
	}
}

func registrationStatusToPortuguese(status types.RegistrationStatus) string {
	switch status {
	case types.RegistrationStatusApproved:
		return "Aprovado"
	case types.RegistrationStatusWaitlisted:
		return "Em espera"
	case types.RegistrationStatusAbsent:
		return "Faltoso"
	case types.RegistrationStatusEnrolled:
		return "Matriculado"
	default:
		return status.String()
	}
}

func coursePeriodToPortugueseCSV(period types.CoursePeriod) string {
	switch period {
	case types.CoursePeriodMorning:
		return "Matutino"
	case types.CoursePeriodEvening:
		return "Noturno"
	default:
		return period.String()
	}
}
