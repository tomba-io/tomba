package output

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

// DisplayTable a table as output.
func DisplayTable(data [][]string, header []string) {
	table := tablewriter.NewWriter(os.Stdout)

	// Set the table headers
	table.SetHeader(header)
	table.SetHeaderColor(tablewriter.Colors{tablewriter.Bold, tablewriter.BgGreenColor})

	table.SetColumnColor(tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiBlackColor})

	// Remove the header row from the data
	data = data[1:]

	// Add data rows to the table
	table.AppendBulk(data)

	// Set table formatting options
	table.SetCenterSeparator("|")
	table.SetAutoWrapText(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)

	// Render the table
	table.Render()
}
