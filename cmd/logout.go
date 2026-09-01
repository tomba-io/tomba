package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tomba-io/tomba/pkg/config"
	"github.com/tomba-io/tomba/pkg/util"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:     "logout",
	Short:   "delete your current KEY & SECRET API session.",
	Long:    Long,
	Example: logoutExample,
	Run:     logoutRun,
}

// logoutRun the actual work logout
func logoutRun(cmd *cobra.Command, args []string) {
	// Clear all credentials (API key + OAuth tokens)
	if err := config.UpdateConfig(config.Config{
		Key:          "",
		Secret:       "",
		AccessToken:  "",
		RefreshToken: "",
		TokenExpiry:  "",
		AuthMethod:   "",
	}); err != nil {
		fmt.Println("Error updating config file:", err)
		return
	}
	fmt.Println(util.Green("Successfully disconnected."))
}
