package output

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tomba-io/tomba/pkg/util"
)

// DisplayText formats JSON data as human-readable text output
func DisplayText(jsonString string, command string) string {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonString), &data); err != nil {
		return jsonString
	}

	switch command {
	case "search":
		return formatSearch(data)
	case "finder":
		return formatFinder(data)
	case "enrich":
		return formatEnrich(data)
	case "verify":
		return formatVerify(data)
	case "count":
		return formatCount(data)
	case "status":
		return formatStatus(data)
	case "sources":
		return formatSources(data)
	case "similar":
		return formatSimilar(data)
	case "technology":
		return formatTechnology(data)
	case "whoami":
		return formatWhoami(data)
	case "usage":
		return formatUsage(data)
	case "linkedin":
		return formatFinder(data)
	case "author":
		return formatFinder(data)
	case "reveal":
		return formatReveal(data)
	case "phone-finder":
		return formatPhoneFinder(data)
	case "phone-validator":
		return formatPhoneValidator(data)
	default:
		return formatGeneric(data, 0)
	}
}

func formatSearch(data map[string]interface{}) string {
	var sb strings.Builder

	// Meta info
	if meta, ok := data["meta"].(map[string]interface{}); ok {
		total := getFloat(meta, "total")
		page := getFloat(meta, "page")
		limit := getFloat(meta, "limit")
		sb.WriteString(fmt.Sprintf("\n%s Domain Search Results\n", util.SuccessIcon()))
		sb.WriteString(fmt.Sprintf("  Total:   %s\n", util.Bold(fmt.Sprintf("%.0f", total))))
		if page > 0 {
			sb.WriteString(fmt.Sprintf("  Page:    %.0f\n", page))
		}
		if limit > 0 {
			sb.WriteString(fmt.Sprintf("  Limit:   %.0f\n", limit))
		}
	}

	if d, ok := data["data"].(map[string]interface{}); ok {
		// Organization details
		if org, ok := d["organization"].(map[string]interface{}); ok {
			sb.WriteString("\n  Organization:\n")
			writeField(&sb, "    ", "Name", getStr(org, "organization"))
			writeField(&sb, "    ", "Website", getStr(org, "website_url"))
			writeField(&sb, "    ", "Blog", getStr(org, "blog_url"))
			writeField(&sb, "    ", "Description", getStr(org, "description"))
			writeField(&sb, "    ", "Industry", getStr(org, "industry"))
			writeField(&sb, "    ", "Type", getStr(org, "type"))
			writeField(&sb, "    ", "Country", getStr(org, "country"))
			writeField(&sb, "    ", "State", getStr(org, "state"))
			writeField(&sb, "    ", "City", getStr(org, "city"))
			writeField(&sb, "    ", "Postal Code", getStr(org, "postal_code"))
			writeField(&sb, "    ", "Street", getStr(org, "street_address"))
			writeField(&sb, "    ", "Size", getStr(org, "size"))
			writeField(&sb, "    ", "Founded", getStr(org, "founded"))
			writeField(&sb, "    ", "Revenue", getStr(org, "revenue"))
			writeField(&sb, "    ", "LinkedIn", getStr(org, "linkedin"))
			writeField(&sb, "    ", "Twitter", getStr(org, "twitter"))
			writeField(&sb, "    ", "Facebook", getStr(org, "facebook"))
			writeField(&sb, "    ", "Instagram", getStr(org, "instagram"))
			writeField(&sb, "    ", "Phone", getStr(org, "phone"))
			writeField(&sb, "    ", "Fax", getStr(org, "fax"))

			// Technologies
			if techs, ok := org["technologies"].([]interface{}); ok && len(techs) > 0 {
				var techNames []string
				for _, t := range techs {
					if s, ok := t.(string); ok {
						techNames = append(techNames, s)
					}
				}
				if len(techNames) > 0 {
					sb.WriteString(fmt.Sprintf("    Technologies: %s\n", util.Gray(strings.Join(techNames, ", "))))
				}
			}

			// SIC/NAICS
			if sic, ok := org["sic"].([]interface{}); ok && len(sic) > 0 {
				var codes []string
				for _, s := range sic {
					if c, ok := s.(string); ok {
						codes = append(codes, c)
					}
				}
				if len(codes) > 0 {
					sb.WriteString(fmt.Sprintf("    SIC Codes: %s\n", strings.Join(codes, ", ")))
				}
			}
			if naics, ok := org["naics"].([]interface{}); ok && len(naics) > 0 {
				var codes []string
				for _, n := range naics {
					if c, ok := n.(string); ok {
						codes = append(codes, c)
					}
				}
				if len(codes) > 0 {
					sb.WriteString(fmt.Sprintf("    NAICS Codes: %s\n", strings.Join(codes, ", ")))
				}
			}
		}

		// Pattern info
		writeField(&sb, "  ", "Email Pattern", getStr(d, "pattern"))
		if accept, ok := d["accept_all"].(bool); ok {
			if accept {
				sb.WriteString(fmt.Sprintf("  Accept All: %s\n", util.Yellow("yes (catch-all)")))
			} else {
				sb.WriteString(fmt.Sprintf("  Accept All: %s\n", util.Green("no")))
			}
		}

		// Email list
		if emails, ok := d["emails"].([]interface{}); ok && len(emails) > 0 {
			sb.WriteString(fmt.Sprintf("\n  Emails (%d):\n", len(emails)))
			sb.WriteString(fmt.Sprintf("  %-4s %-35s %-22s %-18s %-12s %-8s %s\n", "#", "Email", "Name", "Position", "Department", "Type", "Score"))
			sb.WriteString(fmt.Sprintf("  %s\n", strings.Repeat("─", 120)))
			for i, e := range emails {
				if em, ok := e.(map[string]interface{}); ok {
					email := getStr(em, "email")
					first := getStr(em, "first_name")
					last := getStr(em, "last_name")
					position := getStr(em, "position")
					dept := getStr(em, "department")
					emailType := getStr(em, "type")
					name := strings.TrimSpace(first + " " + last)
					score := ""
					if s, ok := em["score"].(float64); ok {
						score = fmt.Sprintf("%.0f%%", s)
					}
					country := getStr(em, "country")
					linkedin := getStr(em, "linkedin")
					twitter := getStr(em, "twitter")

					sb.WriteString(fmt.Sprintf("  %-4d %-35s %-22s %-18s %-12s %-8s %s\n",
						i+1, util.Green(email), name, util.Gray(position), dept, emailType, score))

					// Extra details per email
					if country != "" {
						sb.WriteString(fmt.Sprintf("       Country: %s", country))
						if linkedin != "" {
							sb.WriteString(fmt.Sprintf("  LinkedIn: %s", linkedin))
						}
						if twitter != "" {
							sb.WriteString(fmt.Sprintf("  Twitter: %s", twitter))
						}
						sb.WriteString("\n")
					} else if linkedin != "" || twitter != "" {
						sb.WriteString("      ")
						if linkedin != "" {
							sb.WriteString(fmt.Sprintf(" LinkedIn: %s", linkedin))
						}
						if twitter != "" {
							sb.WriteString(fmt.Sprintf("  Twitter: %s", twitter))
						}
						sb.WriteString("\n")
					}

					// Phone number if available
					if phone := getStr(em, "phone_number"); phone != "" {
						sb.WriteString(fmt.Sprintf("       Phone: %s\n", phone))
					}

					// Sources
					if sources, ok := em["sources"].([]interface{}); ok && len(sources) > 0 {
						sb.WriteString(fmt.Sprintf("       Sources: %d found\n", len(sources)))
					}

					// Verification
					if verif, ok := em["verification"].(map[string]interface{}); ok {
						vDate := getStr(verif, "date")
						vStatus := getStr(verif, "status")
						if vStatus != "" {
							sb.WriteString(fmt.Sprintf("       Verification: %s", vStatus))
							if vDate != "" {
								sb.WriteString(fmt.Sprintf(" (%s)", vDate))
							}
							sb.WriteString("\n")
						}
					}

					// Last updated
					if lastUpdated := getStr(em, "last_updated"); lastUpdated != "" {
						sb.WriteString(fmt.Sprintf("       Last Updated: %s\n", util.Gray(lastUpdated)))
					}
				}
			}
		}
	}

	return sb.String()
}

