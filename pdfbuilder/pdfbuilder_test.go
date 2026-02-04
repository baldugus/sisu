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
		ID:          1,
		Name:        "Seleção de Teste",
		Kind:        1, // ApprovedSelection
		Date:        "2025-01-01T00:00:00Z",
		Institution: "FAETERJ-Rio",
		Course:      "Análise e Desenvolvimento de Sistemas",
	}
}

func createMockApplications() types.Applications {
	return types.Applications{
		{
			ID:             1,
			Status:         "APPROVED",
			EnrollmentID:   "123456789",
			Option:         1,
			CompositeScore: 850.5,
			Ranking:        1,
			Applicant: types.Applicant{
				ID:        1,
				Name:      "João Silva",
				Email:     "joao.silva@exemplo.com",
				CPF:       "12345678901",
				BirthDate: "1990-01-01",
			},
		},
		{
			ID:             2,
			Status:         "APPROVED",
			EnrollmentID:   "987654321",
			Option:         1,
			CompositeScore: 820.3,
			Ranking:        2,
			Applicant: types.Applicant{
				ID:        2,
				Name:      "Maria Santos",
				Email:     "maria.santos@exemplo.com",
				CPF:       "98765432109",
				BirthDate: "1992-05-15",
			},
		},
	}
}

func createMockClassInfo() []*ClassInfo {
	apps := createMockApplications()
	return []*ClassInfo{
		NewClassInfo("Ampla concorrência", apps),
	}
}

func createMockSelectionInfo() *SelectionInfo {
	sel := createMockSelection()
	return NewSelectionInfo(sel, "2025", "1", 0)
}

func createMockBuilder() *Builder {
	return &Builder{
		Period:    "Manhã",
		Selection: createMockSelectionInfo(),
		Classes:   createMockClassInfo(),
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
