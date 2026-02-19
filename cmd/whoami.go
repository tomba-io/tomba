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
	fmt.Println(Long)
	init := start.New(conn)
	result, err := init.Tomba.Account()
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(start.ErrErrInvalidLogin.Error()))
		return
	}
	raw, _ := result.Marshal()
	if init.JSON {
		json, _ := output.DisplayJSON(string(raw))
		fmt.Println(json)
	}
	if init.YAML {
		yaml, _ := output.DisplayYAML(string(raw))
		fmt.Println(yaml)
	}
	if init.Output != "" {
		err := output.CreateOutput(init.Output, string(raw))
		if err != nil {
			fmt.Println("Error creating file:", err)
		}
	}
}
