package pdfbuilder

import (
	"fmt"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/orientation"
	"github.com/johnfercher/maroto/v2/pkg/core"
)

type RegistrationRenderer interface {
	Render(courses []*CourseInfo) ([]core.Row, error)
	header(*SelectionInfo, string) []core.Row
}

type Builder struct {
	Period    string
	Selection *SelectionInfo
	Courses   []*CourseInfo
}

func (b *Builder) BuildWebsitePdf(file string) error {
	return b.buildToFile(&WebsiteRenderer{}, orientation.Vertical, file)
}

func (b *Builder) BuildEnrollmentPdf(file string) error {
	return b.buildToFile(&EnrollmentRenderer{}, orientation.Horizontal, file)
}

func (b *Builder) BuildEmailPdf(file string) error {
	return b.buildToFile(&EmailRenderer{}, orientation.Horizontal, file)
}

func (b *Builder) BuildTeacherPdf(file string) error {
	return b.buildToFile(&TeacherRenderer{}, orientation.Horizontal, file)
}

func (b *Builder) buildToFile(renderer RegistrationRenderer, orientation orientation.Type, file string) error {
	pdf, err := b.build(renderer, orientation)
	if err != nil {
		return err
	}

	if err := pdf.Save(file); err != nil {
		return fmt.Errorf("maroto save: %w", err)
	}

	return nil
}

func (b *Builder) build(renderer RegistrationRenderer, orientation orientation.Type) (core.Document, error) {
	config := config.NewBuilder().WithOrientation(orientation)
	builder := maroto.New(config.Build())

	err := builder.RegisterHeader(renderer.header(b.Selection, b.Period)...)
	if err != nil {
		return nil, fmt.Errorf("register header: %w", err)
	}

	rows, err := renderer.Render(b.Courses)
	if err != nil {
		return nil, fmt.Errorf("applications to website: %w", err)
	}

	builder.AddRows(rows...)

	pdf, err := builder.Generate()
	if err != nil {
		return nil, fmt.Errorf("maroto generate: %w", err)
	}

	return pdf, nil
}
