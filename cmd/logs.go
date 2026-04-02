package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/tomba-io/tomba/pkg/output"
	"github.com/tomba-io/tomba/pkg/start"
	"github.com/tomba-io/tomba/pkg/util"
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Check your last 1,000 requests you made during the last 3 months.",
	Long:  Long,
	Run:   logsRun,
}

// logsRun the actual work logs
func logsRun(cmd *cobra.Command, args []string) {
	init := start.New(conn)

	result, err := init.Tomba.Logs()
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(err.Error()))
		return
	}
	raw, _ := result.Marshal()
	output.Render(string(raw), init.JSON, init.YAML, init.Output, "logs")
}
