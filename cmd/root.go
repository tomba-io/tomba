package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tomba-io/tomba/pkg/config"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tomba",
	Short: "CLI utility to search or verify lists of email addresses in minutes",
	Long:  Long,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(config.InitConfig)
	rootCmd.AddCommand(authorCmd, countCmd, enrichCmd, httpCmd, linkedinCmd, logsCmd, logoutCmd, loginCmd, searchCmd, statusCmd, usageCmd, verifyCmd, versionCmd)
	rootCmd.PersistentFlags().StringVarP(&conn.Key, "key", "k", "", "Tomba API KEY.")
	rootCmd.PersistentFlags().StringVarP(&conn.Secret, "secret", "s", "", "Tomba API SECRET.")
	rootCmd.PersistentFlags().StringVarP(&conn.Target, "target", "t", "", "TARGET SPECIFICATION Can pass email, Domain, URL, Linkedin URL.")
	rootCmd.PersistentFlags().StringVarP(&conn.Output, "output", "o", "", "Save the results to file.")
	rootCmd.PersistentFlags().IntVarP(&conn.Port, "port", "p", 3000, "Sets the port on which the HTTP server should bind.")
	rootCmd.PersistentFlags().BoolVarP(&conn.JSON, "json", "j", true, "output JSON format.")
	rootCmd.PersistentFlags().BoolVarP(&conn.YAML, "yaml", "y", false, "output YAML format.")
}
