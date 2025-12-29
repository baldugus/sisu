package csvparser

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/baldugus/sisu/types"
)

func checkExpectedErr(t testing.TB, err error, expectedErr func(error) bool) {
	t.Helper()
	if expectedErr != nil {
		if err == nil {
			t.Errorf("expected error but got none")
			return
		}
		if !expectedErr(err) {
			t.Errorf("unexpected error: %v", err)
		}
	} else {
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}
}

func TestParseCSVFile(t *testing.T) {
	tests := []struct {
		name        string
		filename    string
		wantLen     int
		expectedErr func(error) bool
		firstName   string
		firstCPF    string
		institution string
		course      string
	}{
		{
			name:        "valid mini csv",
			filename:    "testdata/mock_data_mini.csv",
			wantLen:     4,
			expectedErr: nil,
			firstName:   "Crysta Radborn",
			firstCPF:    "72503459660",
			institution: "FAETERJ-Rio",
			course:      "ANÁLISE E DESENVOLVIMENTO DE SISTEMAS",
		},
		{
			name:     "nonexistent file",
			filename: "testdata/nonexistent.csv",
			wantLen:  0,
			expectedErr: func(err error) bool {
				var e *ErrFileNotFound
				return errors.As(err, &e)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFile(tt.filename)
			checkExpectedErr(t, err, tt.expectedErr)
			if tt.expectedErr != nil {
				return
			}

			if len(got.candidates) != tt.wantLen {
				t.Errorf("expected %d candidates, got %d", tt.wantLen, len(got.candidates))
			}

			if tt.wantLen > 0 {
				if got.candidates[0].Name != tt.firstName {
					t.Errorf("expected first candidate name %q, got %q", tt.firstName, got.candidates[0].Name)
				}
				if got.candidates[0].CPF != tt.firstCPF {
					t.Errorf("expected first candidate CPF %q, got %q", tt.firstCPF, got.candidates[0].CPF)
				}
				if got.candidates[0].Institution != tt.institution {
					t.Errorf("expected institution %q, got %q", tt.institution, got.candidates[0].Institution)
				}
				if got.candidates[0].Course != tt.course {
					t.Errorf("expected course %q, got %q", tt.course, got.candidates[0].Course)
				}
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantLen     int
		expectedErr func(error) bool
	}{
		{
			name:        "valid CSV with data",
			input:       "NU_ETAPA;DS_TURNO;QT_VAGAS_CONCORRENCIA;CO_INSCRICAO_ENEM;NO_INSCRITO;NU_CPF_INSCRITO\nCHAMADA REGULAR;Matutino;20;123456789;Test User;12345678901",
			wantLen:     1,
			expectedErr: nil,
		},
		{
			name:    "empty CSV",
			input:   "",
			wantLen: 0,
			expectedErr: func(err error) bool {
				var e *ErrFileEmpty
				return errors.As(err, &e)
			},
		},
		{
			name:        "only headers",
			input:       "NU_ETAPA;DS_TURNO;QT_VAGAS_CONCORRENCIA;CO_INSCRICAO_ENEM;NO_INSCRITO;NU_CPF_INSCRITO",
			wantLen:     0,
			expectedErr: nil,
		},
		{
			name:    "malformed CSV - unterminated quote",
			input:   "NU_ETAPA;DS_TURNO;QT_VAGAS_CONCORRENCIA\n\"CHAMADA REGULAR;Matutino;20",
			wantLen: 0,
			expectedErr: func(err error) bool {
				var e *ErrFileParse
				return errors.As(err, &e)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			got, err := Parse(reader, "test.csv")
			checkExpectedErr(t, err, tt.expectedErr)
			if tt.expectedErr != nil {
				return
			}

			if len(got.candidates) != tt.wantLen {
				t.Errorf("expected %d candidates, got %d", tt.wantLen, len(got.candidates))
			}
		})
	}
}

func TestRegistrations(t *testing.T) {
	tests := []struct {
		name           string
		pc             *ParsedCsv
		status         types.RegistrationStatus
		expectedOption int32
		expectedScore  int32
		expectedErr    func(error) bool
	}{
		{
			name: "valid applicant",
			pc: &ParsedCsv{
				candidates: []*csvCandidate{
					{
						Name:                 "Crysta Radborn",
						CPF:                  "72503459660",
						EnrollmentID:         "241077018570",
						Option:               "2",
						CompositeScore:       "655,17",
						MinimumScore:         "600,0",
						Seats:                "20",
						Ranking:              "100",
						LanguagesScore:       "500,5",
						HumanitiesScore:      "600,5",
						NaturalSciencesScore: "700,5",
						MathematicsScore:     "800,5",
						EssayScore:           "900,5",
						SchedulePeriod:       "Matutino",
					},
				},
				name: "test.csv",
			},
			status:         types.RegistrationStatusApproved,
			expectedOption: 2,
			expectedScore:  65517,
			expectedErr:    nil,
		},
		{
			name: "invalid option",
			pc: &ParsedCsv{
				candidates: []*csvCandidate{
					{
						MinimumScore:   "600,0",
						Seats:          "20",
						Option:         "invalid",
						SchedulePeriod: "Noturno",
					},
				},
				name: "test.csv",
			},
			status:         types.RegistrationStatusApproved,
			expectedOption: 0,
			expectedScore:  0,
			expectedErr: func(err error) bool {
				var e *ErrFieldValidation
				return errors.As(err, &e) && e.Field == "Option"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.pc.toRegistrationsDomain(tt.status)
			checkExpectedErr(t, err, tt.expectedErr)
			if tt.expectedErr != nil {
				return
			}

			expectedLen := len(tt.pc.candidates)
			if len(got) != expectedLen {
				t.Errorf("expected %d applications, got %d", expectedLen, len(got))
			}

			if len(got) > 0 {
				firstApp := got[0]
				if firstApp.Registration.Status != tt.status {
					t.Errorf("expected status %q, got %q", tt.status, firstApp.Registration.Status)
				}
				if firstApp.Registration.Candidate.Name != tt.pc.candidates[0].Name {
					t.Errorf("expected candidate name %q, got %q", tt.pc.candidates[0].Name, firstApp.Registration.Candidate.Name)
				}
				if firstApp.Registration.Candidate.CPF != tt.pc.candidates[0].CPF {
					t.Errorf("expected candidate CPF %q, got %q", tt.pc.candidates[0].CPF, firstApp.Registration.Candidate.CPF)
				}
				if firstApp.Registration.EnrollmentID != tt.pc.candidates[0].EnrollmentID {
					t.Errorf("expected enrollment ID %q, got %q", tt.pc.candidates[0].EnrollmentID, firstApp.Registration.EnrollmentID)
				}
				if firstApp.Registration.Option != tt.expectedOption {
					t.Errorf("expected option %d, got %d", tt.expectedOption, firstApp.Registration.Option)
				}
				if firstApp.Registration.CompositeScore.Value != tt.expectedScore {
					t.Errorf("expected composite score %d, got %d", tt.expectedScore, firstApp.Registration.CompositeScore.Value)
				}
			}
		})
	}
}

func TestSanitizeCsv(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "BOM removal and header conversion",
			input:    "\ufeffheader1,header2\rline1;line2",
			expected: "header1;header2\nline1;line2",
		},
		{
			name:     "header and CR to LF conversion",
			input:    "header1,header2\rline1;line2\rline3;line4",
			expected: "header1;header2\nline1;line2\nline3;line4",
		},
		{
			name:     "no changes needed",
			input:    "header1;header2\nline1;line2",
			expected: "header1;header2\nline1;line2",
		},
		{
			name:     "empty input",
			input:    "",
			expected: "",
		},
		{
			name:     "only header",
			input:    "header1,header2",
			expected: "header1;header2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			result, err := sanitizeCsv(reader)
			if err != nil {
				t.Fatalf("sanitizeCsv() error = %v", err)
			}

			resultBytes, err := io.ReadAll(result)
			if err != nil {
				t.Fatalf("failed to read result: %v", err)
			}

			if string(resultBytes) != tt.expected {
				t.Errorf("sanitizeCsv() = %q, want %q", string(resultBytes), tt.expected)
			}
		})
	}
}

