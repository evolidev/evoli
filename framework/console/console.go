package console

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/evolidev/evoli/framework/console/color"
	"github.com/olekukonko/tablewriter"
)

type Command struct {
	Definition  string
	Description string
	Execution   func(c *ParsedCommand)
}

func (cmd *Command) GetName() string {
	parts := strings.Split(cmd.Definition, " ")
	return parts[0]
}

func (cmd *Command) GetCommand() string {
	name := cmd.GetName()

	parts := strings.Split(name, ":")
	if len(parts) > 1 {
		name = parts[len(parts)-1]
	}

	return name
}

func (cmd *Command) GetDescription() string {
	return cmd.Description
}

func (cmd *Command) Run(c *ParsedCommand) {
	//return cmd.Name
}

type CommandGroup struct {
	Name        string
	Description string
	Prefix      string
	Commands    []*Command
}

type Console struct {
	Commands map[string]*Command
}

func (c *Console) Run() {
	args := os.Args[1:]

	if len(args) > 0 {
		command := args[0]
		if cmd, ok := c.Commands[command]; ok {
			parsed := Parse(cmd.Definition, strings.Join(args, " "))
			cmd.Execution(parsed)
			return
		} else {
			fmt.Println(color.Text(140, "Command not found"))
		}
	}

	Render(c.Commands)
}

func (c *Console) Add(command *Command) {
	c.Commands[command.GetName()] = command
}

func (c *Console) AddCommand(name string, description string, execution func(c *ParsedCommand)) *Command {
	command := &Command{name, description, execution}
	c.Commands[command.GetName()] = command

	return command
}

func New() *Console {
	return &Console{
		Commands: make(map[string]*Command),
	}
}

func groupCommands(commands map[string]*Command) []CommandGroup {
	groups := make(map[string][]*Command)

	var keys []string
	for _, cmd := range commands {
		commandParts := strings.Split(cmd.GetName(), ":")
		prefix := ""
		if len(commandParts) > 1 {
			prefix = commandParts[0]
		}

		if _, ok := groups[prefix]; !ok {
			keys = append(keys, prefix)
		}

		groups[prefix] = append(groups[prefix], cmd)
	}

	sort.Strings(keys)

	var groupedCommands []CommandGroup
	for _, key := range keys {
		groupedCommands = append(groupedCommands, CommandGroup{
			Name:        key,
			Description: "",
			Prefix:      key,
			Commands:    groups[key],
		})
	}

	return groupedCommands
}

func Render(commands map[string]*Command) {
	table := setupTable()

	addCommandsToTable(commands, table)

	fmt.Println(fmt.Sprintf(
		"Evoli Console %s", color.Text(169, "0.0.1"),
	))
	fmt.Println()
	table.Render()
}

func addCommandsToTable(commands map[string]*Command, table *tablewriter.Table) {
	groupedCommands := groupCommands(commands)
	for _, group := range groupedCommands {
		prefix := ""
		if group.Name != "" {
			table.Rich([]string{group.Name, group.Description}, []tablewriter.Colors{
				{tablewriter.FgHiGreenColor},
				{},
			})

			prefix = color.Text(140, group.Prefix+":")
		}

		for _, cmd := range group.Commands {
			table.Append([]string{
				prefix + color.Text(169, cmd.GetCommand()),
				color.Text(245, cmd.Description),
			})
		}

		table.Append([]string{""})
	}
}

func setupTable() *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Available Commands", ""})
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)

	table.SetColumnColor(
		tablewriter.Colors{tablewriter.FgHiMagentaColor},
		tablewriter.Colors{tablewriter.FgHiBlackColor},
	)

	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.FgHiWhiteColor},
		tablewriter.Colors{tablewriter.FgHiBlackColor},
	)
	return table
}
