package database

import (
	"github.com/baldugus/sisu/database/.gen/model"
	"github.com/baldugus/sisu/types"
)

type selectionResult struct {
	model.Selections
}

func (s *selectionResult) toSelectionDomain() *types.Selection {
	kind, _ := types.ParseSelectionKind(s.Kind)

	return &types.Selection{
		ID:          s.ID,
		Kind:        kind,
		Name:        s.Name,
		Year:        s.Year,
		Institution: s.Institution,
		Degree:      s.Degree,
	}
}

type registrationsResult []registrationResult

func (r registrationsResult) toRegistrationsDomain() []*types.Registration {
	registrations := make([]*types.Registration, len(r))
	for i := range r {
		registrations[i] = r[i].toRegistrationDomain()
	}

	return registrations
}

type registrationResult struct {
	model.Registrations

	Candidate model.Candidates
}

type fullRegistrationResult struct {
	model.Registrations

	Candidate model.Candidates
	Course    model.Courses
	Quota     model.Quotas
	Call      *model.Calls
}

type fullRegistrationsResult []fullRegistrationResult

func (r fullRegistrationsResult) toRegistrationDetails() []*types.RegistrationDetail {
	details := make([]*types.RegistrationDetail, len(r))
	for i := range r {
		details[i] = r[i].toRegistrationDetail()
	}
	return details
}

func (r *fullRegistrationResult) toRegistrationDetail() *types.RegistrationDetail {
	status, _ := types.ParseRegistrationStatus(r.Status)

	return &types.RegistrationDetail{
		Registration: &types.Registration{
			ID:                   r.ID,
			EnrollmentID:         r.EnrollmentID,
			Option:               r.Option,
			LanguagesScore:       &types.Score{Value: r.LanguagesScore},
			HumanitiesScore:      &types.Score{Value: r.HumanitiesScore},
			NaturalSciencesScore: &types.Score{Value: r.NaturalSciencesScore},
			MathematicsScore:     &types.Score{Value: r.MathematicsScore},
			EssayScore:           &types.Score{Value: r.EssayScore},
			CompositeScore:       &types.Score{Value: r.CompositeScore},
			Ranking:              r.Ranking,
			Status:               status,
			Candidate:            toCandidateDomain(&r.Candidate),
		},
		Course: toCourseDomain(&r.Course, r.Quota.Name),
		Call:   toCallDomain(r.Call),
	}
}

func (r *registrationResult) toRegistrationDomain() *types.Registration {
	status, _ := types.ParseRegistrationStatus(r.Status)

	return &types.Registration{
		ID:                   r.ID,
		EnrollmentID:         r.EnrollmentID,
		Option:               r.Option,
		LanguagesScore:       &types.Score{Value: r.LanguagesScore},
		HumanitiesScore:      &types.Score{Value: r.HumanitiesScore},
		NaturalSciencesScore: &types.Score{Value: r.NaturalSciencesScore},
		MathematicsScore:     &types.Score{Value: r.MathematicsScore},
		EssayScore:           &types.Score{Value: r.EssayScore},
		CompositeScore:       &types.Score{Value: r.CompositeScore},
		Ranking:              r.Ranking,
		Status:               status,
		Candidate:            toCandidateDomain(&r.Candidate),
	}
}

func toCandidateDomain(candidate *model.Candidates) *types.Candidate {
	return &types.Candidate{
		ID:           candidate.ID,
		CPF:          candidate.Cpf,
		Name:         candidate.Name,
		SocialName:   candidate.SocialName,
		BirthDate:    candidate.Birthdate,
		Sex:          candidate.Sex,
		MotherName:   candidate.MotherName,
		AddressLine:  candidate.AddressLine,
		AddressLine2: candidate.AddressLine2,
		HouseNumber:  candidate.HouseNumber,
		Neighborhood: candidate.Neighborhood,
		Municipality: candidate.Municipality,
		State:        candidate.State,
		CEP:          candidate.Cep,
		Email:        candidate.Email,
		Phone1:       candidate.Phone1,
		Phone2:       candidate.Phone2,
	}
}

func toCallDomain(c *model.Calls) *types.Call {
	if c == nil {
		return nil
	}

	status, _ := types.ParseCallStatus(c.Status)

	return &types.Call{
		ID:         c.ID,
		Status:     status,
		Number:     c.Number,
		SemesterID: c.SemesterID,
	}
}

type callsResult []*model.Calls

func (c callsResult) toCallsDomain() []*types.Call {
	calls := make([]*types.Call, len(c))

	for i, call := range c {
		calls[i] = toCallDomain(call)
	}

	return calls
}

func toCandidateModel(candidate *types.Candidate) *model.Candidates {
	return &model.Candidates{
		Cpf:          candidate.CPF,
		Name:         candidate.Name,
		SocialName:   candidate.SocialName,
		Birthdate:    candidate.BirthDate,
		Sex:          candidate.Sex,
		MotherName:   candidate.MotherName,
		AddressLine:  candidate.AddressLine,
		AddressLine2: candidate.AddressLine2,
		HouseNumber:  candidate.HouseNumber,
		Neighborhood: candidate.Neighborhood,
		Municipality: candidate.Municipality,
		State:        candidate.State,
		Cep:          candidate.CEP,
		Email:        candidate.Email,
		Phone1:       candidate.Phone1,
		Phone2:       candidate.Phone2,
	}
}

func toCallModel(call *types.Call) *model.Calls {
	return &model.Calls{
		Status:     call.Status.String(),
		Number:     call.Number,
		SemesterID: call.SemesterID,
	}
}

// TODO: enable strict scan in JET.
func toSelectionModel(selection *types.Selection) *model.Selections {
	return &model.Selections{
		Kind:        selection.Kind.String(),
		Name:        selection.Name,
		Year:        selection.Year,
		Institution: selection.Institution,
		Degree:      selection.Degree,
	}
}

func toCourseModel(course *types.Course) *model.Courses {
	return &model.Courses{
		TimeSlot:     course.Period.String(),
		Seats:        course.Seats,
		MinimumScore: course.MinimumScore.Value,
	}
}

func toCourseDomain(course *model.Courses, quotaName string) *types.Course {
	period, _ := types.ParseCoursePeriod(course.TimeSlot)

	return &types.Course{
		ID:           course.ID,
		Seats:        course.Seats,
		MinimumScore: &types.Score{Value: course.MinimumScore},
		Period:       period,
		Quota:        quotaName,
	}
}

func toRegistrationModel(registration *types.Registration) *model.Registrations {
	return &model.Registrations{
		EnrollmentID:         registration.EnrollmentID,
		Option:               registration.Option,
		LanguagesScore:       registration.LanguagesScore.Value,
		HumanitiesScore:      registration.HumanitiesScore.Value,
		NaturalSciencesScore: registration.NaturalSciencesScore.Value,
		MathematicsScore:     registration.MathematicsScore.Value,
		EssayScore:           registration.EssayScore.Value,
		CompositeScore:       registration.CompositeScore.Value,
		Ranking:              registration.Ranking,
		Status:               registration.Status.String(),
	}
}
