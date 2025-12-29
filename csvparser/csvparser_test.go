package csvparser

import (
	"errors"
	"io"
	"strconv"
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
				var e *FileNotFoundError
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

			if len(got.applicants) != tt.wantLen {
				t.Errorf("expected %d applicants, got %d", tt.wantLen, len(got.applicants))
			}

			if tt.wantLen > 0 {
				if got.applicants[0].Name != tt.firstName {
					t.Errorf("expected first applicant name %q, got %q", tt.firstName, got.applicants[0].Name)
				}
				if got.applicants[0].CPF != tt.firstCPF {
					t.Errorf("expected first applicant CPF %q, got %q", tt.firstCPF, got.applicants[0].CPF)
				}
				if got.applicants[0].Institution != tt.institution {
					t.Errorf("expected institution %q, got %q", tt.institution, got.applicants[0].Institution)
				}
				if got.applicants[0].Course != tt.course {
					t.Errorf("expected course %q, got %q", tt.course, got.applicants[0].Course)
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
				var e *EmptyError
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
				var e *ParseError
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

			if len(got.applicants) != tt.wantLen {
				t.Errorf("expected %d applicants, got %d", tt.wantLen, len(got.applicants))
			}
		})
	}
}

func TestToApplications(t *testing.T) {
	tests := []struct {
		name        string
		pc          *ParsedCsv
		status      string
		expectedErr func(error) bool
	}{
		{
			name: "valid applicant",
			pc: &ParsedCsv{
				applicants: []*csvApplicant{
					{
						Name:                 "Crysta Radborn",
						CPF:                  "72503459660",
						EnrollmentID:         "241077018570",
						Option:               "2",
						CompositeScore:       "687,83",
						MinimumScore:         "600,0",
						Seats:                "20",
						Ranking:              "100",
						LanguagesScore:       "500,5",
						HumanitiesScore:      "600,5",
						NaturalSciencesScore: "700,5",
						MathematicsScore:     "800,5",
						EssayScore:           "900,5",
					},
				},
				name: "test.csv",
			},
			status:      "APPROVED",
			expectedErr: nil,
		},
		{
			name: "invalid option",
			pc: &ParsedCsv{
				applicants: []*csvApplicant{
					{
						MinimumScore: "600,0",
						Seats:        "20",
						Option:       "invalid",
					},
				},
				name: "test.csv",
			},
			status: "APPROVED",
			expectedErr: func(err error) bool {
				var e *ValidationError
				return errors.As(err, &e) && e.Field == "Option"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.pc.ToApplications(tt.status)
			checkExpectedErr(t, err, tt.expectedErr)
			if tt.expectedErr != nil {
				return
			}

			expectedLen := len(tt.pc.applicants)
			if len(got) != expectedLen {
				t.Errorf("expected %d applications, got %d", expectedLen, len(got))
			}

			if len(got) > 0 {
				firstApp := got[0]
				if firstApp.Status != tt.status {
					t.Errorf("expected status %q, got %q", tt.status, firstApp.Status)
				}
				if firstApp.Applicant.Name != tt.pc.applicants[0].Name {
					t.Errorf("expected applicant name %q, got %q", tt.pc.applicants[0].Name, firstApp.Applicant.Name)
				}
				if firstApp.Applicant.CPF != tt.pc.applicants[0].CPF {
					t.Errorf("expected applicant CPF %q, got %q", tt.pc.applicants[0].CPF, firstApp.Applicant.CPF)
				}
				if firstApp.EnrollmentID != tt.pc.applicants[0].EnrollmentID {
					t.Errorf("expected enrollment ID %q, got %q", tt.pc.applicants[0].EnrollmentID, firstApp.EnrollmentID)
				}
				expectedOption, _ := strconv.Atoi(tt.pc.applicants[0].Option)
				if firstApp.Option != expectedOption {
					t.Errorf("expected option %d, got %d", expectedOption, firstApp.Option)
				}
				expectedScore, _ := parseLocalizedFloat(tt.pc.applicants[0].CompositeScore)
				if firstApp.CompositeScore != expectedScore {
					t.Errorf("expected composite score %f, got %f", expectedScore, firstApp.CompositeScore)
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
		name       string
		applicants []*csvApplicant
		filename   string
		kind       types.SelectionKind
		date       string
		wantName   string
		wantDate   string
		wantInst   string
		wantCourse string
	}{
		{
			name: "with applicants",
			applicants: []*csvApplicant{
				{
					Institution: "FAETERJ-Rio",
					Course:      "Computer Science",
				},
			},
			filename:   "test.csv",
			kind:       types.ApprovedSelection,
			date:       "2024-01-01",
			wantName:   "test.csv",
			wantDate:   "2024-01-01",
			wantInst:   "FAETERJ-Rio",
			wantCourse: "Computer Science",
		},
		{
			name:       "empty applicants",
			applicants: []*csvApplicant{},
			filename:   "empty.csv",
			kind:       types.ApprovedSelection,
			date:       "2024-01-01",
			wantName:   "empty.csv",
			wantDate:   "",
			wantInst:   "",
			wantCourse: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pc := &ParsedCsv{
				applicants: tt.applicants,
				name:       tt.filename,
			}

			got := pc.ToSelection(tt.kind, tt.date)

			if got.Name != tt.wantName {
				t.Errorf("expected name %q, got %q", tt.wantName, got.Name)
			}
			if got.Kind != tt.kind {
				t.Errorf("expected kind %v, got %v", tt.kind, got.Kind)
			}
			if got.Date != tt.wantDate {
				t.Errorf("expected date %q, got %q", tt.wantDate, got.Date)
			}
			if got.Institution != tt.wantInst {
				t.Errorf("expected institution %q, got %q", tt.wantInst, got.Institution)
			}
			if got.Course != tt.wantCourse {
				t.Errorf("expected course %q, got %q", tt.wantCourse, got.Course)
			}
		})
	}
}

func TestParseLocalizedFloat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
		hasError bool
	}{
		{"valid comma decimal", "123,45", 123.45, false},
		{"valid dot decimal", "123.45", 123.45, false},
		{"zero", "0,0", 0.0, false},
		{"integer", "123", 123.0, false},
		{"empty string", "", 0, true},
		{"invalid number", "invalid", 0, true},
		{"multiple commas", "1,23,45", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseLocalizedFloat(tt.input)
			if (err != nil) != tt.hasError {
				t.Errorf("parseLocalizedFloat() error = %v, wantErr %v", err, tt.hasError)
				return
			}
			if !tt.hasError && got != tt.expected {
				t.Errorf("parseLocalizedFloat() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCsvApplicantParse(t *testing.T) {
	tests := []struct {
		name        string
		applicant   csvApplicant
		status      string
		expectedErr func(error) bool
	}{
		{
			name: "valid applicant",
			applicant: csvApplicant{
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
			status:      "APPROVED",
			expectedErr: nil,
		},
		{
			name: "invalid seats",
			applicant: csvApplicant{
				MinimumScore: "600,0",
				Seats:        "invalid",
			},
			status: "APPROVED",
			expectedErr: func(err error) bool {
				var e *ValidationError
				return errors.As(err, &e) && e.Field == "Seats"
			},
		},
		{
			name: "invalid option",
			applicant: csvApplicant{
				MinimumScore: "600,0",
				Seats:        "20",
				Option:       "invalid",
			},
			status: "APPROVED",
			expectedErr: func(err error) bool {
				var e *ValidationError
				return errors.As(err, &e) && e.Field == "Option"
			},
		},
		{
			name: "invalid ranking",
			applicant: csvApplicant{
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
			},
			status: "APPROVED",
			expectedErr: func(err error) bool {
				var e *ValidationError
				return errors.As(err, &e) && e.Field == "Ranking"
			},
		},
		{
			name: "invalid languages score",
			applicant: csvApplicant{
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
			},
			status: "APPROVED",
			expectedErr: func(err error) bool {
				var e *ValidationError
				return errors.As(err, &e) && e.Field == "LanguagesScore"
			},
		},
		{
			name: "invalid minimum score",
			applicant: csvApplicant{
				Seats:        "20",
				Option:       "1",
				Ranking:      "100",
				MinimumScore: "invalid",
			},
			status: "APPROVED",
			expectedErr: func(err error) bool {
				var e *ValidationError
				return errors.As(err, &e) && e.Field == "MinimumScore"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.applicant.Parse(tt.status)
			checkExpectedErr(t, err, tt.expectedErr)
			if tt.expectedErr != nil {
				return
			}

			if got.Status != tt.status {
				t.Errorf("expected status %q, got %q", tt.status, got.Status)
			}

			if got.EnrollmentID != tt.applicant.EnrollmentID {
				t.Errorf("expected enrollment ID %q, got %q", tt.applicant.EnrollmentID, got.EnrollmentID)
			}
		})
	}
}
