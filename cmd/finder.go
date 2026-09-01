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

// finderCmd represents the finder command
// see https://docs.tomba.io/api/finder#email-finder
var finderCmd = &cobra.Command{
	Use:     "finder",
	Aliases: []string{"s"},
	Short:   "Retrieves the most likely email address from a domain name, a first name and a last name.",
	Long:    Long,
	Run:     finderRun,
	Example: finderExample,
}

// finderRun the actual work finder
func finderRun(cmd *cobra.Command, args []string) {

	init := start.New(conn)
	domain := init.Target

	if !_domain.IsValidDomain(domain) {
		fmt.Println(util.ErrorIcon(), util.Red(start.ErrArgumentsDomain.Error()))
		return
	}

	params := tomba.Params{"domain": domain}

	switch {
	case init.FirstName != "" && init.LastName != "":
		params["first_name"] = init.FirstName
		params["last_name"] = init.LastName
	case init.FullName != "":
		params["full_name"] = init.FullName
	default:
		fmt.Println(util.ErrorIcon(), util.Red(start.ErrArgumentsFinder.Error()))
		return
	}

	if init.EnrichMobile {
		params["enrich_mobile"] = true
	}

	result, err := init.EmailFinder(params)
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(err.Error()))
		return
	}
	finderData := start.GetFinderData(result.Data)
	if finderData != nil && finderData.Email != "" {
		raw, _ := result.Marshal()
		output.Render(string(raw), init.JSON, init.YAML, init.Output, "finder")
		return
	}
	fmt.Println(util.WarningIcon(), util.Yellow("Why doesn't the Email Finder return any result? https://help.tomba.io/en/questions/why-doesn-t-the-email-finder-return-any-result"))
}
