package cmd

import (
	"fmt"

	"github.com/tomba-io/tomba/pkg/start"
	"github.com/tomba-io/tomba/pkg/util"
)

var conn start.Conn
var Long = fmt.Sprintf("cli utility to search or verify lists of email addresses in minutes can significantly improve productivity and efficiency for individuals and businesses dealing with large email databases.\n\n%s", util.RandomBanner())

const (
	authorExample   = `  tomba author --target "https://clearbit.com/blog/company-name-to-domain-api"`
	countExample    = `  tomba count --target "clearbit.com"`
	enrichExample   = `  tomba enrich --target "b.mohamed@tomba.io"`
	linkedinExample = `  tomba linkedin --target "https://www.linkedin.com/in/mohamed-ben-rebia"`
	searchExample   = `  tomba search --target "tomba.io"`
	statusExample   = `  tomba status --target "tomba.io"`
	verifyExample   = `  tomba verify --target "b.mohamed@tomba.io"`
)
