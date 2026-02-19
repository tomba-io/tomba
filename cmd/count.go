package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/tomba-io/tomba/pkg/output"
	"github.com/tomba-io/tomba/pkg/start"
	"github.com/tomba-io/tomba/pkg/util"
	_domain "github.com/tomba-io/tomba/pkg/validation/domain"
)

// countCmd represents the count command
// see https://docs.tomba.io/api/finder#email-count
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
	domain := init.Target

	if !_domain.IsValidDomain(domain) {
		fmt.Println(util.ErrorIcon(), util.Red(start.ErrArgumentsDomain.Error()))
		return
	}
	result, err := init.Tomba.Count(domain)
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(start.ErrErrInvalidLogin.Error()))
		return
	}
	if result.Data.Total > 0 {
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
	fmt.Println(util.WarningIcon(), util.Yellow("TombaPublicWebCrawler is searching the internet for the best leads that relate to this company, but we don't have any for it yet. That said, we're working on it"))
}
