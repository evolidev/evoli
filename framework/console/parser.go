package console

import (
	"regexp"
	"strings"
)

type ParsedCommand struct {
	arguments  map[string]interface{}
	options    map[string]interface{}
	command    string
	name       string
	subCommand string
	prefix     string
}

func (p *ParsedCommand) GetArgument(name string) interface{} {
	argumentValue := p.arguments[name]
	return argumentValue
}

func (p *ParsedCommand) GetOption(name string) interface{} {
	optionValue := p.options[name]
	return optionValue
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

func (p *ParsedCommand) GetOptionWithDefault(s string, defaultValue interface{}) interface{} {
	option := p.GetOption(s)
	if option == nil {
		return defaultValue
	}
	return option
}

var parseRegex = "[\\/-]{0,2}?((\\w+)(?:[=:](\"[^\"]+\"|[^\\s\"]+))?)(?:\\s+|$)"

func Parse(definition string, command string) *ParsedCommand {

	arguments := make(map[string]interface{})
	options := make(map[string]interface{})

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

func parseDefinition(definition string, options map[string]interface{}, arguments map[string]interface{}, argumentsMap []string) []string {
	definitionItems := strings.Split(definition, " ")
	for i := range definitionItems {
		definitionItems[i] = strings.TrimSpace(definitionItems[i])

		// parse definition item
		definitionItem := definitionItems[i]
		// remove curly bracket at the beginning and end of definition item
		definitionItem = strings.Trim(definitionItem, "{}?")

		// split definition item into name and value
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

func parseCommand(command string, options map[string]interface{}, argumentsMap []string, arguments map[string]interface{}) []string {
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

func extractField(item string, prefix string) (string, interface{}) {
	option := strings.TrimPrefix(item, prefix)
	// extract option name and value
	parts := strings.Split(option, "=")
	optionName := parts[0]
	var optionValue interface{}
	if len(parts) > 1 && parts[1] != "" {
		optionValue = parts[1]
	} else {
		optionValue = true
	}
	return optionName, optionValue
}
