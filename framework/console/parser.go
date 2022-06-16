package console

import (
	"github.com/evolidev/evoli/framework/use"
	"regexp"
	"strings"
)

type ParsedCommand struct {
}

func (p *ParsedCommand) GetArgument(name string) interface{} {
	return "foo"
}

func (p *ParsedCommand) GetOption(name string) interface{} {
	return "omer"
}

var parseRegex = "[\\/-]{0,2}?((\\w+)(?:[=:](\"[^\"]+\"|[^\\s\"]+))?)(?:\\s+|$)"

func Parse(definition string, command string) *ParsedCommand {
	r, _ := regexp.Compile(parseRegex)
	items := r.FindAllString(command, -1)
	for i := range items {
		items[i] = strings.TrimSpace(items[i])
	}

	use.D(items)
	return &ParsedCommand{}
}
