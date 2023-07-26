package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tomba-io/tomba/pkg/config"
	"github.com/tomba-io/tomba/pkg/start"
	"github.com/tomba-io/tomba/pkg/util"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "delete your current KEY & SECRET API session.",
	Long:  Long,
	Run:   logoutRun,
}

// loginRun the actual work login
func logoutRun(cmd *cobra.Command, args []string) {
	init := start.New(conn)
	if init.Key == "" && init.Secret == "" {
		fmt.Println(util.WarningIcon(), util.Yellow(start.ErrErrInvalidNoLogin.Error()))
		return
	}
	// update config
	if err := config.UpdateConfig(config.Config{
		Key:    "",
		Secret: "",
	}); err != nil {
		fmt.Println("Error updating config file:", err)
		return
	}
	fmt.Println(util.Green("Successfully disconnected."))
}
