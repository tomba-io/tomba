package cmd

import (
	"fmt"

	"github.com/tomba-io/email/pkg/start"
	"github.com/tomba-io/email/pkg/util"
)

var conn start.Conn
var Long = fmt.Sprintf("cli utility to search or verify lists of email addresses in minutes can significantly improve productivity and efficiency for individuals and businesses dealing with large email databases.\n\n%s", util.RandomBanner())

const (
	authorExample   = `  email author --target "https://clearbit.com/blog/company-name-to-domain-api"`
	countExample    = `  email count --target "clearbit.com"`
	enrichExample   = `  email enrich --target "b.mohamed@tomba.io"`
	linkedinExample = `  email linkedin --target "https://www.linkedin.com/in/alex-maccaw-ab592978"`
	searchExample   = `  email search --target "tomba.io"`
	statusExample   = `  email status --target "tomba.io"`
	verifyExample   = `  email verify --target "b.mohamed@tomba.io"`
)
