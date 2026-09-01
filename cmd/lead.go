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
	leadID        string
	leadEmail     string
	leadListID    string
	leadFirstName string
	leadLastName  string
	leadPage      int
	leadLimit     int
	leadDomain    string
)

// leadCmd represents the lead command
var leadCmd = &cobra.Command{
	Use:     "lead",
	Short:   "Manage leads.",
	Long:    Long,
	Example: leadExample,
}

// leadListCmd represents the lead list subcommand
var leadListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all leads.",
	Long:  Long,
	Run:   leadListRun,
}

// leadGetCmd represents the lead get subcommand
var leadGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a specific lead by ID.",
	Long:  Long,
	Run:   leadGetRun,
}

// leadCreateCmd represents the lead create subcommand
var leadCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new lead.",
	Long:  Long,
	Run:   leadCreateRun,
}

// leadUpdateCmd represents the lead update subcommand
var leadUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing lead.",
	Long:  Long,
	Run:   leadUpdateRun,
}

// leadDeleteCmd represents the lead delete subcommand
var leadDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a lead by ID.",
	Long:  Long,
	Run:   leadDeleteRun,
}

func init() {
	leadCmd.AddCommand(leadListCmd, leadGetCmd, leadCreateCmd, leadUpdateCmd, leadDeleteCmd)

	leadListCmd.Flags().IntVar(&leadPage, "page", 1, "Page number for pagination.")
	leadListCmd.Flags().IntVar(&leadLimit, "limit", 10, "Number of leads per page.")
	leadListCmd.Flags().StringVar(&leadDomain, "domain", "", "Filter leads by domain.")

	leadGetCmd.Flags().StringVar(&leadID, "id", "", "Lead ID.")
	_ = leadGetCmd.MarkFlagRequired("id")

	leadCreateCmd.Flags().StringVar(&leadEmail, "email", "", "Lead email address.")
	leadCreateCmd.Flags().StringVar(&leadListID, "list-id", "", "Lead list ID.")
	leadCreateCmd.Flags().StringVar(&leadFirstName, "first-name", "", "Lead first name.")
	leadCreateCmd.Flags().StringVar(&leadLastName, "last-name", "", "Lead last name.")
	_ = leadCreateCmd.MarkFlagRequired("email")
	_ = leadCreateCmd.MarkFlagRequired("list-id")

	leadUpdateCmd.Flags().StringVar(&leadID, "id", "", "Lead ID.")
	leadUpdateCmd.Flags().StringVar(&leadEmail, "email", "", "Lead email address.")
	leadUpdateCmd.Flags().StringVar(&leadFirstName, "first-name", "", "Lead first name.")
	leadUpdateCmd.Flags().StringVar(&leadLastName, "last-name", "", "Lead last name.")
	_ = leadUpdateCmd.MarkFlagRequired("id")

	leadDeleteCmd.Flags().StringVar(&leadID, "id", "", "Lead ID.")
	_ = leadDeleteCmd.MarkFlagRequired("id")
}

// leadListRun the actual work lead list
func leadListRun(cmd *cobra.Command, args []string) {
	init := start.New(conn)
	params := tomba.Params{}
	if leadPage > 0 {
		params["page"] = fmt.Sprint(leadPage)
	}
	if leadLimit > 0 {
		params["limit"] = fmt.Sprint(leadLimit)
	}
	if leadDomain != "" {
		params["domain"] = leadDomain
	}
	raw, err := init.ListLeads(params)
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(err.Error()))
		return
	}
	output.Render(string(raw), init.JSON, init.YAML, init.Output, "lead-list")
}

// leadGetRun the actual work lead get
func leadGetRun(cmd *cobra.Command, args []string) {
	init := start.New(conn)
	raw, err := init.GetLead(leadID)
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(err.Error()))
		return
	}
	output.Render(string(raw), init.JSON, init.YAML, init.Output, "lead-get")
}

// leadCreateRun the actual work lead create
func leadCreateRun(cmd *cobra.Command, args []string) {
	init := start.New(conn)
	params := tomba.Params{
		"email":   leadEmail,
		"list_id": leadListID,
	}
	if leadFirstName != "" {
		params["first_name"] = leadFirstName
	}
	if leadLastName != "" {
		params["last_name"] = leadLastName
	}
	raw, err := init.CreateLead(params)
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(err.Error()))
		return
	}
	output.Render(string(raw), init.JSON, init.YAML, init.Output, "lead-create")
}

// leadUpdateRun the actual work lead update
func leadUpdateRun(cmd *cobra.Command, args []string) {
	init := start.New(conn)
	params := tomba.Params{}
	if leadEmail != "" {
		params["email"] = leadEmail
	}
	if leadFirstName != "" {
		params["first_name"] = leadFirstName
	}
	if leadLastName != "" {
		params["last_name"] = leadLastName
	}
	raw, err := init.UpdateLead(leadID, params)
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(err.Error()))
		return
	}
	output.Render(string(raw), init.JSON, init.YAML, init.Output, "lead-update")
}

// leadDeleteRun the actual work lead delete
func leadDeleteRun(cmd *cobra.Command, args []string) {
	init := start.New(conn)
	raw, err := init.DeleteLead(leadID)
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(err.Error()))
		return
	}
	output.Render(string(raw), init.JSON, init.YAML, init.Output, "lead-delete")
}
