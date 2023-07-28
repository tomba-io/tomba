package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tomba-io/tomba/pkg/output"
	"github.com/tomba-io/tomba/pkg/start"
	"github.com/tomba-io/tomba/pkg/util"
)

// usageCmd represents the usage command
var usageCmd = &cobra.Command{
	Use:   "usage",
	Short: "Check your monthly requests.",
	Long:  Long,
	Run:   usageRun,
}

// usageRun the actual work usage
func usageRun(cmd *cobra.Command, args []string) {
	fmt.Println(Long)
	init := start.New(conn)
	result, err := init.Tomba.Usage()
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
