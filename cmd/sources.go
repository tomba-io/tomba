package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/tomba-io/tomba/pkg/output"
	"github.com/tomba-io/tomba/pkg/start"
	"github.com/tomba-io/tomba/pkg/util"
	_email "github.com/tomba-io/tomba/pkg/validation/email"
)

// sourcesCmd represents the sources command
var sourcesCmd = &cobra.Command{
	Use:     "sources",
	Aliases: []string{"s"},
	Short:   "Find email address source somewhere on the web",
	Long:    Long,
	Example: sourcesExample,
	Run:     sourcesRun,
}

// sourcesRun the actual work sources
func sourcesRun(cmd *cobra.Command, args []string) {
	fmt.Println(Long)
	init := start.New(conn)
	email := init.Target

	if !_email.IsValidEmail(email) {
		fmt.Println(util.ErrorIcon(), util.Red(start.ErrArgumentEmail.Error()))
		return
	}

	result, err := init.Tomba.Sources(email)
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(start.ErrErrInvalidLogin.Error()))
		return
	}
	if len(result.Sources) > 0 {

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
		return
	}
	fmt.Println(util.WarningIcon(), util.Yellow("We haven't found this email on the web."))
}
