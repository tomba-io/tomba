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
	result, err := init.LinkedinFinder(params)
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(err.Error()))
		return
	}
	finderData := start.GetFinderData(result.Data)
	if finderData != nil && finderData.Email != "" {
		raw, _ := result.Marshal()
		output.Render(string(raw), init.JSON, init.YAML, init.Output, "linkedin")
		return
	}
	fmt.Println(util.WarningIcon(), util.Yellow("Why doesn't the Linkedin return any result? https://help.tomba.io/en/questions/reasons-why-linkedin-don-t-find-any-emails"))
}
