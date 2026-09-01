package cmd

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/tomba-io/go/tomba"
	"github.com/tomba-io/tomba/pkg/auth"
	"github.com/tomba-io/tomba/pkg/config"
	"github.com/tomba-io/tomba/pkg/output"
	"github.com/tomba-io/tomba/pkg/start"
	"github.com/tomba-io/tomba/pkg/util"
	_key "github.com/tomba-io/tomba/pkg/validation/key"
)

var useAPIKey bool

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:     "login",
	Short:   "Sign in to Tomba account",
	Long:    Long,
	Example: loginExample,
	Run:     loginRun,
}

func init() {
	loginCmd.Flags().BoolVar(&useAPIKey, "api-key", false, "Login with API key and secret instead of browser-based OAuth")
}

type Prompts struct {
	Prompt promptui.Prompt
	Name   string
}

// loginRun the actual work login
func loginRun(cmd *cobra.Command, args []string) {
	conn.Use = "login"
	init := start.New(conn)
	if init.Key != "" || init.Secret != "" {
		fmt.Println(util.WarningIcon(), util.Yellow("Please logout to login."))
		return
	}

	// Check if already logged in via OAuth
	conf, _ := config.ReadConfigFile()
	if conf != nil && conf.AuthMethod == "oauth" && conf.AccessToken != "" {
		fmt.Println(util.WarningIcon(), util.Yellow("Please logout to login."))
		return
	}

	if useAPIKey {
		loginWithAPIKey()
	} else {
		loginWithDeviceFlow()
	}
}

// loginWithDeviceFlow performs OAuth 2.0 Device Authorization Grant
func loginWithDeviceFlow() {
	fmt.Println(util.InfoIcon(), "Requesting device authorization...")

	dcResp, err := auth.RequestDeviceCode()
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(err.Error()))
		return
	}

	fmt.Println()
	fmt.Println("  To sign in, visit:", util.Green(dcResp.VerificationURI))
	fmt.Println("  And enter code:   ", util.Green(dcResp.UserCode))
	fmt.Println()
	fmt.Println("  Or open:", util.Green(dcResp.VerificationURIComplete))
	fmt.Println()

	// Try to open browser automatically
	if err := auth.OpenBrowser(dcResp.VerificationURIComplete); err != nil {
		fmt.Println(util.WarningIcon(), util.Yellow("Could not open browser automatically. Please open the URL above."))
	}

	fmt.Println(util.InfoIcon(), "Waiting for authorization...")

	tokenResp, err := auth.PollForToken(dcResp.DeviceCode, dcResp.Interval, dcResp.ExpiresIn)
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(err.Error()))
		return
	}

	// Save tokens
	if err := auth.SaveTokens(tokenResp); err != nil {
		fmt.Println(util.ErrorIcon(), util.Red("Failed to save tokens: "+err.Error()))
		return
	}

	fmt.Println()
	fmt.Println(util.SuccessIcon(), util.Green("You have successfully logged in to Tomba."))
}

// loginWithAPIKey performs the legacy API key + secret login
func loginWithAPIKey() {
	var vars = map[string]string{}
	validateKey := func(key string) error {
		if !_key.IsValidAPI(key) {
			return start.ErrErrInvalidApiKey
		}
		return nil
	}
	validateSecret := func(secret string) error {
		if !_key.IsValidAPI(secret) {
			return start.ErrErrInvalidApiSecret
		}
		return nil
	}

	var prompts = []Prompts{
		{
			Prompt: promptui.Prompt{
				Label:     "API key",
				Validate:  validateKey,
				Mask:      '*',
				IsConfirm: false,
			},
			Name: "key",
		},
		{
			Prompt: promptui.Prompt{
				Label:     "Secret Key",
				Validate:  validateSecret,
				Mask:      '*',
				IsConfirm: false,
			},
			Name: "secret",
		},
	}

	for _, prompt := range prompts {
		vars[prompt.Name], _ = prompt.Prompt.Run()
	}
	tomba := tomba.New(vars["key"], vars["secret"])

	me, err := tomba.Account()
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(err.Error()))
		return
	}

	raw, _ := me.Marshal()
	json, _ := output.DisplayJSON(string(raw))
	fmt.Println(json)
	// update config with vars
	if err := config.UpdateConfig(config.Config{
		Key:        vars["key"],
		Secret:     vars["secret"],
		AuthMethod: "apikey",
	}); err != nil {
		fmt.Println("Error updating config file:", err)
		return
	}
	fmt.Println(util.SuccessIcon(), util.Green("You have successfully logged in to tomba."))
}
