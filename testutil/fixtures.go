package testutil

import (
	"github.com/baldugus/sisu/types"
)

// NewTestSelection creates a test Selection with sensible defaults.
// Pass override functions to customize specific fields.
func NewTestSelection(overrides ...func(*types.Selection)) *types.Selection {
	selection := &types.Selection{
		Name:        "Test Selection",
		Kind:        types.SelectionKindApproved,
		Year:        2025,
		Semester:    1,
		Institution: "FAETERJ-Rio",
		Degree:      "Tecnológico",
	}

	for _, override := range overrides {
		override(selection)
	}

	return selection
}

// NewTestCandidate creates a test Candidate with sensible defaults.
func NewTestCandidate(overrides ...func(*types.Candidate)) *types.Candidate {
	candidate := &types.Candidate{
		CPF:          "12345678901",
		Name:         "Test Candidate",
		SocialName:   "Test Candidate",
		BirthDate:    "2000-01-01 00:00:00",
		Sex:          "M",
		MotherName:   "Test Mother",
		AddressLine:  "Test Street",
		AddressLine2: "Test Complement",
		HouseNumber:  "123",
		Neighborhood: "Test Neighborhood",
		Municipality: "Test City",
		State:        "RJ",
		CEP:          "12345678",
		Phone1:       "1234567890",
		Phone2:       "0987654321",
		Email:        "test@example.com",
	}

	for _, override := range overrides {
		override(candidate)
	}

	return candidate
}

// NewTestCourse creates a test Course with sensible defaults.
func NewTestCourse(overrides ...func(*types.Course)) *types.Course {
	course := &types.Course{
		Period:       types.CoursePeriodMorning,
		Seats:        20,
		MinimumScore: &types.Score{Value: 50000}, // 500.00 in Score format
		Quota:        "Ampla concorrência",
	}

	for _, override := range overrides {
		override(course)
	}

	return course
}

// NewTestCall creates a test Call with sensible defaults.
func NewTestCall(overrides ...func(*types.Call)) *types.Call {
	call := &types.Call{
		Number: 1,
		Status: types.CallStatusCalling,
	}

	for _, override := range overrides {
		override(call)
	}

	return call
}

// NewTestRegistration creates a test Registration with sensible defaults.
func NewTestRegistration(overrides ...func(*types.Registration)) *types.Registration {
	registration := &types.Registration{
		EnrollmentID:         "241077018570",
		Option:               1,
		LanguagesScore:       &types.Score{Value: 70000}, // 700.00 in Score format
		HumanitiesScore:      &types.Score{Value: 75000}, // 750.00 in Score format
		NaturalSciencesScore: &types.Score{Value: 68000}, // 680.00 in Score format
		MathematicsScore:     &types.Score{Value: 72000}, // 720.00 in Score format
		EssayScore:           &types.Score{Value: 80000}, // 800.00 in Score format
		CompositeScore:       &types.Score{Value: 73000}, // 730.00 in Score format
		Ranking:              1,
		Status:               types.RegistrationStatusApproved,
		Candidate:            NewTestCandidate(),
	}

	for _, override := range overrides {
		override(registration)
	}

	return registration
}
