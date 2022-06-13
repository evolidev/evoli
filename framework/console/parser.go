package console

import (
	"github.com/evolidev/evoli/framework/use"
	"regexp"
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
	use.D(r.FindAllString(command, -1))
	return &ParsedCommand{}
}
