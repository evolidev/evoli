package console

import (
	"fmt"
	"os"

	"github.com/evolidev/evoli/framework/console/color"
	"github.com/olekukonko/tablewriter"
)

type Command struct {
	Name        string
	Description string
	Execution   string
}

type CommandGroup struct {
	Name        string
	Description string
	Prefix      string
	Commands    []Command
}

func Commands() {
	commands := []CommandGroup{
		CommandGroup{"Routes", "", "route", []Command{
			Command{"cache", "Create a route cache file for faster route registration", ""},
			Command{"clear", "Remove the route cache file", ""},
			Command{"list", "List all registered routes", ""},
		}},

		CommandGroup{"Config", "", "config", []Command{
			Command{"cache", "Create a route cache file for faster route registration", ""},
			Command{"clear", "Remove the route cache file", ""},
			Command{"list", "List all registered routes", ""},
		}},

		CommandGroup{"Make", "", "make", []Command{
			Command{"cache", "Create a route cache file for faster route registration", ""},
			Command{"clear", "Remove the route cache file", ""},
			Command{"list", "List all registered routes", ""},
		}},

		CommandGroup{"Cache", "", "cache", []Command{
			Command{"clear", "Create a route cache file for faster route registration", ""},
			Command{"forget", "Remove the route cache file", ""},
			Command{"table", "List all registered routes", ""},
		}},

		CommandGroup{"Migrate", "", "migrate", []Command{
			Command{"fresh", "Create a route cache file for faster route registration", ""},
			Command{"generate", "Remove the route cache file", ""},
			Command{"install", "List all registered routes", ""},
			Command{"reset", "List all registered routes", ""},
			Command{"rollback", "List all registered routes", ""},
			Command{"status", "List all registered routes", ""},
		}},
	}

	commandsRender(commands)
}

func commandsRender(commands []CommandGroup) {
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

	for _, group := range commands {
		table.Rich([]string{group.Name, group.Description}, []tablewriter.Colors{
			tablewriter.Colors{tablewriter.FgHiGreenColor},
			tablewriter.Colors{},
		})

		for _, cmd := range group.Commands {
			table.Append([]string{
				color.Text(140, group.Prefix+":") + color.Text(169, cmd.Name),
				color.Text(103, cmd.Description),
			})
		}

		table.Append([]string{""})
	}

	//fmt.Println()
	fmt.Println(fmt.Sprintf("Evoli Console %s", color.Text(169, "0.0.1")))
	fmt.Println()
	table.Render()
	fmt.Println()
}

func Colored() {
	data := [][]string{
		[]string{"Test1Merge", "HelloCol2 - 1", "HelloCol3 - 1", "HelloCol4 - 1"},
		[]string{"Test1Merge", "HelloCol2 - 2", "HelloCol3 - 2", "HelloCol4 - 2"},
		[]string{"Test1Merge", "HelloCol2 - 3", "HelloCol3 - 3", "HelloCol4 - 3"},
		[]string{"Test2Merge", "HelloCol2 - 4", "HelloCol3 - 4", "HelloCol4 - 4"},
		[]string{"Test2Merge", "HelloCol2 - 5", "HelloCol3 - 5", "HelloCol4 - 5"},
		[]string{"Test2Merge", "HelloCol2 - 6", "HelloCol3 - 6", "HelloCol4 - 6"},
		[]string{"Test2Merge", "HelloCol2 - 7", "HelloCol3 - 7", "HelloCol4 - 7"},
		[]string{"Test3Merge", "HelloCol2 - 8", "HelloCol3 - 8", "HelloCol4 - 8"},
		[]string{"Test3Merge", "HelloCol2 - 9", "HelloCol3 - 9", "HelloCol4 - 9"},
		[]string{"Test3Merge", "HelloCol2 - 10", "HelloCol3 -10", "HelloCol4 - 10"},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Col1", "Col2", "Col3", "Col4"})
	table.SetFooter([]string{"", "", "Footer3", "Footer4"})
	table.SetBorder(false)
	table.SetColumnSeparator("")

	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.BgGreenColor},
		tablewriter.Colors{tablewriter.FgHiRedColor, tablewriter.Bold, tablewriter.BgBlackColor},
		tablewriter.Colors{tablewriter.BgRedColor, tablewriter.FgWhiteColor},
		tablewriter.Colors{tablewriter.BgCyanColor, tablewriter.FgWhiteColor},
	)

	table.SetColumnColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiRedColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgBlackColor},
	)

	table.SetFooterColor(
		tablewriter.Colors{},
		tablewriter.Colors{},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgHiRedColor},
	)

	colorData1 := []string{"TestCOLOR1Merge", "HelloCol2 - COLOR1", "HelloCol3 - COLOR1", "HelloCol4 - COLOR1"}
	colorData2 := []string{"TestCOLOR2Merge", "HelloCol2 - COLOR2", "HelloCol3 - COLOR2", "HelloCol4 - COLOR2"}

	for i, row := range data {
		if i == 4 {
			table.Rich(colorData1, []tablewriter.Colors{
				tablewriter.Colors{},
				tablewriter.Colors{tablewriter.Normal, tablewriter.FgCyanColor},
				tablewriter.Colors{tablewriter.Bold, tablewriter.FgWhiteColor},
				tablewriter.Colors{},
			})
			table.Rich(colorData2, []tablewriter.Colors{
				tablewriter.Colors{tablewriter.Normal, tablewriter.FgMagentaColor},
				tablewriter.Colors{},
				tablewriter.Colors{tablewriter.Bold, tablewriter.BgRedColor},
				tablewriter.Colors{tablewriter.FgHiGreenColor, tablewriter.Italic, tablewriter.BgHiCyanColor},
			})
		}
		table.Append(row)
	}

	table.SetAutoMergeCells(true)
	table.Render()
}
