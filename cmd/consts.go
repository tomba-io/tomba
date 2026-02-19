package cmd

import (
	"fmt"

	"github.com/tomba-io/tomba/pkg/start"
	"github.com/tomba-io/tomba/pkg/util"
)

var conn start.Conn
var Long = fmt.Sprintf("cli utility to search or verify lists of email addresses in minutes can significantly improve productivity and efficiency for individuals and businesses dealing with large email databases.\n\n%s", util.RandomBanner())

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
	revealExample = `  tomba reveal --query "Real Estate in France"
  tomba reveal --country US,UK --industry Technology
  tomba reveal --country US --size 101-500,501-1000 --page 2`
	searchExample  = `  tomba search --target "tomba.io"`
	statusExample  = `  tomba status --target "tomba.io"`
	verifyExample  = `  tomba verify --target "b.mohamed@tomba.io"`
	sourcesExample = `  tomba source --target "b.mohamed@tomba.io"`
	whoamiExample  = `  tomba whoami`
)
