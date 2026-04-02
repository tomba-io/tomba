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

	init := start.New(conn)
	domain := init.Target

	if !_domain.IsValidDomain(domain) {
		fmt.Println(util.ErrorIcon(), util.Red(start.ErrArgumentsDomain.Error()))
		return
	}
	result, err := init.Tomba.Count(domain)
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(err.Error()))
		return
	}
	if result.Data.Total > 0 {
		raw, _ := result.Marshal()
		output.Render(string(raw), init.JSON, init.YAML, init.Output, "count")
		return
	}
	fmt.Println(util.WarningIcon(), util.Yellow("TombaPublicWebCrawler is searching the internet for the best leads that relate to this company, but we don't have any for it yet. That said, we're working on it"))
}