func TestToSelection(t *testing.T) {
	tests := []struct {
		name         string
		candidates   []*csvCandidate
		filename     string
		kind         types.SelectionKind
		year         int32
		semester     int32
		wantName     string
		wantYear     int32
		wantSemester int32
		wantInst     string
		wantCourse   string
		wantErr      bool
	}{
		{
			name: "with candidates",
			candidates: []*csvCandidate{
				{
					Institution:          "FAETERJ-Rio",
					Course:               "Computer Science",
					SchedulePeriod:       "Matutino",
					Seats:                "20",
					Option:               "1",
					MinimumScore:         "600,0",
					Ranking:              "1",
					LanguagesScore:       "500,0",
					HumanitiesScore:      "500,0",
					NaturalSciencesScore: "500,0",
					MathematicsScore:     "500,0",
					EssayScore:           "500,0",
					CompositeScore:       "500,0",
				},
			},
			filename:     "test.csv",
			kind:         types.SelectionKindApproved,
			year:         2024,
			semester:     1,
			wantName:     "test.csv",
			wantYear:     2024,
			wantSemester: 1,
			wantInst:     "FAETERJ-Rio",
			wantCourse:   "Computer Science",
			wantErr:      false,
		},
		{
			name:       "empty candidates",
			candidates: []*csvCandidate{},
			filename:   "empty.csv",
			kind:       types.SelectionKindApproved,
			year:       2024,
			semester:   1,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pc := &ParsedCsv{
				candidates: tt.candidates,
				name:       tt.filename,
			}

			got, err := pc.ToSelectionDomain(tt.kind, tt.year, tt.semester)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("ToSelectionDomain() error = %v", err)
			}

			if got.Selection.Name != tt.wantName {
				t.Errorf("expected name %q, got %q", tt.wantName, got.Selection.Name)
			}
			if got.Selection.Kind != tt.kind {
				t.Errorf("expected kind %v, got %v", tt.kind, got.Selection.Kind)
			}
			if got.Selection.Year != tt.wantYear {
				t.Errorf("expected year %d, got %d", tt.wantYear, got.Selection.Year)
			}
			if got.Selection.Semester != tt.wantSemester {
				t.Errorf("expected semester %d, got %d", tt.wantSemester, got.Selection.Semester)
			}
			if got.Selection.Institution != tt.wantInst {
				t.Errorf("expected institution %q, got %q", tt.wantInst, got.Selection.Institution)
			}
			if got.Selection.Degree != tt.wantCourse {
				t.Errorf("expected course %q, got %q", tt.wantCourse, got.Selection.Degree)
			}
		})
	}
}

