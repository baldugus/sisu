package csvparser

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/dimchansky/utfbom"
	"github.com/gocarina/gocsv"
)

func ParseFile(path string) (*ParsedCsv, error) {
	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, &FileNotFoundError{Path: path, Err: err}
		}

		if errors.Is(err, os.ErrPermission) {
			return nil, &PermissionDeniedError{Path: path, Err: err}
		}

		return nil, &FileError{Path: path, Err: err}
	}
	defer file.Close()

	parsedCsv, err := Parse(file, file.Name())
	if err != nil {
		return nil, fmt.Errorf("parsing csv: %w", err)
	}

	return parsedCsv, nil
}

func Parse(reader io.Reader, name string) (*ParsedCsv, error) {
	sanitizedReader, err := sanitizeCsv(reader)
	if err != nil {
		return nil, fmt.Errorf("sanitizing csv: %w", err)
	}

	// We use LazyQuotes because some fields are surrounded by quotes while
	// others don't. Same comment from sanitizeCsv() applies here, we have no clue
	// if this is like this at source.
	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		c := csv.NewReader(in)
		c.Comma = ';'
		c.LazyQuotes = true

		return c
	})

	applicants := []*csvApplicant{}

	if err := gocsv.Unmarshal(sanitizedReader, &applicants); err != nil {
		if errors.Is(err, gocsv.ErrEmptyCSVFile) {
			return nil, &EmptyError{Err: err}
		}

		return nil, &ParseError{Err: err}
	}

	return &ParsedCsv{applicants: applicants, name: name}, nil
}

/*
SISU's CSV file has peculiarities that require handling before parsing:

  - First character in some cases is a UTF-8 BOM, hence the SkipOnly call;

  - Line breaks are encoded using CR instead of LF (go expects LF);

  - The header is separated by commas instead of semicolons (like the values).

Note that we're not sure if the file comes like this or if all of this is made
by someones excel before the CSV gets here, but we have no direct access to
the source file, so we have to deal with it.
*/
func sanitizeCsv(reader io.Reader) (io.Reader, error) {
	rawCsv, err := io.ReadAll(utfbom.SkipOnly(reader))
	if err != nil {
		return nil, &ReadError{Err: err}
	}

	lines := strings.Split(string(rawCsv), "\r")
	lines[0] = strings.ReplaceAll(lines[0], ",", ";")
	sanitizedReader := bytes.NewReader([]byte(strings.Join(lines, "\n")))

	return sanitizedReader, nil
}
