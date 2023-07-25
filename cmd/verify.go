package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tomba-io/email/pkg/output"
	"github.com/tomba-io/email/pkg/start"
	"github.com/tomba-io/email/pkg/util"
	_email "github.com/tomba-io/email/pkg/validation/email"
)

// verifyCmd represents the verify command
// see https://developer.tomba.io/#email-verifier
var verifyCmd = &cobra.Command{
	Use:     "verify",
	Aliases: []string{"t"},
	Short:   "Verify the deliverability of an email address.",
	Run:     verifyRun,
	Example: verifyExample,
}

// verifyRun the actual work verify
func verifyRun(cmd *cobra.Command, args []string) {
	fmt.Println(Long)
	init := start.New(conn)
	if init.Key == "" || init.Secret == "" {
		fmt.Println(util.WarningIcon(), util.Yellow(start.ErrErrInvalidNoLogin.Error()))
		return
	}
	email := init.Target
	if !_email.IsValidEmail(email) {
		fmt.Println(util.ErrorIcon(), util.Red(start.ErrArgumentEmail.Error()))
		return
	}

	result, err := init.Tomba.EmailVerifier(email)
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(start.ErrErrInvalidLogin.Error()))
	}
	if init.JSON {
		raw, _ := result.Marshal()
		json, _ := output.DisplayJSON(string(raw))
		fmt.Println(json)
		return
	}
	if init.YAML {
		raw, _ := result.Marshal()
		json, _ := output.DisplayYAML(string(raw))
		fmt.Println(json)
		return
	}

}