func formatFinder(data map[string]interface{}) string {
	var sb strings.Builder

	if d, ok := data["data"].(map[string]interface{}); ok {
		email := getStr(d, "email")
		if email != "" {
			sb.WriteString(fmt.Sprintf("\n%s Found: %s\n", util.SuccessIcon(), util.Green(email)))
		}

		sb.WriteString("\n  Contact Details:\n")
		writeField(&sb, "    ", "Email", email)
		first := getStr(d, "first_name")
		last := getStr(d, "last_name")
		if first != "" || last != "" {
			sb.WriteString(fmt.Sprintf("    Name:         %s\n", util.Bold(strings.TrimSpace(first+" "+last))))
		}
		writeField(&sb, "    ", "First Name", first)
		writeField(&sb, "    ", "Last Name", last)
		writeField(&sb, "    ", "Gender", getStr(d, "gender"))
		writeField(&sb, "    ", "Position", getStr(d, "position"))
		writeField(&sb, "    ", "Department", getStr(d, "department"))
		writeField(&sb, "    ", "Seniority", getStr(d, "seniority"))
		writeField(&sb, "    ", "Company", getStr(d, "company"))
		writeField(&sb, "    ", "Website", getStr(d, "website_url"))
		writeField(&sb, "    ", "Country", getStr(d, "country"))
		writeField(&sb, "    ", "City", getStr(d, "city"))
		writeField(&sb, "    ", "State", getStr(d, "state"))

		// Score
		if score, ok := d["score"].(float64); ok {
			sb.WriteString(fmt.Sprintf("    Score:        %s\n", util.Bold(fmt.Sprintf("%.0f%%", score))))
		}

		// Type
		writeField(&sb, "    ", "Type", getStr(d, "type"))

		// Social profiles
		sb.WriteString("\n  Social Profiles:\n")
		writeField(&sb, "    ", "LinkedIn", getStr(d, "linkedin"))
		writeField(&sb, "    ", "Twitter", getStr(d, "twitter"))
		writeField(&sb, "    ", "Facebook", getStr(d, "facebook"))
		writeField(&sb, "    ", "GitHub", getStr(d, "github"))

		// Phone
		if phone := getStr(d, "phone_number"); phone != "" {
			sb.WriteString(fmt.Sprintf("\n  Phone: %s\n", util.Bold(phone)))
		}

		// Accept all
		if accept, ok := d["accept_all"].(bool); ok {
			if accept {
				sb.WriteString(fmt.Sprintf("    Accept All:   %s\n", util.Yellow("yes (catch-all domain)")))
			} else {
				sb.WriteString(fmt.Sprintf("    Accept All:   %s\n", util.Green("no")))
			}
		}

		// Pattern
		writeField(&sb, "    ", "Pattern", getStr(d, "pattern"))

		// Verification
		if verif, ok := d["verification"].(map[string]interface{}); ok {
			sb.WriteString("\n  Verification:\n")
			writeField(&sb, "    ", "Date", getStr(verif, "date"))
			writeField(&sb, "    ", "Status", getStr(verif, "status"))
		}

		// Sources
		if sources, ok := d["sources"].([]interface{}); ok && len(sources) > 0 {
			sb.WriteString(fmt.Sprintf("\n  Sources (%d):\n", len(sources)))
			for i, s := range sources {
				if src, ok := s.(map[string]interface{}); ok {
					uri := getStr(src, "uri")
					extracted := getStr(src, "extracted_on")
					sb.WriteString(fmt.Sprintf("    %d. %s", i+1, util.Cyan(uri)))
					if extracted != "" {
						sb.WriteString(fmt.Sprintf(" (%s)", util.Gray(extracted)))
					}
					sb.WriteString("\n")
				}
			}
		}

		// Last updated
		writeField(&sb, "    ", "Last Updated", getStr(d, "last_updated"))
	}

	return sb.String()
}

