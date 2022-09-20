package console

import (
	"fmt"
	"github.com/spf13/cast"
	"regexp"
	"strings"
)

type ParsedCommand struct {
	arguments  map[string]any
	options    map[string]any
	command    string
	name       string
	subCommand string
	prefix     string
}

type Option struct {
	fmt.Stringer
	Value any
}

func (o *Option) Bool() bool {
	return cast.ToBool(o.Value)
}

func (o *Option) Integer() int {
	return cast.ToInt(o.Value)
}

func (o *Option) String() string {
	return cast.ToString(o.Value)
}

func (p *ParsedCommand) GetArgument(name string) any {
	argumentValue := p.arguments[name]
	return argumentValue
}

func (p *ParsedCommand) GetOption(name string) *Option {
	optionValue := p.options[name]

	if optionValue == nil {
		return nil
	}

	return &Option{Value: optionValue}
}

func (p *ParsedCommand) GetName() string {
	return p.name
}

func (p *ParsedCommand) GetPrefix() string {
	return p.prefix
}

func (p *ParsedCommand) GetSubCommand() string {
	return p.subCommand
}

func (p *ParsedCommand) GetOptionWithDefault(name string, defaultValue any) *Option {
	optionValue := p.options[name]
	if optionValue == nil || optionValue == "" {
		return &Option{Value: defaultValue}
	}
	return &Option{Value: optionValue}
}

var parseRegex = "[\\/-]{0,2}?((\\w+)(?:[=:](\"[^\"]+\"|[^\\s\"]+))?)(?:\\s+|$)"

func Parse(definition string, command string) *ParsedCommand {
	arguments := make(map[string]any)
	options := make(map[string]any)

	var argumentsMap []string

	// parse definition
	argumentsMap = parseDefinition(definition, options, arguments, argumentsMap)

	items := parseCommand(command, options, argumentsMap, arguments)

	name := items[0]
	// split name into prefix and subcommand
	nameParts := strings.Split(name, ":")
	prefix := nameParts[0]
	subCommand := ""
	if len(nameParts) > 1 {
		subCommand = nameParts[1]
	}

	return &ParsedCommand{
		arguments:  arguments,
		options:    options,
		command:    command,
		name:       name,
		subCommand: subCommand,
		prefix:     prefix,
	}
}

func parseDefinition(definition string, options map[string]any, arguments map[string]any, argumentsMap []string) []string {
	definitionItems := strings.Split(definition, " ")
	for i := range definitionItems {
		definitionItems[i] = strings.TrimSpace(definitionItems[i])

		// parse definition item
		definitionItem := definitionItems[i]
		// remove curly bracket at the beginning and end of definition item
		definitionItem = strings.Trim(definitionItem, "{}?")

		// split definition item into name and Value
		definitionItemParts := strings.Split(definitionItem, "=")
		definitionItemName := definitionItemParts[0]
		definitionItemValue := ""
		if len(definitionItemParts) > 1 {
			definitionItemValue = definitionItemParts[1]
		}

		if strings.HasPrefix(definitionItemName, "--") {
			optionName := strings.TrimPrefix(definitionItemName, "--")

			optionNameParts := strings.Split(optionName, "|")
			for index := range optionNameParts {
				optionName = strings.TrimSpace(optionNameParts[index])
				options[optionName] = definitionItemValue
			}
		} else {
			arguments[definitionItemName] = definitionItemValue
			argumentsMap = append(argumentsMap, definitionItemName)
		}
	}
	return argumentsMap
}

func parseCommand(command string, options map[string]any, argumentsMap []string, arguments map[string]any) []string {
	// extract all arguments and options
	r, _ := regexp.Compile(parseRegex)
	items := r.FindAllString(command, -1)
	for i := range items {
		items[i] = strings.TrimSpace(items[i])
	}

	for index, item := range items {
		if strings.HasPrefix(item, "--") {
			optionName, optionValue := extractField(item, "--")
			options[optionName] = optionValue
		} else if strings.HasPrefix(item, "-") {
			optionName, optionValue := extractField(item, "-")
			options[optionName] = optionValue
		} else {
			if index > 0 && index < len(argumentsMap) {
				arguments[argumentsMap[index]] = item
			}
		}
	}
	return items
}

func extractField(item string, prefix string) (string, any) {
	option := strings.TrimPrefix(item, prefix)
	// extract option name and Value
	parts := strings.Split(option, "=")
	optionName := parts[0]
	var optionValue any
	if len(parts) > 1 && parts[1] != "" {
		optionValue = parts[1]
	} else {
		optionValue = true
	}
	return optionName, optionValue
}
