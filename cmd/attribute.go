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
	attributeID   string
	attributeName string
	attributeType string
)

// attributeCmd represents the attribute command
var attributeCmd = &cobra.Command{
	Use:     "attribute",
	Short:   "Manage custom attributes.",
	Long:    Long,
	Example: attributeExample,
}

// attributeListCmd represents the attribute list subcommand
var attributeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all custom attributes.",
	Long:  Long,
	Run:   attributeListRun,
}

// attributeGetCmd represents the attribute get subcommand
var attributeGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a specific attribute by ID.",
	Long:  Long,
	Run:   attributeGetRun,
}

// attributeCreateCmd represents the attribute create subcommand
var attributeCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new custom attribute.",
	Long:  Long,
	Run:   attributeCreateRun,
}

// attributeUpdateCmd represents the attribute update subcommand
var attributeUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing attribute.",
	Long:  Long,
	Run:   attributeUpdateRun,
}

// attributeDeleteCmd represents the attribute delete subcommand
var attributeDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an attribute by ID.",
	Long:  Long,
	Run:   attributeDeleteRun,
}

func init() {
	attributeCmd.AddCommand(attributeListCmd, attributeGetCmd, attributeCreateCmd, attributeUpdateCmd, attributeDeleteCmd)

	attributeGetCmd.Flags().StringVar(&attributeID, "id", "", "Attribute ID.")
	_ = attributeGetCmd.MarkFlagRequired("id")

	attributeCreateCmd.Flags().StringVar(&attributeName, "name", "", "Attribute name.")
	attributeCreateCmd.Flags().StringVar(&attributeType, "type", "", "Attribute type.")
	_ = attributeCreateCmd.MarkFlagRequired("name")
	_ = attributeCreateCmd.MarkFlagRequired("type")

	attributeUpdateCmd.Flags().StringVar(&attributeID, "id", "", "Attribute ID.")
	attributeUpdateCmd.Flags().StringVar(&attributeName, "name", "", "Attribute name.")
	_ = attributeUpdateCmd.MarkFlagRequired("id")
	_ = attributeUpdateCmd.MarkFlagRequired("name")

	attributeDeleteCmd.Flags().StringVar(&attributeID, "id", "", "Attribute ID.")
	_ = attributeDeleteCmd.MarkFlagRequired("id")
}

// attributeListRun the actual work attribute list
func attributeListRun(cmd *cobra.Command, args []string) {
	init := start.New(conn)
	params := tomba.Params{}
	raw, err := init.ListAttributes(params)
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(err.Error()))
		return
	}
	output.Render(string(raw), init.JSON, init.YAML, init.Output, "attribute-list")
}

// attributeGetRun the actual work attribute get
func attributeGetRun(cmd *cobra.Command, args []string) {
	init := start.New(conn)
	raw, err := init.GetAttribute(attributeID)
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(err.Error()))
		return
	}
	output.Render(string(raw), init.JSON, init.YAML, init.Output, "attribute-get")
}

// attributeCreateRun the actual work attribute create
func attributeCreateRun(cmd *cobra.Command, args []string) {
	init := start.New(conn)
	params := tomba.Params{
		"name": attributeName,
		"type": attributeType,
	}
	raw, err := init.CreateAttribute(params)
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(err.Error()))
		return
	}
	output.Render(string(raw), init.JSON, init.YAML, init.Output, "attribute-create")
}

// attributeUpdateRun the actual work attribute update
func attributeUpdateRun(cmd *cobra.Command, args []string) {
	init := start.New(conn)
	params := tomba.Params{"name": attributeName}
	raw, err := init.UpdateAttribute(attributeID, params)
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(err.Error()))
		return
	}
	output.Render(string(raw), init.JSON, init.YAML, init.Output, "attribute-update")
}

// attributeDeleteRun the actual work attribute delete
func attributeDeleteRun(cmd *cobra.Command, args []string) {
	init := start.New(conn)
	raw, err := init.DeleteAttribute(attributeID)
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(err.Error()))
		return
	}
	output.Render(string(raw), init.JSON, init.YAML, init.Output, "attribute-delete")
}