func formatEnrich(data map[string]interface{}) string {
	var sb strings.Builder

	if d, ok := data["data"].(map[string]interface{}); ok {
		email := getStr(d, "email")
		if email != "" {
			sb.WriteString(fmt.Sprintf("\n%s Enriched: %s\n", util.SuccessIcon(), util.Green(email)))
		}

		sb.WriteString("\n  Contact Information:\n")
		writeField(&sb, "    ", "Email", email)
		first := getStr(d, "first_name")
		last := getStr(d, "last_name")
		if first != "" || last != "" {
			sb.WriteString(fmt.Sprintf("    Full Name:    %s\n", util.Bold(strings.TrimSpace(first+" "+last))))
		}
		writeField(&sb, "    ", "First Name", first)
		writeField(&sb, "    ", "Last Name", last)
		writeField(&sb, "    ", "Gender", getStr(d, "gender"))

		sb.WriteString("\n  Professional:\n")
		writeField(&sb, "    ", "Position", getStr(d, "position"))
		writeField(&sb, "    ", "Seniority", getStr(d, "seniority"))
		writeField(&sb, "    ", "Department", getStr(d, "department"))
		writeField(&sb, "    ", "Company", getStr(d, "company"))
		writeField(&sb, "    ", "Website", getStr(d, "website_url"))
		writeField(&sb, "    ", "Industry", getStr(d, "industry"))

		sb.WriteString("\n  Location:\n")
		writeField(&sb, "    ", "Country", getStr(d, "country"))
		writeField(&sb, "    ", "State", getStr(d, "state"))
		writeField(&sb, "    ", "City", getStr(d, "city"))

		sb.WriteString("\n  Social Profiles:\n")
		writeField(&sb, "    ", "LinkedIn", getStr(d, "linkedin"))
		writeField(&sb, "    ", "Twitter", getStr(d, "twitter"))
		writeField(&sb, "    ", "Facebook", getStr(d, "facebook"))
		writeField(&sb, "    ", "GitHub", getStr(d, "github"))

		// Phone
		if phone := getStr(d, "phone_number"); phone != "" {
			sb.WriteString(fmt.Sprintf("\n  Phone: %s\n", util.Bold(phone)))
		}

		// Score
		if score, ok := d["score"].(float64); ok {
			sb.WriteString(fmt.Sprintf("\n  Score:          %s\n", util.Bold(fmt.Sprintf("%.0f%%", score))))
		}

		writeField(&sb, "  ", "Type", getStr(d, "type"))
		writeField(&sb, "  ", "Pattern", getStr(d, "pattern"))

		// Accept all
		if accept, ok := d["accept_all"].(bool); ok {
			if accept {
				sb.WriteString(fmt.Sprintf("  Accept All:     %s\n", util.Yellow("yes")))
			} else {
				sb.WriteString(fmt.Sprintf("  Accept All:     %s\n", util.Green("no")))
			}
		}

		// Verification
		if verif, ok := d["verification"].(map[string]interface{}); ok {
			sb.WriteString("\n  Verification:\n")
			writeField(&sb, "    ", "Date", getStr(verif, "date"))
			writeField(&sb, "    ", "Status", getStr(verif, "status"))
		}

		// Sources
		if sources, ok := d["sources"].([]interface{}); ok && len(sources) > 0 {
			sb.WriteString(fmt.Sprintf("\n  Sources (%d):\n", len(sources)))
			for i, s := range sources {
				if src, ok := s.(map[string]interface{}); ok {
					uri := getStr(src, "uri")
					extracted := getStr(src, "extracted_on")
					sb.WriteString(fmt.Sprintf("    %d. %s", i+1, util.Cyan(uri)))
					if extracted != "" {
						sb.WriteString(fmt.Sprintf(" (%s)", util.Gray(extracted)))
					}
					sb.WriteString("\n")
				}
			}
		}

		writeField(&sb, "  ", "Last Updated", getStr(d, "last_updated"))
	}

	return sb.String()
}

