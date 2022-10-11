package command

import (
	"github.com/evolidev/evoli/framework/console"
	"github.com/evolidev/evoli/framework/use"
)

func Migrate() *console.Command {
	return &console.Command{
		Definition:  "migrate",
		Description: "Migrate the database",
		Execution:   migrateRun,
	}
}

func migrateRun(cmd *console.ParsedCommand) {
	use.Migration().Migrate(use.DB())
}
