package command

import (
	"evoli.dev/framework/console"
	"evoli.dev/framework/use"
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