func formatVerify(data map[string]interface{}) string {
	var sb strings.Builder

	if d, ok := data["data"].(map[string]interface{}); ok {
		if emailData, ok := d["email"].(map[string]interface{}); ok {
			email := getStr(emailData, "email")
			result := getStr(emailData, "result")

			icon := util.SuccessIcon()
			statusColor := util.Green(result)
			if result == "undeliverable" {
				icon = util.ErrorIcon()
				statusColor = util.Red(result)
			} else if result == "risky" {
				icon = util.WarningIcon()
				statusColor = util.Yellow(result)
			}

			sb.WriteString(fmt.Sprintf("\n%s %s is %s\n", icon, util.Bold(email), statusColor))

			sb.WriteString("\n  Email Details:\n")
			writeField(&sb, "    ", "Email", email)
			writeField(&sb, "    ", "Result", result)
			writeField(&sb, "    ", "Status", getStr(emailData, "status"))

			if score, ok := emailData["score"].(float64); ok {
				sb.WriteString(fmt.Sprintf("    Score:        %.0f%%\n", score))
			}

			writeField(&sb, "    ", "Regex", formatBoolStr(emailData, "regex"))
			writeField(&sb, "    ", "Gibberish", formatBoolStr(emailData, "gibberish"))
			writeField(&sb, "    ", "Disposable", formatBoolStr(emailData, "disposable"))
			writeField(&sb, "    ", "Webmail", formatBoolStr(emailData, "webmail"))

			sb.WriteString("\n  Server Checks:\n")
			writeBoolField(&sb, "    ", "MX Records", emailData, "mx_found")
			writeBoolField(&sb, "    ", "SMTP Server", emailData, "smtp_server")
			writeBoolField(&sb, "    ", "SMTP Check", emailData, "smtp_check")
			writeBoolField(&sb, "    ", "Accept All", emailData, "accept_all")
			writeBoolField(&sb, "    ", "Block", emailData, "block")

			sb.WriteString("\n  DNS:\n")
			writeField(&sb, "    ", "MX Host", getStr(emailData, "mx_host"))
			writeField(&sb, "    ", "MX IP", getStr(emailData, "mx_ip"))
			writeField(&sb, "    ", "MX Priority", getStr(emailData, "mx_priority"))

			// Sources
			if sources, ok := emailData["sources"].([]interface{}); ok && len(sources) > 0 {
				sb.WriteString(fmt.Sprintf("\n  Sources (%d):\n", len(sources)))
				for i, s := range sources {
					if src, ok := s.(map[string]interface{}); ok {
						uri := getStr(src, "uri")
						extracted := getStr(src, "extracted_on")
						sb.WriteString(fmt.Sprintf("    %d. %s", i+1, util.Cyan(uri)))
						if extracted != "" {
							sb.WriteString(fmt.Sprintf(" (%s)", util.Gray(extracted)))
						}
						sb.WriteString("\n")
					}
				}
			}

			writeField(&sb, "    ", "Whois Created", getStr(emailData, "whois_created_date"))
			writeField(&sb, "    ", "Whois Referral", getStr(emailData, "whois_referral_url"))
		}
	}

	return sb.String()
}

func formatCount(data map[string]interface{}) string {
	var sb strings.Builder

	if d, ok := data["data"].(map[string]interface{}); ok {
		if total, ok := d["total"].(float64); ok {
			sb.WriteString(fmt.Sprintf("\n%s Email Count\n", util.SuccessIcon()))
			sb.WriteString(fmt.Sprintf("  Total:     %s\n", util.Bold(fmt.Sprintf("%.0f", total))))
		}
		if personal, ok := d["personal_emails"].(float64); ok {
			sb.WriteString(fmt.Sprintf("  Personal:  %.0f\n", personal))
		}
		if generic, ok := d["generic_emails"].(float64); ok {
			sb.WriteString(fmt.Sprintf("  Generic:   %.0f\n", generic))
		}

		if dept, ok := d["department"].(map[string]interface{}); ok {
			sb.WriteString("\n  Departments:\n")
			for k, v := range dept {
				if count, ok := v.(float64); ok {
					bar := ""
					if total, ok := d["total"].(float64); ok && total > 0 {
						pct := count / total * 100
						filled := int(pct / 5)
						if filled > 20 {
							filled = 20
						}
						bar = fmt.Sprintf(" %s %.0f%%", strings.Repeat("█", filled)+strings.Repeat("░", 20-filled), pct)
					}
					sb.WriteString(fmt.Sprintf("    %-20s %5.0f%s\n", k, count, util.Gray(bar)))
				}
			}
		}
	}

	return sb.String()
}

func formatStatus(data map[string]interface{}) string {
	var sb strings.Builder

	// API returns flat: {"domain":"...","webmail":true,"disposable":false}
	// or wrapped: {"data":{"domain":"...",...}}
	d := data

	domain := getStr(d, "domain")
	sb.WriteString(fmt.Sprintf("\n  Domain Status: %s\n\n", util.Bold(domain)))

	writeBoolField(&sb, "  ", "Webmail", d, "webmail")
	writeBoolField(&sb, "  ", "Disposable", d, "disposable")

	return sb.String()
}

