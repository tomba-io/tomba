package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/tomba-io/tomba/pkg/output"
	"github.com/tomba-io/tomba/pkg/start"
	"github.com/tomba-io/tomba/pkg/util"
	_domain "github.com/tomba-io/tomba/pkg/validation/domain"
)

// locationCmd represents the location command
var locationCmd = &cobra.Command{
	Use:     "location",
	Short:   "Get the employees count by location for a domain.",
	Long:    Long,
	Example: locationExample,
	Run:     locationRun,
}

// locationRun the actual work location
func locationRun(cmd *cobra.Command, args []string) {
	init := start.New(conn)
	domain := init.Target

	if !_domain.IsValidDomain(domain) {
		fmt.Println(util.ErrorIcon(), util.Red(start.ErrArgumentsDomain.Error()))
		return
	}
	result, err := init.EmployeesCount(domain)
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(err.Error()))
		return
	}
	raw, _ := result.Marshal()
	output.Render(string(raw), init.JSON, init.YAML, init.Output, "location")
}
