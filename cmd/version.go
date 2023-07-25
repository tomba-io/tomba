package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tomba-io/email/pkg/util"
	"github.com/tomba-io/email/pkg/version"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version number and build information.",
	Long:  fmt.Sprintf("utility to search or verify lists of email addresses in minutes can significantly improve productivity and efficiency for individuals and businesses dealing with large email databases.\n\n%s", util.RandomBanner()),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Long)
		fmt.Println(version.String())
	},
}
