package sql

import (
	"embed"
	"fmt"
)

//go:embed *.sql
var queries embed.FS

func Query(name string) string {
	query, err := queries.ReadFile(fmt.Sprintf("%s.sql", name))
	if err != nil {
		panic(err)
	}

	return string(query)
}
