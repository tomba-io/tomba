package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/tomba-io/tomba/pkg/output"
	"github.com/tomba-io/tomba/pkg/start"
	"github.com/tomba-io/tomba/pkg/util"
	_domain "github.com/tomba-io/tomba/pkg/validation/domain"
)

// technologyCmd represents the technology command
// see https://docs.tomba.io/api/domain#technology
var technologyCmd = &cobra.Command{
	Use:     "technology",
	Aliases: []string{"tech"},
	Short:   "Discover technologies detected for a domain.",
	Long:    Long,
	Run:     technologyRun,
	Example: technologyExample,
}

// technologyRun the actual work technology
func technologyRun(cmd *cobra.Command, args []string) {

	init := start.New(conn)
	domain := init.Target

	if !_domain.IsValidDomain(domain) {
		fmt.Println(util.ErrorIcon(), util.Red(start.ErrArgumentsDomain.Error()))
		return
	}

	result, err := init.TechnologyCheck(domain)
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(err.Error()))
		return
	}

	raw, _ := result.Marshal()
	output.Render(string(raw), init.JSON, init.YAML, init.Output, "technology")
}
