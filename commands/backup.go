package commands

import (
	"fmt"

	"github.com/baldugus/sisu/database"
)

type BackupCommand struct {
	FilePath string
}

func (cmd *BackupCommand) Execute(db *database.Database) error {
	if err := db.Backup(cmd.FilePath); err != nil {
		return fmt.Errorf("backup: %w", err)
	}

	return nil
}

type RestoreCommand struct {
	FilePath string
}

func (cmd *RestoreCommand) Execute(db *database.Database) error {
	if err := db.Restore(cmd.FilePath); err != nil {
		return fmt.Errorf("restore: %w", err)
	}

	return nil
}
