package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/tomba-io/tomba/pkg/output"
	"github.com/tomba-io/tomba/pkg/start"
	"github.com/tomba-io/tomba/pkg/util"
	_domain "github.com/tomba-io/tomba/pkg/validation/domain"
)

// similarCmd represents the similar command
// see https://docs.tomba.io/api/~endpoints#similar
var similarCmd = &cobra.Command{
	Use:     "similar",
	Aliases: []string{"sim"},
	Short:   "Retrieve domains similar to a specific domain.",
	Long:    Long,
	Run:     similarRun,
	Example: similarExample,
}

// similarRun the actual work similar
func similarRun(cmd *cobra.Command, args []string) {

	init := start.New(conn)
	domain := init.Target

	if !_domain.IsValidDomain(domain) {
		fmt.Println(util.ErrorIcon(), util.Red(start.ErrArgumentsDomain.Error()))
		return
	}

	result, err := init.SimilarDomains(domain)
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(err.Error()))
		return
	}

	raw, _ := result.Marshal()
	output.Render(string(raw), init.JSON, init.YAML, init.Output, "similar")
}
