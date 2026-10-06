package database

import (
	"fmt"

	"github.com/baldugus/sisu/database/.gen/model"
	"github.com/baldugus/sisu/types"
)

// derivePlacement turns a registration's current placement (from the
// registration_placements view) into its status and semester.
func derivePlacement(p *model.RegistrationPlacements) (types.RegistrationStatus, *int32) {
	if p == nil || p.Outcome == nil {
		return types.RegistrationStatusWaitlisted, nil
	}

	switch types.CallEntryOutcome(*p.Outcome) {
	case types.CallEntryOutcomePending:
		return types.RegistrationStatusApproved, p.Semester
	case types.CallEntryOutcomeEnrolled:
		return types.RegistrationStatusEnrolled, p.Semester
	default:
		return types.RegistrationStatusAbsent, nil
	}
}

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
	Placement *model.RegistrationPlacements
}

type fullRegistrationResult struct {
	model.Registrations

	Candidate model.Candidates
	Placement *model.RegistrationPlacements
	Course    model.Courses
	Quota     model.Quotas
}

type fullRegistrationsResult []fullRegistrationResult

func (r fullRegistrationsResult) toRegistrationDetails() ([]*types.RegistrationDetail, error) {
	details := make([]*types.RegistrationDetail, len(r))
	for i := range r {
		detail, err := r[i].toRegistrationDetail()
		if err != nil {
			return nil, err
		}
		details[i] = detail
	}
	return details, nil
}

func (r *fullRegistrationResult) toRegistrationDetail() (*types.RegistrationDetail, error) {
	course, err := toCourseDomain(&r.Course, r.Quota.Name)
	if err != nil {
		return nil, err
	}

	return &types.RegistrationDetail{
		Registration: toRegistrationDomain(&r.Registrations, &r.Candidate, r.Placement),
		Course:       course,
	}, nil
}

func (r *registrationResult) toRegistrationDomain() *types.Registration {
	return toRegistrationDomain(&r.Registrations, &r.Candidate, r.Placement)
}

func toRegistrationDomain(
	r *model.Registrations,
	candidate *model.Candidates,
	placement *model.RegistrationPlacements,
) *types.Registration {
	status, semester := derivePlacement(placement)

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
		Semester:             semester,
		Candidate:            toCandidateDomain(candidate),
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
		ID:     c.ID,
		Status: status,
		Number: c.Number,
	}
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
		Status: call.Status.String(),
		Number: call.Number,
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
		Seats:        course.Seats.Total(),
		MinimumScore: course.MinimumScore.Value,
	}
}

func toCourseDomain(course *model.Courses, quotaName string) (*types.Course, error) {
	period, _ := types.ParseCoursePeriod(course.TimeSlot)

	// The courses.seats CHECK constraint keeps odd seat counts out of the
	// database, so this only fails on a database written outside the app.
	seats, err := types.NewSeats(course.Seats)
	if err != nil {
		return nil, fmt.Errorf("map course %d: %w", course.ID, err)
	}

	return &types.Course{
		ID:           course.ID,
		Seats:        seats,
		MinimumScore: &types.Score{Value: course.MinimumScore},
		Period:       period,
		Quota:        quotaName,
	}, nil
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
	}
}

func toCallEntryDomain(e *model.CallEntries, callNumber int32) *types.CallEntry {
	return &types.CallEntry{
		CallID:         e.CallID,
		CallNumber:     callNumber,
		RegistrationID: e.RegistrationID,
		Kind:           types.CallEntryKind(e.Kind),
		Semester:       e.Semester,
		Outcome:        types.CallEntryOutcome(e.Outcome),
		WantsPromotion: e.WantsPromotion != 0,
	}
}

func toCallEntryModel(e *types.CallEntry) *model.CallEntries {
	var wants int32
	if e.WantsPromotion {
		wants = 1
	}

	return &model.CallEntries{
		CallID:         e.CallID,
		RegistrationID: e.RegistrationID,
		Kind:           e.Kind.String(),
		Semester:       e.Semester,
		Outcome:        e.Outcome.String(),
		WantsPromotion: wants,
	}
}
