package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tomba-io/go/tomba"

	"github.com/tomba-io/tomba/pkg/output"
	"github.com/tomba-io/tomba/pkg/start"
	"github.com/tomba-io/tomba/pkg/util"
)

// enrichmentCmd represents the enrichment command
var enrichmentCmd = &cobra.Command{
	Use:     "enrichment",
	Short:   "Enrich person, company, or combined data.",
	Long:    Long,
	Example: enrichmentExample,
}

// enrichmentPersonCmd represents the enrichment person subcommand
var enrichmentPersonCmd = &cobra.Command{
	Use:   "person",
	Short: "Look up person data based on an email address.",
	Long:  Long,
	Run:   enrichmentPersonRun,
}

// enrichmentCompanyCmd represents the enrichment company subcommand
var enrichmentCompanyCmd = &cobra.Command{
	Use:   "company",
	Short: "Look up company data based on a domain name.",
	Long:  Long,
	Run:   enrichmentCompanyRun,
}

// enrichmentCombinedCmd represents the enrichment combined subcommand
var enrichmentCombinedCmd = &cobra.Command{
	Use:   "combined",
	Short: "Look up both person and company data based on an email address.",
	Long:  Long,
	Run:   enrichmentCombinedRun,
}

func init() {
	enrichmentCmd.AddCommand(enrichmentPersonCmd, enrichmentCompanyCmd, enrichmentCombinedCmd)
}

// enrichmentPersonRun the actual work enrichment person
func enrichmentPersonRun(cmd *cobra.Command, args []string) {
	init := start.New(conn)
	email := init.Target

	if email == "" {
		fmt.Println(util.ErrorIcon(), util.Red("email is required"))
		return
	}
	result, err := init.PersonFind(tomba.Params{"email": email})
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(err.Error()))
		return
	}
	raw, _ := result.Marshal()
	output.Render(string(raw), init.JSON, init.YAML, init.Output, "enrichment-person")
}

// enrichmentCompanyRun the actual work enrichment company
func enrichmentCompanyRun(cmd *cobra.Command, args []string) {
	init := start.New(conn)
	domain := init.Target

	if domain == "" {
		fmt.Println(util.ErrorIcon(), util.Red("domain is required"))
		return
	}
	result, err := init.CompanyFind(tomba.Params{"domain": domain})
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(err.Error()))
		return
	}
	raw, _ := result.Marshal()
	output.Render(string(raw), init.JSON, init.YAML, init.Output, "enrichment-company")
}

// enrichmentCombinedRun the actual work enrichment combined
func enrichmentCombinedRun(cmd *cobra.Command, args []string) {
	init := start.New(conn)
	email := init.Target

	if email == "" {
		fmt.Println(util.ErrorIcon(), util.Red("email is required"))
		return
	}
	result, err := init.CombinedFind(tomba.Params{"email": email})
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(err.Error()))
		return
	}
	raw, _ := result.Marshal()
	output.Render(string(raw), init.JSON, init.YAML, init.Output, "enrichment-combined")
}
