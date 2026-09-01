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
	flagPage   int
	flagLimit  int
	flagEmail  string
	flagReason string
)

// flagCmd represents the flag command
var flagCmd = &cobra.Command{
	Use:     "flag",
	Short:   "Manage flags.",
	Long:    Long,
	Example: flagExample,
}

// flagListCmd represents the flag list subcommand
var flagListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all flags.",
	Long:  Long,
	Run:   flagListRun,
}

// flagCreateCmd represents the flag create subcommand
var flagCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new flag.",
	Long:  Long,
	Run:   flagCreateRun,
}

func init() {
	flagCmd.AddCommand(flagListCmd, flagCreateCmd)

	flagListCmd.Flags().IntVar(&flagPage, "page", 1, "Page number for pagination.")
	flagListCmd.Flags().IntVar(&flagLimit, "limit", 10, "Number of flags per page.")

	flagCreateCmd.Flags().StringVar(&flagEmail, "email", "", "Email address to flag.")
	flagCreateCmd.Flags().StringVar(&flagReason, "reason", "", "Reason for flagging.")
	_ = flagCreateCmd.MarkFlagRequired("email")
}

// flagListRun the actual work flag list
func flagListRun(cmd *cobra.Command, args []string) {
	init := start.New(conn)
	params := tomba.Params{}
	if flagPage > 0 {
		params["page"] = fmt.Sprint(flagPage)
	}
	if flagLimit > 0 {
		params["limit"] = fmt.Sprint(flagLimit)
	}
	raw, err := init.ListFlags(params)
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(err.Error()))
		return
	}
	output.Render(string(raw), init.JSON, init.YAML, init.Output, "flag-list")
}

// flagCreateRun the actual work flag create
func flagCreateRun(cmd *cobra.Command, args []string) {
	init := start.New(conn)
	params := tomba.Params{"email": flagEmail}
	if flagReason != "" {
		params["reason"] = flagReason
	}
	raw, err := init.CreateFlag(params)
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(err.Error()))
		return
	}
	output.Render(string(raw), init.JSON, init.YAML, init.Output, "flag-create")
}