func formatSources(data map[string]interface{}) string {
	var sb strings.Builder

	if d, ok := data["data"].([]interface{}); ok {
		sb.WriteString(fmt.Sprintf("\n%s Found %d sources\n\n", util.SuccessIcon(), len(d)))
		sb.WriteString(fmt.Sprintf("  %-4s %-60s %-20s %s\n", "#", "URL", "Extracted On", "Still On Page"))
		sb.WriteString(fmt.Sprintf("  %s\n", strings.Repeat("─", 100)))
		for i, s := range d {
			if src, ok := s.(map[string]interface{}); ok {
				uri := getStr(src, "uri")
				extracted := getStr(src, "extracted_on")
				lastSeen := getStr(src, "last_seen_on")
				stillOnPage := formatBoolStr(src, "still_on_page")

				// Truncate long URIs
				displayURI := uri
				if len(displayURI) > 58 {
					displayURI = displayURI[:55] + "..."
				}

				sb.WriteString(fmt.Sprintf("  %-4d %-60s %-20s %s\n", i+1, util.Cyan(displayURI), util.Gray(extracted), stillOnPage))
				if lastSeen != "" {
					sb.WriteString(fmt.Sprintf("       Last Seen: %s\n", lastSeen))
				}
				if website := getStr(src, "website_url"); website != "" {
					sb.WriteString(fmt.Sprintf("       Website: %s\n", website))
				}
			}
		}
	}

	return sb.String()
}

func formatSimilar(data map[string]interface{}) string {
	var sb strings.Builder

	// Meta: total, pageSize, current, total_pages
	if meta, ok := data["meta"].(map[string]interface{}); ok {
		total := getFloat(meta, "total")
		current := getFloat(meta, "current")
		totalPages := getFloat(meta, "total_pages")
		sb.WriteString(fmt.Sprintf("\n%s Found %.0f similar domains", util.SuccessIcon(), total))
		if totalPages > 0 {
			sb.WriteString(fmt.Sprintf("  (page %.0f/%.0f)", current, totalPages))
		}
		sb.WriteString("\n")
	}

	if d, ok := data["data"].([]interface{}); ok {
		if _, hasMeta := data["meta"]; !hasMeta {
			sb.WriteString(fmt.Sprintf("\n%s Found %d similar domains\n", util.SuccessIcon(), len(d)))
		}

		sb.WriteString(fmt.Sprintf("\n  %-4s %-30s %-30s %s\n", "#", "Domain", "Name", "Industry"))
		sb.WriteString(fmt.Sprintf("  %s\n", strings.Repeat("─", 100)))
		for i, s := range d {
			if sim, ok := s.(map[string]interface{}); ok {
				// website_url or host
				domain := getStr(sim, "website_url")
				if domain == "" {
					domain = getStr(sim, "host")
				}
				name := getStr(sim, "name")
				// industries or country
				industry := getStr(sim, "industries")
				country := getStr(sim, "country")

				extra := industry
				if extra == "" {
					extra = country
				}

				nameDisplay := name
				if nameDisplay == "" {
					nameDisplay = util.Gray("-")
				}

				sb.WriteString(fmt.Sprintf("  %-4d %-30s %-30s %s\n",
					i+1, util.Green(domain), nameDisplay, util.Gray(extra)))
			}
		}
	}

	return sb.String()
}

func formatTechnology(data map[string]interface{}) string {
	var sb strings.Builder

	// Domain field at top level or inside data
	domain := getStr(data, "domain")

	// Technologies: data[] (flat array) or data.technologies[]
	var techs []interface{}
	if d, ok := data["data"].([]interface{}); ok {
		techs = d
	} else if d, ok := data["data"].(map[string]interface{}); ok {
		if domain == "" {
			domain = getStr(d, "domain")
		}
		if t, ok := d["technologies"].([]interface{}); ok {
			techs = t
		}
	}

	if domain != "" {
		writeField(&sb, "\n  ", "Domain", domain)
	}

	if len(techs) > 0 {
		sb.WriteString(fmt.Sprintf("\n  Technologies (%d):\n\n", len(techs)))
		sb.WriteString(fmt.Sprintf("  %-4s %-22s %-35s %s\n", "#", "Name", "Categories", "Website"))
		sb.WriteString(fmt.Sprintf("  %s\n", strings.Repeat("─", 100)))
		for i, t := range techs {
			if tech, ok := t.(map[string]interface{}); ok {
				name := getStr(tech, "name")
				website := getStr(tech, "website")

				// categories: string or []interface{} of {name, slug}
				var cats string
				if catStr := getStr(tech, "categories"); catStr != "" {
					cats = catStr
				} else if catArr, ok := tech["categories"].([]interface{}); ok {
					var catNames []string
					for _, c := range catArr {
						if cm, ok := c.(map[string]interface{}); ok {
							if cn := getStr(cm, "name"); cn != "" {
								catNames = append(catNames, cn)
							}
						}
					}
					cats = strings.Join(catNames, ", ")
				}

				sb.WriteString(fmt.Sprintf("  %-4d %-22s %-35s %s\n",
					i+1, util.Bold(name), util.Gray(cats), util.Cyan(website)))

				// Description
				if desc := getStr(tech, "description"); desc != "" {
					displayDesc := desc
					if len(displayDesc) > 95 {
						displayDesc = displayDesc[:92] + "..."
					}
					sb.WriteString(fmt.Sprintf("       %s\n", util.Gray(displayDesc)))
				}
			}
		}
	}

	return sb.String()
}

