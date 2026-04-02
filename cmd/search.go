package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tomba-io/go/tomba"

	"github.com/tomba-io/tomba/pkg/output"
	"github.com/tomba-io/tomba/pkg/start"
	"github.com/tomba-io/tomba/pkg/util"
	_domain "github.com/tomba-io/tomba/pkg/validation/domain"
)

// searchCmd represents the search command
// see https://docs.tomba.io/api/finder#domain-search
var searchCmd = &cobra.Command{
	Use:     "search",
	Aliases: []string{"s"},
	Short:   "Instantly locate email addresses from any company name or website.",
	Long:    Long,
	Run:     searchRun,
	Example: searchExample,
}

// searchRun the actual work search
func searchRun(cmd *cobra.Command, args []string) {

	init := start.New(conn)
	domain := init.Target

	if !_domain.IsValidDomain(domain) {
		fmt.Println(util.ErrorIcon(), util.Red(start.ErrArgumentsDomain.Error()))
		return

	}
	params := tomba.Params{"domain": domain}
	if init.Page != 0 {
		params["page"] = fmt.Sprint(init.Page)
	}
	switch init.Limit {
	case 10, 20, 50:
		params["limit"] = fmt.Sprint(init.Limit)
	default:
		fmt.Println(util.ErrorIcon(), util.Red(start.ErrArgumentsDomainLimit.Error()), init.Page)
		return
	}
	result, err := init.Tomba.DomainSearch(params)
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(err.Error()))
		return
	}
	if result.Meta.Total > 0 {
		raw, _ := result.Marshal()
		output.Render(string(raw), init.JSON, init.YAML, init.Output, "search")
		return
	}
	fmt.Println(util.WarningIcon(), util.Yellow("TombaPublicWebCrawler is searching the internet for the best leads that relate to this company, but we don't have any for it yet. That said, we're working on it"))
}
