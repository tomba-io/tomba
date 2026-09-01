package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/tomba-io/tomba/pkg/output"
	"github.com/tomba-io/tomba/pkg/start"
	"github.com/tomba-io/tomba/pkg/util"
)

// autocompleteCmd represents the autocomplete command
var autocompleteCmd = &cobra.Command{
	Use:     "autocomplete",
	Short:   "Get domain suggestions from a partial query.",
	Long:    Long,
	Example: autocompleteExample,
	Run:     autocompleteRun,
}

// autocompleteRun the actual work autocomplete
func autocompleteRun(cmd *cobra.Command, args []string) {
	init := start.New(conn)
	query := init.Target

	if query == "" {
		fmt.Println(util.ErrorIcon(), util.Red("query is required"))
		return
	}
	result, err := init.AutoComplete(query)
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(err.Error()))
		return
	}
	raw, _ := result.Marshal()
	output.Render(string(raw), init.JSON, init.YAML, init.Output, "autocomplete")
}
