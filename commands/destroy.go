package commands

import (
	"fmt"

	"github.com/baldugus/sisu/database"
)

type DestroyCommand struct{}

func (cmd *DestroyCommand) Execute(db *database.Database) error {
	if err := db.Destroy(); err != nil {
		return fmt.Errorf("destroy: %w", err)
	}

	return nil
}
