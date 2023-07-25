package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tomba-io/email/pkg/version"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version number and build information.",
	Long:  Long,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Long)
		fmt.Println(version.String())
	},
}