func formatWhoami(data map[string]interface{}) string {
	var sb strings.Builder

	if d, ok := data["data"].(map[string]interface{}); ok {
		email := getStr(d, "email")
		sb.WriteString(fmt.Sprintf("\n%s Authenticated as %s\n", util.SuccessIcon(), util.Green(email)))

		sb.WriteString("\n  Account Details:\n")
		writeField(&sb, "    ", "Email", email)
		first := getStr(d, "first_name")
		last := getStr(d, "last_name")
		if first != "" || last != "" {
			sb.WriteString(fmt.Sprintf("    Name:         %s\n", util.Bold(strings.TrimSpace(first+" "+last))))
		}
		writeField(&sb, "    ", "First Name", first)
		writeField(&sb, "    ", "Last Name", last)
		writeField(&sb, "    ", "Phone", getStr(d, "phone"))
		writeField(&sb, "    ", "Country", getStr(d, "country"))
		writeField(&sb, "    ", "Job Title", getStr(d, "job_title"))
		writeField(&sb, "    ", "Company", getStr(d, "company"))
		writeField(&sb, "    ", "Website", getStr(d, "website"))
		writeField(&sb, "    ", "Image", getStr(d, "image"))

		if id, ok := d["user_id"].(float64); ok {
			sb.WriteString(fmt.Sprintf("    User ID:      %.0f\n", id))
		}

		// IP info
		writeField(&sb, "    ", "IP", getStr(d, "ip"))

		// Pricing / Plan
		if plan, ok := d["pricing"].(map[string]interface{}); ok {
			sb.WriteString("\n  Plan:\n")
			writeField(&sb, "    ", "Name", getStr(plan, "name"))
			if searches, ok := plan["search"].(map[string]interface{}); ok {
				used := getFloat(searches, "used")
				limit := getFloat(searches, "limit")
				sb.WriteString(fmt.Sprintf("    Search:       %.0f/%.0f\n", used, limit))
			}
			if verify, ok := plan["verification"].(map[string]interface{}); ok {
				used := getFloat(verify, "used")
				limit := getFloat(verify, "limit")
				sb.WriteString(fmt.Sprintf("    Verification: %.0f/%.0f\n", used, limit))
			}
		}

		writeField(&sb, "    ", "Created At", getStr(d, "created_at"))
	}

	return sb.String()
}

func formatUsage(data map[string]interface{}) string {
	var sb strings.Builder

	if d, ok := data["data"].(map[string]interface{}); ok {
		sb.WriteString(fmt.Sprintf("\n%s API Usage\n", util.SuccessIcon()))

		if usage, ok := d["usage"].(map[string]interface{}); ok {
			sb.WriteString(fmt.Sprintf("\n  %-25s %-10s %-10s %s\n", "Endpoint", "Used", "Limit", "Progress"))
			sb.WriteString(fmt.Sprintf("  %s\n", strings.Repeat("─", 70)))
			for key, val := range usage {
				if v, ok := val.(map[string]interface{}); ok {
					used := getFloat(v, "used")
					limit := getFloat(v, "limit")
					pct := float64(0)
					if limit > 0 {
						pct = used / limit * 100
					}
					filled := int(pct / 5)
					if filled > 20 {
						filled = 20
					}
					bar := strings.Repeat("█", filled) + strings.Repeat("░", 20-filled)
					color := util.Green
					if pct > 80 {
						color = util.Red
					} else if pct > 50 {
						color = util.Yellow
					}
					sb.WriteString(fmt.Sprintf("  %-25s %-10.0f %-10.0f %s %.0f%%\n",
						util.Bold(key), used, limit, color(bar), pct))
				}
			}
		}
	}

	return sb.String()
}

