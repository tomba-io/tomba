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
	leadsListID   string
	leadsListName string
)

// leadsListCmd represents the leads-list command
var leadsListCmd = &cobra.Command{
	Use:     "leads-list",
	Short:   "Manage leads lists.",
	Long:    Long,
	Example: leadsListExample,
}

// leadsListListCmd represents the leads-list list subcommand
var leadsListListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all leads lists.",
	Long:  Long,
	Run:   leadsListListRun,
}

// leadsListGetCmd represents the leads-list get subcommand
var leadsListGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a specific leads list by ID.",
	Long:  Long,
	Run:   leadsListGetRun,
}

// leadsListCreateCmd represents the leads-list create subcommand
var leadsListCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new leads list.",
	Long:  Long,
	Run:   leadsListCreateRun,
}

// leadsListUpdateCmd represents the leads-list update subcommand
var leadsListUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing leads list.",
	Long:  Long,
	Run:   leadsListUpdateRun,
}

// leadsListDeleteCmd represents the leads-list delete subcommand
var leadsListDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a leads list by ID.",
	Long:  Long,
	Run:   leadsListDeleteRun,
}

func init() {
	leadsListCmd.AddCommand(leadsListListCmd, leadsListGetCmd, leadsListCreateCmd, leadsListUpdateCmd, leadsListDeleteCmd)

	leadsListGetCmd.Flags().StringVar(&leadsListID, "id", "", "Leads list ID.")
	_ = leadsListGetCmd.MarkFlagRequired("id")

	leadsListCreateCmd.Flags().StringVar(&leadsListName, "name", "", "Leads list name.")
	_ = leadsListCreateCmd.MarkFlagRequired("name")

	leadsListUpdateCmd.Flags().StringVar(&leadsListID, "id", "", "Leads list ID.")
	leadsListUpdateCmd.Flags().StringVar(&leadsListName, "name", "", "Leads list name.")
	_ = leadsListUpdateCmd.MarkFlagRequired("id")
	_ = leadsListUpdateCmd.MarkFlagRequired("name")

	leadsListDeleteCmd.Flags().StringVar(&leadsListID, "id", "", "Leads list ID.")
	_ = leadsListDeleteCmd.MarkFlagRequired("id")
}

// leadsListListRun the actual work leads-list list
func leadsListListRun(cmd *cobra.Command, args []string) {
	init := start.New(conn)
	params := tomba.Params{}
	raw, err := init.ListLeadsLists(params)
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(err.Error()))
		return
	}
	output.Render(string(raw), init.JSON, init.YAML, init.Output, "leads-list-list")
}

// leadsListGetRun the actual work leads-list get
func leadsListGetRun(cmd *cobra.Command, args []string) {
	init := start.New(conn)
	raw, err := init.GetLeadsList(leadsListID)
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(err.Error()))
		return
	}
	output.Render(string(raw), init.JSON, init.YAML, init.Output, "leads-list-get")
}

// leadsListCreateRun the actual work leads-list create
func leadsListCreateRun(cmd *cobra.Command, args []string) {
	init := start.New(conn)
	params := tomba.Params{"name": leadsListName}
	raw, err := init.CreateLeadsList(params)
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(err.Error()))
		return
	}
	output.Render(string(raw), init.JSON, init.YAML, init.Output, "leads-list-create")
}

// leadsListUpdateRun the actual work leads-list update
func leadsListUpdateRun(cmd *cobra.Command, args []string) {
	init := start.New(conn)
	params := tomba.Params{"name": leadsListName}
	raw, err := init.UpdateLeadsList(leadsListID, params)
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(err.Error()))
		return
	}
	output.Render(string(raw), init.JSON, init.YAML, init.Output, "leads-list-update")
}

// leadsListDeleteRun the actual work leads-list delete
func leadsListDeleteRun(cmd *cobra.Command, args []string) {
	init := start.New(conn)
	raw, err := init.DeleteLeadsList(leadsListID)
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(err.Error()))
		return
	}
	output.Render(string(raw), init.JSON, init.YAML, init.Output, "leads-list-delete")
}
