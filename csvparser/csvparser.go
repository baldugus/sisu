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

func ParseBytes(content []byte, path string) (*ParsedCsv, error) {
	r := bytes.NewReader(content)
	parsedCsv, err := Parse(r, path)
	if err != nil {
		return nil, fmt.Errorf("parsing csv: %w", err)
	}

	return parsedCsv, nil
}

func ParseFile(path string) (p *ParsedCsv, err error) {
	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, &ErrFileNotFound{Path: path, Err: err}
		}

		if errors.Is(err, os.ErrPermission) {
			return nil, &ErrPermissionDenied{Path: path, Err: err}
		}

		return nil, &ErrFileOpen{Path: path, Err: err}
	}

	defer func() {
		if cerr := file.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

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

	candidates := []*csvCandidate{}

	if err := gocsv.Unmarshal(sanitizedReader, &candidates); err != nil {
		if errors.Is(err, gocsv.ErrEmptyCSVFile) {
			return nil, &ErrFileEmpty{Err: err}
		}

		return nil, &ErrFileParse{Err: err}
	}

	return &ParsedCsv{candidates: candidates, name: name}, nil
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
		return nil, &ErrFileRead{Err: err}
	}

	lines := strings.Split(string(rawCsv), "\r")
	lines[0] = strings.ReplaceAll(lines[0], ",", ";")
	sanitizedReader := bytes.NewReader([]byte(strings.Join(lines, "\n")))

	return sanitizedReader, nil
}
