package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tomba-io/go/tomba"

	"github.com/tomba-io/tomba/pkg/output"
	"github.com/tomba-io/tomba/pkg/start"
	"github.com/tomba-io/tomba/pkg/util"
)

var (
	validatorPhone       string
	validatorCountryCode string
)

// phoneValidatorCmd represents the phone-validator command
// see https://docs.tomba.io/api/phone#phone-validator
var phoneValidatorCmd = &cobra.Command{
	Use:     "phone-validator",
	Aliases: []string{"pv"},
	Short:   "Validate phone numbers.",
	Long:    Long,
	Run:     phoneValidatorRun,
	Example: phoneValidatorExample,
}

func init() {
	phoneValidatorCmd.Flags().StringVar(&validatorPhone, "phone", "", "Phone number to validate (e.g., +14155552671).")
	phoneValidatorCmd.Flags().StringVar(&validatorCountryCode, "country-code", "", "Country code for parsing local numbers (e.g., US).")
	phoneValidatorCmd.MarkFlagRequired("phone")
}

// phoneValidatorRun the actual work phone-validator
func phoneValidatorRun(cmd *cobra.Command, args []string) {
	fmt.Println(Long)
	init := start.New(conn)

	if validatorPhone == "" {
		fmt.Println(util.ErrorIcon(), util.Red("--phone flag is required."))
		return
	}

	params := tomba.Params{
		"phone": validatorPhone,
	}

	if validatorCountryCode != "" {
		params["country_code"] = validatorCountryCode
	}

	result, err := init.Tomba.PhoneValidator(params)
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