func formatReveal(data map[string]interface{}) string {
	var sb strings.Builder

	// Pagination from meta or data.data
	var total, page, limit, pages float64
	if meta, ok := data["meta"].(map[string]interface{}); ok {
		total = getFloat(meta, "total")
		page = getFloat(meta, "page")
		limit = getFloat(meta, "limit")
		pages = getFloat(meta, "pages")
	}
	// data.data may also contain total/page/limit/pages
	d, _ := data["data"].(map[string]interface{})
	if d != nil {
		if total == 0 {
			total = getFloat(d, "total")
		}
		if page == 0 {
			page = getFloat(d, "page")
		}
		if limit == 0 {
			limit = getFloat(d, "limit")
		}
		if pages == 0 {
			pages = getFloat(d, "pages")
		}
	}

	sb.WriteString(fmt.Sprintf("\n%s Found %.0f companies", util.SuccessIcon(), total))
	if pages > 0 {
		sb.WriteString(fmt.Sprintf("  (page %.0f/%.0f, %.0f per page)", page, pages, limit))
	}
	sb.WriteString("\n")

	// Active filters from meta
	if meta, ok := data["meta"].(map[string]interface{}); ok {
		if filters, ok := meta["filters"].(map[string]interface{}); ok {
			var activeFilters []string
			for filterName, filterVal := range filters {
				if fv, ok := filterVal.(map[string]interface{}); ok {
					if inc, ok := fv["include"].([]interface{}); ok && len(inc) > 0 {
						var vals []string
						for _, v := range inc {
							if s, ok := v.(string); ok {
								vals = append(vals, s)
							}
						}
						if len(vals) > 0 {
							activeFilters = append(activeFilters, fmt.Sprintf("%s: %s", filterName, strings.Join(vals, ", ")))
						}
					}
				}
			}
			if len(activeFilters) > 0 {
				sb.WriteString(fmt.Sprintf("  Filters: %s\n", util.Gray(strings.Join(activeFilters, "  |  "))))
			}
		}
	}

	// Companies list: data.data.companies[] or data.data[] (older format)
	var companies []interface{}
	if d != nil {
		if c, ok := d["companies"].([]interface{}); ok {
			companies = c
		}
	}
	if companies == nil {
		if c, ok := data["data"].([]interface{}); ok {
			companies = c
		}
	}

	if len(companies) > 0 {
		sb.WriteString(fmt.Sprintf("\n  %-4s %-28s %-18s %-15s %-12s %-10s %s\n",
			"#", "Company", "Industry", "Country", "Size", "Founded", "Website"))
		sb.WriteString(fmt.Sprintf("  %s\n", strings.Repeat("─", 120)))

		for i, c := range companies {
			company, ok := c.(map[string]interface{})
			if !ok {
				continue
			}

			// name or organization
			name := getStr(company, "name")
			if name == "" {
				name = getStr(company, "organization")
			}
			industry := getStr(company, "industry")
			country := getStr(company, "country")
			// company_size or size
			size := getStr(company, "company_size")
			if size == "" {
				size = getStr(company, "size")
			}
			founded := getStr(company, "founded")
			// website_url or domain
			website := getStr(company, "website_url")
			if website == "" {
				website = getStr(company, "domain")
			}

			sb.WriteString(fmt.Sprintf("  %-4d %-28s %-18s %-15s %-12s %-10s %s\n",
				i+1, util.Bold(name), industry, country, size, founded, util.Cyan(website)))

			// Description
			if desc := getStr(company, "description"); desc != "" {
				displayDesc := desc
				if len(displayDesc) > 100 {
					displayDesc = displayDesc[:97] + "..."
				}
				sb.WriteString(fmt.Sprintf("       %s\n", util.Gray(displayDesc)))
			}

			// Location
			state := getStr(company, "state")
			city := getStr(company, "city")
			street := getStr(company, "street_address")
			postal := getStr(company, "postal_code")
			var location []string
			if street != "" {
				location = append(location, street)
			}
			if city != "" {
				location = append(location, city)
			}
			if state != "" {
				location = append(location, state)
			}
			if postal != "" {
				location = append(location, postal)
			}
			if len(location) > 0 {
				sb.WriteString(fmt.Sprintf("       Location: %s\n", strings.Join(location, ", ")))
			}

			// Extra details
			compType := getStr(company, "type")
			revenue := getStr(company, "revenue")
			var details []string
			if compType != "" {
				details = append(details, fmt.Sprintf("Type: %s", compType))
			}
			if revenue != "" {
				details = append(details, fmt.Sprintf("Revenue: %s", revenue))
			}
			if te, ok := company["total_emails"].(float64); ok && te > 0 {
				details = append(details, fmt.Sprintf("Emails: %.0f", te))
			}
			if ts, ok := company["total_similar"].(float64); ok && ts > 0 {
				details = append(details, fmt.Sprintf("Similar: %.0f", ts))
			}
			if pn, ok := company["phone_number"].(bool); ok && pn {
				details = append(details, "Phone: available")
			}
			if len(details) > 0 {
				sb.WriteString(fmt.Sprintf("       %s\n", strings.Join(details, "  |  ")))
			}

			// Social links
			linkedin := getStr(company, "linkedin_url")
			if linkedin == "" {
				linkedin = getStr(company, "linkedin")
			}
			twitter := getStr(company, "twitter_url")
			if twitter == "" {
				twitter = getStr(company, "twitter")
			}
			facebook := getStr(company, "facebook_url")
			if facebook == "" {
				facebook = getStr(company, "facebook")
			}
			var social []string
			if linkedin != "" {
				social = append(social, fmt.Sprintf("LinkedIn: %s", linkedin))
			}
			if twitter != "" {
				social = append(social, fmt.Sprintf("Twitter: %s", twitter))
			}
			if facebook != "" {
				social = append(social, fmt.Sprintf("Facebook: %s", facebook))
			}
			if len(social) > 0 {
				sb.WriteString(fmt.Sprintf("       %s\n", util.Gray(strings.Join(social, "  |  "))))
			}
		}
	}

	return sb.String()
}

