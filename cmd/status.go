package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/tomba-io/tomba/pkg/output"
	"github.com/tomba-io/tomba/pkg/start"
	"github.com/tomba-io/tomba/pkg/util"
	_domain "github.com/tomba-io/tomba/pkg/validation/domain"
)

// statusCmd represents the status command
// see https://docs.tomba.io/api/~endpoints#domain-status
var statusCmd = &cobra.Command{
	Use:     "status",
	Aliases: []string{"t"},
	Short:   "Returns domain status if is webmail or disposable.",
	Long:    Long,
	Run:     statusRun,
	Example: statusExample,
}

// statusRun the actual work status
func statusRun(cmd *cobra.Command, args []string) {
	init := start.New(conn)
	domain := init.Target

	if !_domain.IsValidDomain(domain) {
		fmt.Println(util.ErrorIcon(), util.Red(start.ErrArgumentsDomain.Error()))
		return
	}
	result, err := init.Tomba.Status(domain)
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(start.ErrErrInvalidLogin.Error()))
		return
	}
	raw, _ := result.Marshal()
	output.Render(string(raw), init.JSON, init.YAML, init.Output, "status")
}
