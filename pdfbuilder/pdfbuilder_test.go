package pdfbuilder

import (
	"bytes"
	"image/png"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/baldugus/sisu/types"
	fitz "github.com/gen2brain/go-fitz"
)

func createMockSelection() *types.Selection {
	return &types.Selection{
		Name:        "Seleção de Teste",
		Kind:        types.SelectionKindApproved,
		Year:        2025,
		Semester:    1,
		Institution: "FAETERJ-Rio",
		Degree:      "Análise e Desenvolvimento de Sistemas",
	}
}

func createMockRegistrations() types.Registrations {
	return types.Registrations{
		{
			Status:         types.RegistrationStatusApproved,
			EnrollmentID:   "123456789",
			Option:         1,
			CompositeScore: &types.Score{Value: 85050},
			Ranking:        1,
			Candidate: &types.Candidate{
				Name:      "João Silva",
				Email:     "joao.silva@exemplo.com",
				CPF:       "12345678901",
				BirthDate: "1990-01-01",
			},
		},
		{
			Status:         types.RegistrationStatusApproved,
			EnrollmentID:   "987654321",
			Option:         1,
			CompositeScore: &types.Score{Value: 82030},
			Ranking:        2,
			Candidate: &types.Candidate{
				Name:      "Maria Santos",
				Email:     "maria.santos@exemplo.com",
				CPF:       "98765432109",
				BirthDate: "1992-05-15",
			},
		},
	}
}

func createMockCourseInfo() []*CourseInfo {
	registrations := createMockRegistrations()
	return []*CourseInfo{
		NewCourseInfo("Ampla concorrência", registrations),
	}
}

func createMockSelectionInfo() *SelectionInfo {
	sel := createMockSelection()
	return NewSelectionInfo(sel, 0)
}

func createMockBuilder() *Builder {
	return &Builder{
		Period:    "Manhã",
		Selection: createMockSelectionInfo(),
		Courses:   createMockCourseInfo(),
	}
}

func generateAndCompareGolden(t *testing.T, pdfPath, goldenPath string) {
	t.Helper()

	doc, err := fitz.New(pdfPath)
	if err != nil {
		t.Fatalf("Failed to open PDF: %v", err)
	}
	defer doc.Close()

	img, err := doc.Image(0) // First page
	if err != nil {
		t.Fatalf("Failed to convert PDF to image: %v", err)
	}

	// Encode to PNG bytes
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("Failed to encode PNG: %v", err)
	}
	newPNG := buf.Bytes()

	if _, err := os.Stat(goldenPath); os.IsNotExist(err) {
		// Golden doesn't exist, create it
		if err := os.WriteFile(goldenPath, newPNG, 0644); err != nil {
			t.Fatalf("Failed to write golden file: %v", err)
		}
		t.Logf("Golden file created: %s", goldenPath)
	} else {
		// Golden exists, compare
		goldenPNG, err := os.ReadFile(goldenPath)
		if err != nil {
			t.Fatalf("Failed to read golden file: %v", err)
		}
		if !reflect.DeepEqual(newPNG, goldenPNG) {
			t.Errorf("Generated image does not match golden file: %s", goldenPath)
		}
	}
}

func TestBuildWebsitePdf_Golden(t *testing.T) {
	builder := createMockBuilder()
	pdfPath := filepath.Join(t.TempDir(), "website.pdf")
	goldenPath := "testdata/website_golden.png"

	if err := builder.BuildWebsitePdf(pdfPath); err != nil {
		t.Fatalf("Failed to build website PDF: %v", err)
	}

	generateAndCompareGolden(t, pdfPath, goldenPath)
}

func TestBuildEnrollmentPdf_Golden(t *testing.T) {
	builder := createMockBuilder()
	pdfPath := filepath.Join(t.TempDir(), "enrollment.pdf")
	goldenPath := "testdata/enrollment_golden.png"

	if err := builder.BuildEnrollmentPdf(pdfPath); err != nil {
		t.Fatalf("Failed to build enrollment PDF: %v", err)
	}

	generateAndCompareGolden(t, pdfPath, goldenPath)
}

func TestBuildEmailPdf_Golden(t *testing.T) {
	builder := createMockBuilder()
	pdfPath := filepath.Join(t.TempDir(), "email.pdf")
	goldenPath := "testdata/email_golden.png"

	if err := builder.BuildEmailPdf(pdfPath); err != nil {
		t.Fatalf("Failed to build email PDF: %v", err)
	}

	generateAndCompareGolden(t, pdfPath, goldenPath)
}

func TestBuildTeacherPdf_Golden(t *testing.T) {
	builder := createMockBuilder()
	pdfPath := filepath.Join(t.TempDir(), "teacher.pdf")
	goldenPath := "testdata/teacher_golden.png"

	if err := builder.BuildTeacherPdf(pdfPath); err != nil {
		t.Fatalf("Failed to build teacher PDF: %v", err)
	}

	generateAndCompareGolden(t, pdfPath, goldenPath)
}
