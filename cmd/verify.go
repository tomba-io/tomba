package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tomba-io/go/tomba"

	"github.com/tomba-io/tomba/pkg/output"
	"github.com/tomba-io/tomba/pkg/start"
	"github.com/tomba-io/tomba/pkg/util"
	_email "github.com/tomba-io/tomba/pkg/validation/email"
)

// verifyCmd represents the verify command
// see https://docs.tomba.io/api/verifier#email-verifier
var verifyCmd = &cobra.Command{
	Use:     "verify",
	Aliases: []string{"t"},
	Short:   "Verify the deliverability of an email address.",
	Long:    Long,
	Run:     verifyRun,
	Example: verifyExample,
}

// verifyRun the actual work verify
func verifyRun(cmd *cobra.Command, args []string) {
	init := start.New(conn)
	email := init.Target

	if !_email.IsValidEmail(email) {
		fmt.Println(util.ErrorIcon(), util.Red(start.ErrArgumentEmail.Error()))
		return
	}

	result, err := init.Tomba.EmailVerifier(tomba.Params{"email": email})
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(err.Error()))
		return
	}
	if result.Data.Email.Email != "" {
		if result.Data.Email.Disposable {
			fmt.Println(util.WarningIcon(), util.Bold("The domain name is used by a disposable email addresses provider."))
			fmt.Println(util.WarningIcon(), util.Yellow("Tomba is designed to contact other professionals. This email is used to create personal email addresses so we don't the verification. 💡"))
			return
		}
		if result.Data.Email.Webmail {
			fmt.Println(util.WarningIcon(), util.Bold("The domain name  is webmail provider."))
			fmt.Println(util.WarningIcon(), util.Yellow("Tomba is designed to contact other professionals. This email is used to create personal email addresses so we don't the verification. 💡"))
			return
		}
		raw, _ := result.Marshal()
		output.Render(string(raw), init.JSON, init.YAML, init.Output, "verify")
		return
	}
	fmt.Println(util.WarningIcon(), util.Yellow("The Email Verification failed because of an unexpected response from the remote SMTP server. This failure is outside of our control. We advise you to retry later."))
}
