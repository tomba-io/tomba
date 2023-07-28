package cmd

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/tomba-io/go/tomba"
	"github.com/tomba-io/tomba/pkg/config"
	"github.com/tomba-io/tomba/pkg/output"
	"github.com/tomba-io/tomba/pkg/start"
	"github.com/tomba-io/tomba/pkg/util"
	_key "github.com/tomba-io/tomba/pkg/validation/key"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Sign in to Tomba account",
	Long:  Long,
	Run:   loginRun,
}

type Prompts struct {
	Prompt promptui.Prompt
	Name   string
}

// loginRun the actual work login
func loginRun(cmd *cobra.Command, args []string) {
	fmt.Println(Long)
	conn.Use = "login"
	init := start.New(conn)
	if init.Key != "" || init.Secret != "" {
		fmt.Println(util.WarningIcon(), util.Yellow("Please logout to login."))
		return
	}

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
		fmt.Println(util.ErrorIcon(), util.Red(start.ErrErrInvalidLogin.Error()))
		return
	}

	raw, _ := me.Marshal()
	json, _ := output.DisplayJSON(string(raw))
	fmt.Println(json)
	// update config with vars
	if err := config.UpdateConfig(config.Config{
		Key:    vars["key"],
		Secret: vars["secret"],
	}); err != nil {
		fmt.Println("Error updating config file:", err)
		return
	}
	fmt.Println(util.SuccessIcon(), util.Green("You have successfully logged in to tomba."))
}
