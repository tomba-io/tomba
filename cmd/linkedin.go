package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tomba-io/go/tomba"

	"github.com/tomba-io/tomba/pkg/output"
	"github.com/tomba-io/tomba/pkg/start"
	"github.com/tomba-io/tomba/pkg/util"
	_url "github.com/tomba-io/tomba/pkg/validation/url"
)

// linkedinCmd represents the linkedin command
// see https://docs.tomba.io/api/finder#linkedin-finder
var linkedinCmd = &cobra.Command{
	Use:     "linkedin",
	Aliases: []string{"l"},
	Short:   "Instantly discover the email addresses of Linkedin URLs.",
	Long:    Long,
	Run:     linkedinRun,
	Example: linkedinExample,
}

// linkedinRun the actual work linkedin
func linkedinRun(cmd *cobra.Command, args []string) {
	fmt.Println(Long)
	init := start.New(conn)
	url := init.Target

	if !_url.IsValidURL(url) {
		fmt.Println(util.ErrorIcon(), util.Red(start.ErrArgumentsURL.Error()))
		return
	}
	params := tomba.Params{"url": url}
	if init.EnrichMobile {
		params["enrich_mobile"] = true
	}
	result, err := init.Tomba.LinkedinFinder(params)
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
	fmt.Println(util.WarningIcon(), util.Yellow("Why doesn't the Linkedin return any result? https://help.tomba.io/en/questions/reasons-why-linkedin-don-t-find-any-emails"))
}
