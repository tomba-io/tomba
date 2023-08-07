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
// see https://developer.tomba.io/#email-finder
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
	fmt.Println(Long)
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

	result, err := init.Tomba.EmailFinder(params)
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(start.ErrErrInvalidLogin.Error()))
		return
	}
	if result.Data.Email != "" {
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
	fmt.Println(util.WarningIcon(), util.Yellow("Why doesn't the Email Finder return any result? https://help.tomba.io/en/questions/why-doesn-t-the-email-finder-return-any-result"))
}
