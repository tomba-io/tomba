package cmd

import (
	"fmt"

	"github.com/tomba-io/tomba/pkg/start"
	"github.com/tomba-io/tomba/pkg/util"
)

var conn start.Conn
var Long = fmt.Sprintf("CLI utility to search or verify lists of email addresses in seconds can significantly improve productivity and efficiency for individuals and businesses dealing with large email databases.\n\n%s", util.RandomBanner())

const (
	authorExample      = `  tomba author --target "https://clearbit.com/blog/company-name-to-domain-api"`
	countExample       = `  tomba count --target "clearbit.com"`
	enrichExample      = `  tomba enrich --target "b.mohamed@tomba.io"`
	finderExample      = `  tomba finder --target "tomba.io" --fist "mohamed" --last "ben rebia"`
	linkedinExample    = `  tomba linkedin --target "https://www.linkedin.com/in/mohamed-ben-rebia"`
	phoneFinderExample = `  tomba phone-finder --email "info@stripe.com"
  tomba phone-finder --domain "tomba.io"
  tomba phone-finder --linkedin "https://www.linkedin.com/in/alex-maccaw-ab592978"
  tomba phone-finder --domain "stripe.com" --full`
	phoneValidatorExample = `  tomba phone-validator --phone "+14155552671"
  tomba phone-validator --phone "4155552671" --country-code US`
	revealExample = `  tomba reveal --query "SaaS startups in Europe"
  tomba reveal --country US,UK --industry Technology
  tomba reveal --country US --size 101-500,501-1000 --page 2`
	searchExample     = `  tomba search --target "tomba.io"`
	similarExample    = `  tomba similar --target "tomba.io"`
	statusExample     = `  tomba status --target "tomba.io"`
	technologyExample = `  tomba technology --target "tomba.io"`
	verifyExample     = `  tomba verify --target "b.mohamed@tomba.io"`
	sourcesExample    = `  tomba source --target "b.mohamed@tomba.io"`
	whoamiExample     = `  tomba whoami`
	chatExample       = `  tomba chat
  # Then type: Find the VP Sales at Stripe and get their email
  # Or: Verify the email address b.mohamed@tomba.io

  # Bulk mode: process a CSV file via AI chat
  tomba chat --file contacts.csv --type enrich
  tomba chat --file leads.csv --type verify
  tomba chat --file prospects.csv --type finder`
	skillExample = `  # List all available skills
  tomba skill list

  # Run a skill
  tomba skill run find-email --target stripe.com first_name=David last_name=Singleton
  tomba skill run company-intel --target tomba.io
  tomba skill run enrich-verify --target b.mohamed@tomba.io
  tomba skill run lead-gen --target stripe.com department=sales
  tomba skill run linkedin-intel --target "https://www.linkedin.com/in/someone"

  # Skill details
  tomba skill info company-intel

  # Create a custom skill
  tomba skill create

  # Import / Export
  tomba skill export find-email > my-skill.yaml
  tomba skill import my-skill.yaml

  # Remove a custom skill
  tomba skill remove my-custom-skill`
	bulkExample = `  tomba bulk --file contacts.csv --type enrich
  tomba bulk --file leads.csv --type verify --column "Email Address"
  tomba bulk --file prospects.csv --type finder --domain-col company --first-col first --last-col last
  tomba bulk --file domains.csv --type search --column domain -o results.csv`
)
