package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tomba-io/email/pkg/output"
	"github.com/tomba-io/email/pkg/start"
	"github.com/tomba-io/email/pkg/util"
	_key "github.com/tomba-io/email/pkg/validation/key"
	_url "github.com/tomba-io/email/pkg/validation/url"
)

// authorCmd represents the author command
// see https://developer.tomba.io/#author-finder
var authorCmd = &cobra.Command{
	Use:     "author",
	Aliases: []string{"a"},
	Short:   "Instantly discover the email addresses of article authors.",
	Run:     authorRun,
	Example: authorExample,
}

// authorRun the actual work author
func authorRun(cmd *cobra.Command, args []string) {
	fmt.Println(Long)
	init := start.New(conn)
	if init.Key == "" || init.Secret == "" {
		fmt.Println(util.WarningIcon(), util.Yellow(start.ErrErrInvalidNoLogin.Error()))
		return
	}
	if !_key.IsValidAPI(init.Key) && !_key.IsValidAPI(init.Secret) {
		fmt.Println(util.WarningIcon(), util.Yellow(start.ErrErrInvalidLogin.Error()))
		return
	}
	url := init.Target
	if !_url.IsValidURL(url) {
		fmt.Println(util.ErrorIcon(), util.Red(start.ErrArgumentsURL.Error()))
		return
	}
	result, err := init.Tomba.AuthorFinder(url)
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
