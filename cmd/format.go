package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/tomba-io/tomba/pkg/output"
	"github.com/tomba-io/tomba/pkg/start"
	"github.com/tomba-io/tomba/pkg/util"
	_domain "github.com/tomba-io/tomba/pkg/validation/domain"
)

// formatCmd represents the format command
var formatCmd = &cobra.Command{
	Use:     "format",
	Short:   "Get the email format used by a domain.",
	Long:    Long,
	Example: formatExample,
	Run:     formatRun,
}

// formatRun the actual work format
func formatRun(cmd *cobra.Command, args []string) {
	init := start.New(conn)
	domain := init.Target

	if !_domain.IsValidDomain(domain) {
		fmt.Println(util.ErrorIcon(), util.Red(start.ErrArgumentsDomain.Error()))
		return
	}
	result, err := init.EmailFormat(domain)
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(err.Error()))
		return
	}
	raw, _ := result.Marshal()
	output.Render(string(raw), init.JSON, init.YAML, init.Output, "format")
}