func formatPhoneFinder(data map[string]interface{}) string {
	var sb strings.Builder

	if d, ok := data["data"].(map[string]interface{}); ok {
		// Contact info
		writeField(&sb, "\n  ", "Email", getStr(d, "email"))
		writeField(&sb, "  ", "Domain", getStr(d, "domain"))
		writeField(&sb, "  ", "Name", getStr(d, "full_name"))
		writeField(&sb, "  ", "Company", getStr(d, "company"))
		writeField(&sb, "  ", "Country", getStr(d, "country"))

		if phones, ok := d["phones"].([]interface{}); ok {
			sb.WriteString(fmt.Sprintf("\n%s Found %d phone number(s)\n\n", util.SuccessIcon(), len(phones)))
			sb.WriteString(fmt.Sprintf("  %-4s %-20s %-12s %-15s %-12s %s\n", "#", "Number", "Type", "Country", "Carrier", "Valid"))
			sb.WriteString(fmt.Sprintf("  %s\n", strings.Repeat("─", 80)))
			for i, p := range phones {
				if phone, ok := p.(map[string]interface{}); ok {
					number := getStr(phone, "phone")
					phoneType := getStr(phone, "type")
					phoneCountry := getStr(phone, "country")
					carrier := getStr(phone, "carrier")
					valid := formatBoolStr(phone, "valid")
					international := getStr(phone, "international_format")
					national := getStr(phone, "national_format")

					sb.WriteString(fmt.Sprintf("  %-4d %-20s %-12s %-15s %-12s %s\n",
						i+1, util.Bold(number), util.Gray(phoneType), phoneCountry, carrier, valid))
					if international != "" {
						sb.WriteString(fmt.Sprintf("       International: %s\n", international))
					}
					if national != "" {
						sb.WriteString(fmt.Sprintf("       National: %s\n", national))
					}
				}
			}
		}
	}

	return sb.String()
}

func formatPhoneValidator(data map[string]interface{}) string {
	var sb strings.Builder

	if d, ok := data["data"].(map[string]interface{}); ok {
		phone := getStr(d, "phone_number")
		valid, hasValid := d["valid"].(bool)

		if hasValid && valid {
			sb.WriteString(fmt.Sprintf("\n%s %s is %s\n", util.SuccessIcon(), util.Bold(phone), util.Green("valid")))
		} else {
			sb.WriteString(fmt.Sprintf("\n%s %s is %s\n", util.ErrorIcon(), util.Bold(phone), util.Red("invalid")))
		}

		sb.WriteString("\n  Phone Details:\n")
		writeField(&sb, "    ", "Phone", phone)
		writeField(&sb, "    ", "Local Format", getStr(d, "local_format"))
		writeField(&sb, "    ", "Intl Format", getStr(d, "intl_format"))
		writeField(&sb, "    ", "E164 Format", getStr(d, "e164_format"))
		writeField(&sb, "    ", "Country Prefix", getStr(d, "country_prefix"))
		writeField(&sb, "    ", "Country Code", getStr(d, "country_code"))
		writeField(&sb, "    ", "Country Name", getStr(d, "country_name"))
		writeField(&sb, "    ", "Carrier", getStr(d, "carrier"))
		writeField(&sb, "    ", "Line Type", getStr(d, "line_type"))
		writeField(&sb, "    ", "Location", getStr(d, "location"))
		writeField(&sb, "    ", "Timezone", getStr(d, "timezone"))
	}

	return sb.String()
}

func formatGeneric(data map[string]interface{}, indent int) string {
	var sb strings.Builder
	prefix := strings.Repeat("  ", indent)

	for k, v := range data {
		switch val := v.(type) {
		case map[string]interface{}:
			sb.WriteString(fmt.Sprintf("%s%s:\n", prefix, util.Bold(k)))
			sb.WriteString(formatGeneric(val, indent+1))
		case []interface{}:
			sb.WriteString(fmt.Sprintf("%s%s: (%d items)\n", prefix, util.Bold(k), len(val)))
			for i, item := range val {
				if m, ok := item.(map[string]interface{}); ok {
					sb.WriteString(fmt.Sprintf("%s  [%d]:\n", prefix, i+1))
					sb.WriteString(formatGeneric(m, indent+2))
				} else {
					sb.WriteString(fmt.Sprintf("%s  [%d]: %v\n", prefix, i+1, item))
				}
			}
		case string:
			if val != "" {
				sb.WriteString(fmt.Sprintf("%s%-18s %s\n", prefix, k+":", val))
			}
		case float64:
			sb.WriteString(fmt.Sprintf("%s%-18s %.0f\n", prefix, k+":", val))
		case bool:
			if val {
				sb.WriteString(fmt.Sprintf("%s%-18s %s\n", prefix, k+":", util.Green("true")))
			} else {
				sb.WriteString(fmt.Sprintf("%s%-18s %s\n", prefix, k+":", util.Red("false")))
			}
		case nil:
			// skip nil values
		default:
			sb.WriteString(fmt.Sprintf("%s%-18s %v\n", prefix, k+":", val))
		}
	}

	return sb.String()
}

// writeField writes a labeled field only if the value is non-empty
func writeField(sb *strings.Builder, prefix, label, value string) {
	if value != "" {
		sb.WriteString(fmt.Sprintf("%s%-14s %s\n", prefix, label+":", value))
	}
}

// writeBoolField writes a boolean field with color coding
func writeBoolField(sb *strings.Builder, prefix, label string, m map[string]interface{}, key string) {
	if v, ok := m[key].(bool); ok {
		if v {
			sb.WriteString(fmt.Sprintf("%s%-14s %s\n", prefix, label+":", util.Green("yes")))
		} else {
			sb.WriteString(fmt.Sprintf("%s%-14s %s\n", prefix, label+":", util.Red("no")))
		}
	}
}

// formatBoolStr returns a colored string for a boolean field
func formatBoolStr(m map[string]interface{}, key string) string {
	if v, ok := m[key].(bool); ok {
		if v {
			return util.Green("yes")
		}
		return util.Red("no")
	}
	return ""
}

func getStr(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func getFloat(m map[string]interface{}, key string) float64 {
	if v, ok := m[key].(float64); ok {
		return v
	}
	return 0
}
