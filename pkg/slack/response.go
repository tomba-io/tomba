package slack

import (
	"fmt"

	"github.com/tomba-io/go/tomba/models"
)

// FinderResponse
func FinderResponse(data models.Finder) Model {
	email := data.Data.Email
	fields := make([]Text, 0, 5)

	fields = append(fields, Text{
		Type: "mrkdwn",
		Text: fmt.Sprintf("\n\n*Name*\n\n %s \n\n", data.Data.FullName),
	}, Text{
		Type: "mrkdwn",
		Text: fmt.Sprintf("\n\n*Email*\n\n %s \n\n", email),
	})

	if data.Data.Position != "" {
		fields = append(fields, Text{
			Type: "mrkdwn",
			Text: fmt.Sprintf("\n\n*Bio*\n\n %s at %s \n\n", data.Data.Position, data.Data.Company),
		})
	}

	if data.Data.Country != "" {
		fields = append(fields, Text{
			Type: "mrkdwn",
			Text: fmt.Sprintf("\n\n*Location*\n\n %s \n\n", data.Data.Country),
		})
	}

	if data.Data.Linkedin != "" {
		fields = append(fields, Text{
			Type: "mrkdwn",
			Text: fmt.Sprintf("\n\n*Linkedin*\n\n %s \n\n", data.Data.Linkedin),
		})
	}

	if data.Data.Twitter != nil {
		fields = append(fields, Text{
			Type: "mrkdwn",
			Text: fmt.Sprintf("\n\n*Twitter*\n\n %s \n\n", data.Data.Twitter),
		})
	}

	sources := ""
	if len(data.Data.Sources) > 0 {
		for i := 0; i < len(data.Data.Sources); i++ {
			sources += data.Data.Sources[i].URI + "\n"
		}
	}

	return Model{
		Attachments: []Attachment{
			{
				Color: "#3f77e8",
				Blocks: []Block{
					{
						Type: "section",
						Text: Text{
							Type: "mrkdwn",
							Text: fmt.Sprintf("Email *%s* \n\n<https://app.tomba.io/|Tomba web app>\n\n", email),
						},
						Fields: fields,
					},
					{
						Type: "divider",
					},
					{
						Type: "section",
						Text: Text{
							Type: "mrkdwn",
							Text: fmt.Sprintf("We found `%d` sources for *%s* on the web.\n %s", len(data.Data.Sources), email, sources),
						},
					},
					{
						Type: "divider",
					},
				},
			},
		},
	}
}

// SearchResponse
func SearchResponse(data models.Search) []Model {

	attachment := make([]Attachment, 0, 5)
	emails := data.Data.Emails
	for i := 0; i < len(emails); i++ {
		attachment = append(attachment, Attachment{
			Color:  "",
			Blocks: []Block{},
		})
	}

	return []Model{
		{
			Attachments: []Attachment{},
		},
	}
}

// VerifyResponse
func VerifyResponse(data models.Verifier) Model {
	return Model{}
}
