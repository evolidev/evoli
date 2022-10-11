package command

import (
	"evoli.dev/framework/console"
	"evoli.dev/framework/console/color"
	"evoli.dev/framework/use"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
	"runtime"
)

func About() *console.Command {
	return &console.Command{
		Definition:  "about",
		Description: "Get information about the application",
		Execution:   aboutRun,
	}
}

func aboutRun(cmd *console.ParsedCommand) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Environment", ""})
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
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

	info := map[string]any{
		"Evoli version": "0.0.1",
		"Go version":    fmt.Sprintf("%s %s", runtime.Version(), runtime.GOOS),
		"Environment":   color.Text(120, "development"),
		"Debug":         "false",
		"Host":          "http://localhost:8080",
		"Root":          use.BasePath(),
	}

	for k, v := range info {
		table.Append([]string{
			color.Text(140, k+":"),
			color.Text(245, v),
		})
	}

	fmt.Println()
	table.Render()
	fmt.Println()
}
