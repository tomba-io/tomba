package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tomba-io/go/tomba/models"

	"github.com/tomba-io/tomba/pkg/output"
	"github.com/tomba-io/tomba/pkg/start"
	"github.com/tomba-io/tomba/pkg/util"
)

var (
	revealQuery        string
	revealPage         int
	revealCountry      []string
	revealCity         []string
	revealState        []string
	revealIndustry     []string
	revealSize         []string
	revealType         []string
	revealKeywords     []string
	revealFounded      []string
	revealTechnologies []string
	revealSimilar      []string
	revealRevenue      []string
	revealSIC          []string
	revealNAICS        []string
)

// revealCmd represents the reveal command
// see https://docs.tomba.io/api/reveal
var revealCmd = &cobra.Command{
	Use:     "reveal",
	Aliases: []string{"r"},
	Short:   "Search for companies using natural language or structured filters.",
	Long:    Long,
	Run:     revealRun,
	Example: revealExample,
}

func init() {
	revealCmd.Flags().StringVarP(&revealQuery, "query", "q", "", "Natural language query (e.g., 'Real Estate in France').")
	revealCmd.Flags().IntVar(&revealPage, "page", 1, "Page number for pagination.")
	revealCmd.Flags().StringSliceVar(&revealCountry, "country", nil, "Filter by country codes (e.g., US,UK).")
	revealCmd.Flags().StringSliceVar(&revealCity, "city", nil, "Filter by city names.")
	revealCmd.Flags().StringSliceVar(&revealState, "state", nil, "Filter by state/region.")
	revealCmd.Flags().StringSliceVar(&revealIndustry, "industry", nil, "Filter by industry (e.g., Technology).")
	revealCmd.Flags().StringSliceVar(&revealSize, "size", nil, "Filter by company size (e.g., 101-500,501-1000).")
	revealCmd.Flags().StringSliceVar(&revealType, "type", nil, "Filter by company type.")
	revealCmd.Flags().StringSliceVar(&revealKeywords, "keywords", nil, "Filter by keywords.")
	revealCmd.Flags().StringSliceVar(&revealFounded, "founded", nil, "Filter by founding year.")
	revealCmd.Flags().StringSliceVar(&revealTechnologies, "technologies", nil, "Filter by technologies used.")
	revealCmd.Flags().StringSliceVar(&revealSimilar, "similar", nil, "Filter by similar companies.")
	revealCmd.Flags().StringSliceVar(&revealRevenue, "revenue", nil, "Filter by revenue range.")
	revealCmd.Flags().StringSliceVar(&revealSIC, "sic", nil, "Filter by SIC code.")
	revealCmd.Flags().StringSliceVar(&revealNAICS, "naics", nil, "Filter by NAICS code.")
}

// revealRun the actual work reveal
func revealRun(cmd *cobra.Command, args []string) {
	fmt.Println(Long)
	init := start.New(conn)

	// Build request
	request := &models.RevealSearchRequest{
		Page: revealPage,
	}

	if revealQuery != "" {
		request.Query = revealQuery
	}

	// Build filters if any filter flags are provided
	hasFilters := len(revealCountry) > 0 || len(revealCity) > 0 || len(revealState) > 0 ||
		len(revealIndustry) > 0 || len(revealSize) > 0 || len(revealType) > 0 ||
		len(revealKeywords) > 0 || len(revealFounded) > 0 || len(revealTechnologies) > 0 ||
		len(revealSimilar) > 0 || len(revealRevenue) > 0 || len(revealSIC) > 0 || len(revealNAICS) > 0

	if hasFilters {
		request.Filters = &models.RevealSearchFilters{
			Company: &models.RevealCompanyFilters{},
		}

		if len(revealCountry) > 0 {
			request.Filters.Company.LocationCountry = &models.RevealCircularFilter{
				Include: revealCountry,
			}
		}

		if len(revealCity) > 0 {
			request.Filters.Company.LocationCity = &models.RevealCircularFilter{
				Include: revealCity,
			}
		}

		if len(revealState) > 0 {
			request.Filters.Company.LocationState = &models.RevealCircularFilter{
				Include: revealState,
			}
		}

		if len(revealIndustry) > 0 {
			request.Filters.Company.Industry = &models.RevealCircularFilter{
				Include: revealIndustry,
			}
		}

		if len(revealSize) > 0 {
			request.Filters.Company.Size = &models.RevealCircularFilter{
				Include: revealSize,
			}
		}

		if len(revealType) > 0 {
			request.Filters.Company.Type = &models.RevealCircularFilter{
				Include: revealType,
			}
		}

		if len(revealKeywords) > 0 {
			request.Filters.Company.Keywords = &models.RevealCircularFilter{
				Include: revealKeywords,
			}
		}

		if len(revealFounded) > 0 {
			request.Filters.Company.Founded = &models.RevealCircularFilter{
				Include: revealFounded,
			}
		}

		if len(revealTechnologies) > 0 {
			request.Filters.Company.Technologies = &models.RevealCircularFilter{
				Include: revealTechnologies,
			}
		}

		if len(revealSimilar) > 0 {
			request.Filters.Company.Similar = &models.RevealCircularFilter{
				Include: revealSimilar,
			}
		}

		if len(revealRevenue) > 0 {
			request.Filters.Company.Revenue = &models.RevealCircularFilter{
				Include: revealRevenue,
			}
		}

		if len(revealSIC) > 0 {
			request.Filters.Company.SIC = &models.RevealCircularFilter{
				Include: revealSIC,
			}
		}

		if len(revealNAICS) > 0 {
			request.Filters.Company.NAICS = &models.RevealCircularFilter{
				Include: revealNAICS,
			}
		}
	}

	result, err := init.Tomba.SearchCompanies(request)
	if err != nil {
		fmt.Println(util.ErrorIcon(), util.Red(start.ErrErrInvalidLogin.Error()))
		return
	}

	raw, _ := result.Marshal()
	if init.JSON {
		json, _ := output.DisplayJSON(string(raw))
		fmt.Println(json)
	}
	if init.YAML {
		yaml, _ := output.DisplayYAML(string(raw))
		fmt.Println(yaml)
	}
	if init.Output != "" {
		err := output.CreateOutput(init.Output, string(raw))
		if err != nil {
			fmt.Println("Error creating file:", err)
		}
	}
}
