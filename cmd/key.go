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
	keyID string
)

// keyCmd represents the key command
var keyCmd = &cobra.Command{
	Use:     "key",
	Short:   "Manage API keys.",
	Long:    Long,
	Example: keyExample,
}

// keyListCmd represents the key list subcommand
var keyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all API keys.",
	Long:  Long,
	Run:   keyListRun,
}

// keyGetCmd represents the key get subcommand
var keyGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a specific API key by ID.",
	Long:  Long,
	Run:   keyGetRun,
}

// keyCreateCmd represents the key create subcommand
var keyCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new API key.",
	Long:  Long,
	Run:   keyCreateRun,
}

// keyDeleteCmd represents the key delete subcommand
var keyDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an API key by ID.",
	Long:  Long,
	Run:   keyDeleteRun,
}

// keyResetCmd represents the key reset subcommand
var keyResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset an API key by ID.",
	Long:  Long,
	Run:   keyResetRun,
}

func init() {
	keyCmd.AddCommand(keyListCmd, keyGetCmd, keyCreateCmd, keyDeleteCmd, keyResetCmd)

	keyGetCmd.Flags().StringVar(&keyID, "id", "", "API key ID.")
	_ = keyGetCmd.MarkFlagRequired("id")

	keyDeleteCmd.Flags().StringVar(&keyID, "id", "", "API key ID.")
	_ = keyDeleteCmd.MarkFlagRequired("id")

	keyResetCmd.Flags().StringVar(&keyID, "id", "", "API key ID.")
	_ = keyResetCmd.MarkFlagRequired("id")
}

// keyListRun the actual work key list
func keyListRun(cmd *cobra.Command, args []string) {
	init := start.New(conn)
	raw, err := init.ListKeys()
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(err.Error()))
		return
	}
	output.Render(string(raw), init.JSON, init.YAML, init.Output, "key-list")
}

// keyGetRun the actual work key get
func keyGetRun(cmd *cobra.Command, args []string) {
	init := start.New(conn)
	raw, err := init.GetKey(keyID)
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(err.Error()))
		return
	}
	output.Render(string(raw), init.JSON, init.YAML, init.Output, "key-get")
}

// keyCreateRun the actual work key create
func keyCreateRun(cmd *cobra.Command, args []string) {
	init := start.New(conn)
	params := tomba.Params{}
	raw, err := init.CreateKey(params)
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(err.Error()))
		return
	}
	output.Render(string(raw), init.JSON, init.YAML, init.Output, "key-create")
}

// keyDeleteRun the actual work key delete
func keyDeleteRun(cmd *cobra.Command, args []string) {
	init := start.New(conn)
	raw, err := init.DeleteKey(keyID)
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(err.Error()))
		return
	}
	output.Render(string(raw), init.JSON, init.YAML, init.Output, "key-delete")
}

// keyResetRun the actual work key reset
func keyResetRun(cmd *cobra.Command, args []string) {
	init := start.New(conn)
	raw, err := init.ResetKey(keyID)
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(err.Error()))
		return
	}
	output.Render(string(raw), init.JSON, init.YAML, init.Output, "key-reset")
}
