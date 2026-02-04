package main

import (
	"fmt"
	"os"

	"changeme/types"

	"github.com/gocarina/gocsv"
)

func ApplicationsToCSV(applications []*types.Application, filename string) error {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	defer file.Close()
	if err != nil {
		return fmt.Errorf("os open file: %w", err)
	}

	if err := gocsv.Marshal(applications, file); err != nil {
		return fmt.Errorf("csv marshal: %w", err)
	}

	return nil
}
