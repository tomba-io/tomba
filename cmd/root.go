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
	rootCmd.AddCommand(authorCmd, countCmd, enrichCmd, finderCmd, httpCmd, linkedinCmd, logsCmd, logoutCmd, loginCmd, searchCmd, statusCmd, usageCmd, verifyCmd, versionCmd)
	rootCmd.PersistentFlags().StringVarP(&conn.Key, "key", "k", "", "Tomba API KEY.")
	rootCmd.PersistentFlags().StringVarP(&conn.Secret, "secret", "s", "", "Tomba API SECRET.")
	rootCmd.PersistentFlags().StringVarP(&conn.Target, "target", "t", "", "TARGET SPECIFICATION Can pass email, Domain, URL, Linkedin URL.")
	rootCmd.PersistentFlags().StringVarP(&conn.Output, "output", "o", "", "Save the results to file.")
	rootCmd.PersistentFlags().IntVarP(&conn.Port, "port", "p", 3000, "Sets the port on which the HTTP server should bind.")
	rootCmd.PersistentFlags().BoolVarP(&conn.JSON, "json", "j", true, "output JSON format.")
	rootCmd.PersistentFlags().BoolVarP(&conn.YAML, "yaml", "y", false, "output YAML format.")
	searchCmd.PersistentFlags().IntVar(&conn.Page, "page", 1, "Specifies the number of email addresses to skip. The default is 1.")
	searchCmd.PersistentFlags().IntVar(&conn.Limit, "limit", 10, "Specifies the max number of email addresses to return. The default is 10. valid number(10,20,50)")
	searchCmd.PersistentFlags().StringVar(&conn.Department, "department", "", "Get only email addresses for people working in the selected department(s).")
	finderCmd.PersistentFlags().StringVar(&conn.FullName, "full", "", "The person's full name")
	finderCmd.PersistentFlags().StringVarP(&conn.FirstName, "fist", "f", "", "The person's first name. It doesn't need to be in lowercase..")
	finderCmd.PersistentFlags().StringVarP(&conn.LastName, "last", "l", "", "The person's last name. It doesn't need to be in lowercase..")
}
