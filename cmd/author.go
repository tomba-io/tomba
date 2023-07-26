package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tomba-io/tomba/pkg/output"
	"github.com/tomba-io/tomba/pkg/start"
	"github.com/tomba-io/tomba/pkg/util"
	_url "github.com/tomba-io/tomba/pkg/validation/url"
)

// authorCmd represents the author command
// see https://developer.tomba.io/#author-finder
var authorCmd = &cobra.Command{
	Use:     "author",
	Aliases: []string{"a"},
	Short:   "Instantly discover the email addresses of article authors.",
	Long:    Long,
	Run:     authorRun,
	Example: authorExample,
}

// authorRun the actual work author
func authorRun(cmd *cobra.Command, args []string) {
	fmt.Println(Long)
	init := start.New(conn)
	url := init.Target

	if !_url.IsValidURL(url) {
		fmt.Println(util.ErrorIcon(), util.Red(start.ErrArgumentsURL.Error()))
		return
	}
	result, err := init.Tomba.AuthorFinder(url)
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
	fmt.Println(util.WarningIcon(), util.Yellow("Why doesn't the author finder return any result? https://help.tomba.io/en/questions/reasons-why-author-finder-don-t-find-any-result"))
}
