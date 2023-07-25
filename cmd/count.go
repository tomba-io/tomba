package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tomba-io/email/pkg/output"
	"github.com/tomba-io/email/pkg/start"
	"github.com/tomba-io/email/pkg/util"
	_domain "github.com/tomba-io/email/pkg/validation/domain"
	_key "github.com/tomba-io/email/pkg/validation/key"
)

// countCmd represents the count command
// see https://developer.tomba.io/#email-count
var countCmd = &cobra.Command{
	Use:     "count",
	Aliases: []string{"c"},
	Short:   "Returns total email addresses we have for one domain.",
	Long:    Long,
	Run:     countRun,
	Example: countExample,
}

// countRun the actual work count
func countRun(cmd *cobra.Command, args []string) {
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
	url := init.Target

	if !_domain.IsValidDomain(url) {
		fmt.Println(util.ErrorIcon(), util.Red(start.ErrArgumentsDomain.Error()))
		return
	}
	result, err := init.Tomba.Count(url)
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(start.ErrErrInvalidLogin.Error()))
		return
	}
	if result.Data.Total > 0 {
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
	fmt.Println(util.WarningIcon(), util.Yellow("TombaPublicWebCrawler is searching the internet for the best leads that relate to this company, but we don't have any for it yet. That said, we're working on it"))
}
