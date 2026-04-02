package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/tomba-io/tomba/pkg/output"
	"github.com/tomba-io/tomba/pkg/start"
	"github.com/tomba-io/tomba/pkg/util"
)

// whoamiCmd represents the whoami command
var whoamiCmd = &cobra.Command{
	Use:     "whoami",
	Aliases: []string{"w"},
	Short:   "Print current account information.",
	Long:    Long,
	Run:     whoamiRun,
	Example: whoamiExample,
}

// whoamiRun the actual work whoami
func whoamiRun(cmd *cobra.Command, args []string) {
	init := start.New(conn)
	result, err := init.Tomba.Account()
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(start.ErrErrInvalidLogin.Error()))
		return
	}
	raw, _ := result.Marshal()
	output.Render(string(raw), init.JSON, init.YAML, init.Output, "whoami")
}
