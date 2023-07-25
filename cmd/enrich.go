package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tomba-io/email/pkg/output"
	"github.com/tomba-io/email/pkg/start"
	"github.com/tomba-io/email/pkg/util"
	_email "github.com/tomba-io/email/pkg/validation/email"
	_key "github.com/tomba-io/email/pkg/validation/key"
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
	if init.Key == "" || init.Secret == "" {
		fmt.Println(util.WarningIcon(), util.Yellow(start.ErrErrInvalidNoLogin.Error()))
		return
	}
	if !_key.IsValidAPI(init.Key) && !_key.IsValidAPI(init.Secret) {
		fmt.Println(util.WarningIcon(), util.Yellow(start.ErrErrInvalidLogin.Error()))
		return
	}

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
		if init.JSON {
			raw, _ := result.Marshal()
			json, _ := output.DisplayJSON(string(raw))
			fmt.Println(json)
			return
		}
		if init.YAML {
			raw, _ := result.Marshal()
			yaml, _ := output.DisplayYAML(string(raw))
			fmt.Println(yaml)
			return
		}
		return
	}
	fmt.Println(util.WarningIcon(), util.Yellow("Why doesn't the enrichment return any result? https://help.tomba.io/en/questions/reasons-why-enrichment-don-t-find-any-emails"))
}
