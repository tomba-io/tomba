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
	finderExample      = `  tomba finder --target "tomba.io" --first "mohamed" --last "ben rebia"`
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
  tomba bulk --file prospects.csv --type finder --domain-col company --full-name-col "full_name"
  tomba bulk --file prospects.csv --type finder --domain-col company --first-col first --enrich-mobile
  tomba bulk --file domains.csv --type search --column domain -o results.csv
  tomba bulk --file articles.csv --type author --url-col article_url
  tomba bulk --file profiles.csv --type linkedin --url-col linkedin_url
  tomba bulk --file contacts.csv --type phone --column email
  tomba bulk --file companies.csv --type phone --domain-col domain
  tomba bulk --file profiles.csv --type phone --url-col linkedin_url --full
  tomba bulk --file emails.csv --type sources --column email
  tomba bulk --file domains.csv --type company --domain-col domain
  tomba bulk --file domains.csv --type similar --domain-col domain
  tomba bulk --file phones.csv --type phone-validator --phone-col phone
  tomba bulk --file phones.csv --type phone-validator --phone-col phone --country-code-col country`

	attributeExample = `  tomba attribute list
  tomba attribute create --name "Company Size" --type "text"
  tomba attribute get --id 123
  tomba attribute delete --id 123`

	autocompleteExample = `  tomba autocomplete -t "googl"
  tomba autocomplete -t "tomba" -j`

	enrichmentExample = `  tomba enrichment person -t "user@example.com"
  tomba enrichment company -t "example.com"
  tomba enrichment combined -t "user@example.com"`

	formatExample = `  tomba format -t "tomba.io"
  tomba format -t "stripe.com" -j`

	keyExample = `  tomba key list
  tomba key create
  tomba key delete --id 123
  tomba key reset --id 123`

	leadExample = `  tomba lead list
  tomba lead list --page 2 --limit 20
  tomba lead create --email "user@example.com" --list-id 1
  tomba lead get --id 123
  tomba lead delete --id 123`

	leadsListExample = `  tomba leads-list list
  tomba leads-list create --name "My List"
  tomba leads-list update --id 1 --name "New Name"
  tomba leads-list delete --id 1`

	locationExample = `  tomba location -t "tomba.io"
  tomba location -t "stripe.com" -j`

	loginExample = `  tomba login`

	logoutExample = `  tomba logout`

	logsExample = `  tomba logs
  tomba logs -j`

	usageExample = `  tomba usage
  tomba usage -j`

	updateExample  = `  tomba update`
	versionExample = `  tomba version`

	flagExample = `  tomba flag list
  tomba flag create --email "bounce@example.com" --reason "hard bounce"`
)
