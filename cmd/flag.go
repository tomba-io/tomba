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
	flagPage     int
	flagLimit    int
	flagType     string
	flagValue    string
	flagReason   string
	flagComment  string
)

// flagCmd represents the flag command
var flagCmd = &cobra.Command{
	Use:     "flag",
	Short:   "Report incorrect data for credit recovery.",
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
	Short: "Create a new flag to report incorrect data.",
	Long: `Create a new flag to report incorrect data.

Flag types: email, organization, phone, author_url, website

Reasons by type:
  email:        hard_bounce, invalid_email, wrong_person, outdated, other
  organization: wrong_company, outdated, other
  phone:        wrong_phone, outdated, other
  author_url:   broken_url, wrong_person, outdated, other
  website:      broken_url, wrong_company, outdated, other`,
	Run: flagCreateRun,
	Example: `  tomba flag create --type email --value bounce@example.com --reason hard_bounce
  tomba flag create --type organization --value example.com --reason wrong_company
  tomba flag create --type phone --value "+1234567890" --reason wrong_phone --comment "Number disconnected"`,
}

func init() {
	flagCmd.AddCommand(flagListCmd, flagCreateCmd)

	flagListCmd.Flags().IntVar(&flagPage, "page", 1, "Page number for pagination.")
	flagListCmd.Flags().IntVar(&flagLimit, "limit", 10, "Number of flags per page.")

	flagCreateCmd.Flags().StringVar(&flagType, "type", "", "Flag type: email, organization, phone, author_url, website (required).")
	flagCreateCmd.Flags().StringVar(&flagValue, "value", "", "The flagged item: email, domain, phone, URL (required).")
	flagCreateCmd.Flags().StringVar(&flagReason, "reason", "", "Reason for flagging (required, depends on type).")
	flagCreateCmd.Flags().StringVar(&flagComment, "comment", "", "Optional additional details (max 1000 chars).")
	_ = flagCreateCmd.MarkFlagRequired("type")
	_ = flagCreateCmd.MarkFlagRequired("value")
	_ = flagCreateCmd.MarkFlagRequired("reason")
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
	params := tomba.Params{
		"flag_type": flagType,
		"value":     flagValue,
		"reason":    flagReason,
	}
	if flagComment != "" {
		params["comment"] = flagComment
	}
	raw, err := init.CreateFlag(params)
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(err.Error()))
		return
	}
	output.Render(string(raw), init.JSON, init.YAML, init.Output, "flag-create")
}
