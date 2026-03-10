package slack

import (
	"fmt"

	"github.com/tomba-io/go/tomba/models"
)

// FinderResponse
func FinderResponse(data models.Finder) Model {
	// Type assert Data to FinderData (can be FinderData or []FinderData)
	var finderData models.FinderData
	switch v := data.Data.(type) {
	case models.FinderData:
		finderData = v
	case map[string]interface{}:
		// Handle JSON unmarshaled as map
		if email, ok := v["email"].(string); ok {
			finderData.Email = email
		}
		if fullName, ok := v["full_name"].(string); ok {
			finderData.FullName = fullName
		}
		if position, ok := v["position"].(string); ok {
			finderData.Position = position
		}
		if company, ok := v["company"].(string); ok {
			finderData.Company = company
		}
		if country, ok := v["country"].(string); ok {
			finderData.Country = &country
		}
		if linkedin, ok := v["linkedin"].(string); ok {
			finderData.Linkedin = &linkedin
		}
		if twitter, ok := v["twitter"].(string); ok {
			finderData.Twitter = &twitter
		}
		if sources, ok := v["sources"].([]interface{}); ok {
			for _, s := range sources {
				if src, ok := s.(map[string]interface{}); ok {
					if uri, ok := src["uri"].(string); ok {
						finderData.Sources = append(finderData.Sources, models.SourceElement{URI: uri})
					}
				}
			}
		}
	default:
		return Model{}
	}

	email := finderData.Email
	fields := make([]Text, 0, 5)

	fields = append(fields, Text{
		Type: "mrkdwn",
		Text: fmt.Sprintf("\n\n*Name*\n\n %s \n\n", finderData.FullName),
	}, Text{
		Type: "mrkdwn",
		Text: fmt.Sprintf("\n\n*Email*\n\n %s \n\n", email),
	})

	if finderData.Position != "" {
		fields = append(fields, Text{
			Type: "mrkdwn",
			Text: fmt.Sprintf("\n\n*Bio*\n\n %s at %s \n\n", finderData.Position, finderData.Company),
		})
	}

	if finderData.Country != nil && *finderData.Country != "" {
		fields = append(fields, Text{
			Type: "mrkdwn",
			Text: fmt.Sprintf("\n\n*Location*\n\n %s \n\n", *finderData.Country),
		})
	}

	if finderData.Linkedin != nil && *finderData.Linkedin != "" {
		fields = append(fields, Text{
			Type: "mrkdwn",
			Text: fmt.Sprintf("\n\n*Linkedin*\n\n %s \n\n", *finderData.Linkedin),
		})
	}

	if finderData.Twitter != nil && *finderData.Twitter != "" {
		fields = append(fields, Text{
			Type: "mrkdwn",
			Text: fmt.Sprintf("\n\n*Twitter*\n\n %s \n\n", *finderData.Twitter),
		})
	}

	sources := ""
	if len(finderData.Sources) > 0 {
		for i := 0; i < len(finderData.Sources); i++ {
			sources += finderData.Sources[i].URI + "\n"
		}
	}

	return Model{
		Attachments: []Attachment{
			{
				Color: "#3f77e8",
				Blocks: []Block{
					{
						Type: "section",
						Text: &Text{
							Type: "mrkdwn",
							Text: fmt.Sprintf("Email *%s* \n\n<https://app.tomba.io/|Tomba web app>\n\n", email),
						},
						Fields: &fields,
					},
					{
						Type: "divider",
					},
					{
						Type: "section",
						Text: &Text{
							Type: "mrkdwn",
							Text: fmt.Sprintf("We found `%d` sources for *%s* on the web.\n %s", len(finderData.Sources), email, sources),
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
func SearchResponse(data models.Search) Model {

	block := make([]Block, 0, 5)
	emails := data.Data.Emails
	block = append(block, Block{
		Type: "section",
		Text: &Text{
			Type: "mrkdwn",
			Text: fmt.Sprintf("%d results for your search %s\n", data.Meta.Total, *data.Data.Organization.WebsiteURL),
		},
	})

	for i := 0; i < len(emails); i++ {
		fields := make([]Text, 0, 5)
		if emails[i].FullName != nil {
			fields = append(fields, Text{
				Type: "mrkdwn",
				Text: fmt.Sprintf("\n\n*Name*\n\n %s at \n\n", *emails[i].FullName),
			})
		}
		if emails[i].Email != "" {
			fields = append(fields, Text{
				Type: "mrkdwn",
				Text: fmt.Sprintf("\n\n*Email*\n\n %s  \n\n", emails[i].Email),
			})
		}
		if emails[i].Position != nil {
			fields = append(fields, Text{
				Type: "mrkdwn",
				Text: fmt.Sprintf("\n\n*Bio*\n\n %s at %s \n\n", *emails[i].Position, data.Data.Organization.Organization),
			})
		}

		if emails[i].Country != nil {
			fields = append(fields, Text{
				Type: "mrkdwn",
				Text: fmt.Sprintf("\n\n*Location*\n\n %s \n\n", *emails[i].Country),
			})
		}

		if emails[i].Linkedin != nil && *emails[i].Linkedin != "" {
			fields = append(fields, Text{
				Type: "mrkdwn",
				Text: fmt.Sprintf("\n\n*Linkedin*\n\n %s \n\n", *emails[i].Linkedin),
			})
		}

		if emails[i].Twitter != nil && *emails[i].Twitter != "" {
			fields = append(fields, Text{
				Type: "mrkdwn",
				Text: fmt.Sprintf("\n\n*Twitter*\n\n %s \n\n", *emails[i].Twitter),
			})
		}
		sources := ""
		if len(emails[i].Sources) > 0 {
			for i := 0; i < len(emails[i].Sources); i++ {
				sources += emails[i].Sources[i].URI + "\n"
			}
		}
		block = append(block, Block{Type: "divider"})
		block = append(block, Block{
			Type: "section",
			Text: &Text{
				Type: "mrkdwn",
				Text: fmt.Sprintf("Email *%s* \n\n<https://app.tomba.io/|Tomba web app>\n\n", emails[i].Email),
			},
			Fields: &fields,
		})
		block = append(block, Block{
			Type: "section",
			Text: &Text{
				Type: "mrkdwn",
				Text: fmt.Sprintf("We found `%d` sources for *%s* on the web.\n %s", len(emails[i].Sources), emails[i].Email, sources),
			},
		})
		block = append(block, Block{Type: "divider"})
	}

	return Model{Attachments: []Attachment{
		{
			Color:  "#3f77e8",
			Blocks: block,
		}}}
}

// VerifyResponse
func VerifyResponse(data models.Verifier) Model {
	email := data.Data.Email.Email
	fields := make([]Text, 0, 5)

	fields = append(fields, Text{
		Type: "mrkdwn",
		Text: fmt.Sprintf("\n\n*Email*\n\n %s \n\n", email),
	})

	Format := "Invalid"
	if data.Data.Email.Regex && !data.Data.Email.Gibberish {
		Format = "Valid"
	}
	if data.Data.Email.Gibberish {
		Format = "Gibberish"
	}
	fields = append(fields, Text{
		Type: "mrkdwn",
		Text: fmt.Sprintf("\n\n*Format*\n\n %s \n\n", Format),
	})
	ServerStatus := "Invalid"
	if data.Data.Email.MXRecords {
		ServerStatus = "Valid"
	}
	fields = append(fields, Text{
		Type: "mrkdwn",
		Text: fmt.Sprintf("\n\n*Server status*\n\n %s \n\n", ServerStatus),
	})
	Type := "Professional"
	if data.Data.Email.Webmail {
		Type = "Webmail"
	}
	if data.Data.Email.Disposable {
		Type = "Disposable"
	}
	fields = append(fields, Text{
		Type: "mrkdwn",
		Text: fmt.Sprintf("\n\n*Email Type*\n\n %s \n\n", Type),
	})
	EmailStatus := "Invalid"
	if data.Data.Email.Result == "deliverable" {
		EmailStatus = "Valid"
	}
	if data.Data.Email.AcceptAll {
		EmailStatus = "Accept all"
	}
	if data.Data.Email.Block {
		EmailStatus = "Blocked"
	}
	fields = append(fields, Text{
		Type: "mrkdwn",
		Text: fmt.Sprintf("\n\n*Email status*\n\n %s \n\n", EmailStatus),
	})

	if data.Data.Email.Whois.CreatedDate != "" {
		fields = append(fields, Text{
			Type: "mrkdwn",
			Text: fmt.Sprintf("\n\n*Whois Creation Date*\n\n %s \n\n", data.Data.Email.Whois.CreatedDate),
		})
	}
	if data.Data.Email.Whois.ReferralURL != "" {
		fields = append(fields, Text{
			Type: "mrkdwn",
			Text: fmt.Sprintf("\n\n*RWhois*\n\n %s \n\n", data.Data.Email.Whois.ReferralURL),
		})
	}
	if data.Data.Email.Whois.RegistrarName != "" {
		fields = append(fields, Text{
			Type: "mrkdwn",
			Text: fmt.Sprintf("\n\n*Whois Name*\n\n %s \n\n", data.Data.Email.Whois.RegistrarName),
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
						Text: &Text{
							Type: "mrkdwn",
							Text: fmt.Sprintf("Email *%s* \n\n<https://app.tomba.io/|Tomba web app>\n\n", email),
						},
						Fields: &fields,
					},
					{
						Type: "divider",
					},
					{
						Type: "section",
						Text: &Text{
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
