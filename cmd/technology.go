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
// see https://docs.tomba.io/api/~endpoints#technology
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
	fmt.Println(Long)
	init := start.New(conn)
	domain := init.Target

	if !_domain.IsValidDomain(domain) {
		fmt.Println(util.ErrorIcon(), util.Red(start.ErrArgumentsDomain.Error()))
		return
	}

	result, err := init.Tomba.TechnologyCheck(domain)
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(start.ErrErrInvalidLogin.Error()))
		return
	}

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
}
