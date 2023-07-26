package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tomba-io/tomba/pkg/output"
	"github.com/tomba-io/tomba/pkg/start"
	"github.com/tomba-io/tomba/pkg/util"
	_email "github.com/tomba-io/tomba/pkg/validation/email"
)

// enrichCmd represents the enrich command
// see https://developer.tomba.io/#author-finder
var enrichCmd = &cobra.Command{
	Use:     "enrich",
	Aliases: []string{"e"},
	Short:   "Locate and include data in your emails.",
	Long:    Long,
	Run:     enrichRun,
	Example: enrichExample,
}

// enrichRun the actual work enrich
func enrichRun(cmd *cobra.Command, args []string) {
	fmt.Println(Long)
	init := start.New(conn)
	email := init.Target

	if !_email.IsValidEmail(email) {
		fmt.Println(util.ErrorIcon(), util.Red(start.ErrArgumentEmail.Error()))
		return
	}
	result, err := init.Tomba.Enrichment(email)
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(start.ErrErrInvalidLogin.Error()))
		return
	}
	if result.Data.Email != "" {
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
	fmt.Println(util.WarningIcon(), util.Yellow("Why doesn't the enrichment return any result? https://help.tomba.io/en/questions/reasons-why-enrichment-don-t-find-any-emails"))
}
