package main

import (
	"fmt"
	"os"

	"github.com/baldugus/sisu/types"

	"github.com/gocarina/gocsv"
)

func ApplicationsToCSV(applications []*types.Application, filename string) error {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return fmt.Errorf("os open file: %w", err)
	}
	defer file.Close()

	if err := gocsv.Marshal(applications, file); err != nil {
		return fmt.Errorf("csv marshal: %w", err)
	}

	return nil
}
