package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tomba-io/go/tomba"

	"github.com/tomba-io/tomba/pkg/output"
	"github.com/tomba-io/tomba/pkg/start"
	"github.com/tomba-io/tomba/pkg/util"
	_domain "github.com/tomba-io/tomba/pkg/validation/domain"
	_email "github.com/tomba-io/tomba/pkg/validation/email"
	_url "github.com/tomba-io/tomba/pkg/validation/url"
)

var (
	phoneEmail    string
	phoneDomain   string
	phoneLinkedin string
	phoneFull     bool
)

// phoneFinderCmd represents the phone-finder command
// see https://docs.tomba.io/api/phone#phone-finder
var phoneFinderCmd = &cobra.Command{
	Use:     "phone-finder",
	Aliases: []string{"pf"},
	Short:   "Search for phone numbers given an email, domain, or LinkedIn URL.",
	Long:    Long,
	Run:     phoneFinderRun,
	Example: phoneFinderExample,
}

func init() {
	phoneFinderCmd.Flags().StringVar(&phoneEmail, "email", "", "Email address to search phone for.")
	phoneFinderCmd.Flags().StringVar(&phoneDomain, "domain", "", "Domain to search phone for.")
	phoneFinderCmd.Flags().StringVar(&phoneLinkedin, "linkedin", "", "LinkedIn URL to search phone for.")
	phoneFinderCmd.Flags().BoolVar(&phoneFull, "full", false, "Get all phone numbers.")
}

// phoneFinderRun the actual work phone-finder
func phoneFinderRun(cmd *cobra.Command, args []string) {
	fmt.Println(Long)
	init := start.New(conn)

	// Build params based on provided flags
	params := tomba.Params{}

	if phoneEmail != "" {
		if !_email.IsValidEmail(phoneEmail) {
			fmt.Println(util.ErrorIcon(), util.Red(start.ErrArgumentEmail.Error()))
			return
		}
		params["email"] = phoneEmail
	} else if phoneDomain != "" {
		if !_domain.IsValidDomain(phoneDomain) {
			fmt.Println(util.ErrorIcon(), util.Red(start.ErrArgumentsDomain.Error()))
			return
		}
		params["domain"] = phoneDomain
	} else if phoneLinkedin != "" {
		if !_url.IsValidLinkedInProfile(phoneLinkedin) {
			fmt.Println(util.ErrorIcon(), util.Red(start.ErrArgumentsLinkedinURL.Error()))
			return
		}
		params["linkedin"] = phoneLinkedin
	} else {
		fmt.Println(util.ErrorIcon(), util.Red("At least one of --email, --domain, or --linkedin is required."))
		return
	}

	if phoneFull {
		params["full"] = true
	}

	result, err := init.Tomba.PhoneFinder(params)
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