func TestParseLocalizedFloat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int32
		hasError bool
	}{
		{"valid comma decimal", "123,45", 12345, false},
		{"valid dot decimal", "123.45", 12345, false},
		{"zero", "0,0", 0, false},
		{"integer", "123", 12300, false},
		{"empty string", "", 0, true},
		{"invalid number", "invalid", 0, true},
		{"multiple commas", "1,23,45", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseScoreString(tt.input)
			if (err != nil) != tt.hasError {
				t.Errorf("parseLocalizedFloat() error = %v, wantErr %v", err, tt.hasError)
				return
			}
			if !tt.hasError && got.Value != tt.expected {
				t.Errorf("parseLocalizedFloat() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCsvCandidateParse(t *testing.T) {
	tests := []struct {
		name        string
		candidate   csvCandidate
		status      types.RegistrationStatus
		expectedErr func(error) bool
	}{
		{
			name: "valid candidate",
			candidate: csvCandidate{
				Stage:                "CHAMADA REGULAR",
				SchedulePeriod:       "Matutino",
				Seats:                "20",
				EnrollmentID:         "123456789",
				Name:                 "Test User",
				CPF:                  "12345678901",
				Date:                 "2024-01-01",
				SocialName:           "",
				BirthDate:            "1990-01-01",
				Sex:                  "M",
				MotherName:           "Test Mother",
				AddressLine:          "Test Street",
				HouseNumber:          "123",
				AddressLine2:         "Apt 1",
				State:                "RJ",
				Municipality:         "Rio de Janeiro",
				Neighborhood:         "Test Neighborhood",
				CEP:                  "12345678",
				Phone1:               "123456789",
				Phone2:               "987654321",
				Email:                "test@example.com",
				LanguagesScore:       "500,5",
				HumanitiesScore:      "600,5",
				NaturalSciencesScore: "700,5",
				MathematicsScore:     "800,5",
				EssayScore:           "900,5",
				Option:               "1",
				Quota:                "Ampla concorrência",
				CompositeScore:       "650,5",
				MinimumScore:         "600,0",
				Ranking:              "100",
				Institution:          "FAETERJ-Rio",
				Course:               "Computer Science",
			},
			status:      types.RegistrationStatusApproved,
			expectedErr: nil,
		},
		{
			name: "invalid seats",
			candidate: csvCandidate{
				MinimumScore:   "600,0",
				Seats:          "invalid",
				SchedulePeriod: "Noturno",
			},
			status: types.RegistrationStatusApproved,
			expectedErr: func(err error) bool {
				var e *ErrFieldValidation
				return errors.As(err, &e) && e.Field == "Seats"
			},
		},
		{
			name: "invalid option",
			candidate: csvCandidate{
				MinimumScore:   "600,0",
				Seats:          "20",
				Option:         "invalid",
				SchedulePeriod: "Noturno",
			},
			status: types.RegistrationStatusApproved,
			expectedErr: func(err error) bool {
				var e *ErrFieldValidation
				return errors.As(err, &e) && e.Field == "Option"
			},
		},
		{
			name: "invalid ranking",
			candidate: csvCandidate{
				MinimumScore:         "600,0",
				Seats:                "20",
				Option:               "1",
				LanguagesScore:       "500,5",
				HumanitiesScore:      "600,5",
				NaturalSciencesScore: "700,5",
				MathematicsScore:     "800,5",
				EssayScore:           "900,5",
				CompositeScore:       "650,5",
				Ranking:              "invalid",
				SchedulePeriod:       "Noturno",
			},
			status: types.RegistrationStatusApproved,
			expectedErr: func(err error) bool {
				var e *ErrFieldValidation
				return errors.As(err, &e) && e.Field == "Ranking"
			},
		},
		{
			name: "invalid languages score",
			candidate: csvCandidate{
				MinimumScore:         "600,0",
				Seats:                "20",
				Option:               "1",
				LanguagesScore:       "invalid",
				HumanitiesScore:      "600,5",
				NaturalSciencesScore: "700,5",
				MathematicsScore:     "800,5",
				EssayScore:           "900,5",
				CompositeScore:       "650,5",
				Ranking:              "100",
				SchedulePeriod:       "Noturno",
			},
			status: types.RegistrationStatusApproved,
			expectedErr: func(err error) bool {
				var e *ErrFieldValidation
				return errors.As(err, &e) && e.Field == "LanguagesScore"
			},
		},
		{
			name: "invalid minimum score",
			candidate: csvCandidate{
				Seats:          "20",
				Option:         "1",
				Ranking:        "100",
				MinimumScore:   "invalid",
				SchedulePeriod: "Noturno",
			},
			status: types.RegistrationStatusApproved,
			expectedErr: func(err error) bool {
				var e *ErrFieldValidation
				return errors.As(err, &e) && e.Field == "MinimumScore"
			},
		},
		{
			name: "invalid period",
			candidate: csvCandidate{
				SchedulePeriod: "invalid",
			},
			status: types.RegistrationStatusApproved,
			expectedErr: func(err error) bool {
				var e *ErrFieldValidation
				return errors.As(err, &e) && e.Field == "Period"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.candidate.Parse(tt.status)
			checkExpectedErr(t, err, tt.expectedErr)
			if tt.expectedErr != nil {
				return
			}

			if got.Registration.Status != tt.status {
				t.Errorf("expected status %q, got %q", tt.status, got.Registration.Status)
			}

			if got.Registration.EnrollmentID != tt.candidate.EnrollmentID {
				t.Errorf("expected enrollment ID %q, got %q", tt.candidate.EnrollmentID, got.Registration.EnrollmentID)
			}
		})
	}
}
