package skill

// BuiltinSkills returns the pre-installed skills that ship with tomba
func BuiltinSkills() []*Skill {
	return []*Skill{
		{
			Name:        "find-email",
			Description: "Find and verify someone's email from their name and company domain",
			Author:      "tomba",
			Version:     "1.0.0",
			Tags:        []string{"email", "finder", "verify"},
			Inputs: []SkillInput{
				{Name: "domain", Description: "Company domain (e.g., stripe.com)", Type: "domain", Required: true},
				{Name: "first_name", Description: "Person's first name", Type: "string", Required: true},
				{Name: "last_name", Description: "Person's last name", Type: "string", Required: true},
			},
			Steps: []SkillStep{
				{
					ID:     "find",
					Name:   "Find email",
					Action: "finder",
					Params: map[string]string{
						"domain":     "{{input.domain}}",
						"first_name": "{{input.first_name}}",
						"last_name":  "{{input.last_name}}",
					},
				},
				{
					ID:     "verify",
					Name:   "Verify email",
					Action: "verify",
					Params: map[string]string{
						"email": "{{steps.find.data.email}}",
					},
				},
			},
			Output: SkillOutput{
				Format: "summary",
				Title:  "Email Found & Verified",
				Fields: []string{"email", "score", "verification", "position", "company"},
			},
		},
		{
			Name:        "company-emails",
			Description: "Search all emails for a company and count by department",
			Author:      "tomba",
			Version:     "1.0.0",
			Tags:        []string{"search", "company", "domain"},
			Inputs: []SkillInput{
				{Name: "domain", Description: "Company domain (e.g., tomba.io)", Type: "domain", Required: true},
			},
			Steps: []SkillStep{
				{
					ID:     "count",
					Name:   "Count emails",
					Action: "count",
					Params: map[string]string{
						"domain": "{{input.domain}}",
					},
				},
				{
					ID:     "search",
					Name:   "Search emails",
					Action: "search",
					Params: map[string]string{
						"domain": "{{input.domain}}",
						"limit":  "50",
					},
				},
			},
			Output: SkillOutput{
				Format: "table",
				Title:  "Company Email Directory",
				Fields: []string{"email", "first_name", "last_name", "position", "department", "score"},
			},
		},
		{
			Name:        "enrich-verify",
			Description: "Enrich an email with contact data then verify deliverability",
			Author:      "tomba",
			Version:     "1.0.0",
			Tags:        []string{"enrich", "verify", "email"},
			Inputs: []SkillInput{
				{Name: "email", Description: "Email address to enrich and verify", Type: "email", Required: true},
			},
			Steps: []SkillStep{
				{
					ID:     "enrich",
					Name:   "Enrich contact",
					Action: "enrich",
					Params: map[string]string{
						"email": "{{input.email}}",
					},
				},
				{
					ID:     "verify",
					Name:   "Verify email",
					Action: "verify",
					Params: map[string]string{
						"email": "{{input.email}}",
					},
				},
				{
					ID:     "sources",
					Name:   "Find sources",
					Action: "sources",
					Params: map[string]string{
						"email": "{{input.email}}",
					},
				},
			},
			Output: SkillOutput{
				Format: "summary",
				Title:  "Email Intelligence Report",
				Fields: []string{"email", "name", "position", "company", "country", "linkedin", "twitter", "verification", "sources"},
			},
		},
		{
			Name:        "company-intel",
			Description: "Full company intelligence — emails, tech stack, similar domains",
			Author:      "tomba",
			Version:     "1.0.0",
			Tags:        []string{"company", "technology", "similar", "search"},
			Inputs: []SkillInput{
				{Name: "domain", Description: "Company domain", Type: "domain", Required: true},
			},
			Steps: []SkillStep{
				{
					ID:     "status",
					Name:   "Domain status",
					Action: "status",
					Params: map[string]string{
						"domain": "{{input.domain}}",
					},
				},
				{
					ID:     "count",
					Name:   "Email count",
					Action: "count",
					Params: map[string]string{
						"domain": "{{input.domain}}",
					},
				},
				{
					ID:     "search",
					Name:   "Search emails",
					Action: "search",
					Params: map[string]string{
						"domain": "{{input.domain}}",
						"limit":  "20",
					},
				},
				{
					ID:     "tech",
					Name:   "Technology stack",
					Action: "technology",
					Params: map[string]string{
						"domain": "{{input.domain}}",
					},
				},
				{
					ID:     "similar",
					Name:   "Similar companies",
					Action: "similar",
					Params: map[string]string{
						"domain": "{{input.domain}}",
					},
				},
			},
			Output: SkillOutput{
				Format: "summary",
				Title:  "Company Intelligence Report",
				Fields: []string{"domain", "status", "total_emails", "departments", "technologies", "similar_domains", "emails"},
			},
		},
		{
			Name:        "lead-gen",
			Description: "Find decision makers at a company by searching and filtering by department",
			Author:      "tomba",
			Version:     "1.0.0",
			Tags:        []string{"leads", "search", "sales"},
			Inputs: []SkillInput{
				{Name: "domain", Description: "Company domain", Type: "domain", Required: true},
				{Name: "department", Description: "Target department (e.g., executive, sales, marketing)", Type: "string", Required: false, Default: "executive"},
			},
			Steps: []SkillStep{
				{
					ID:     "search",
					Name:   "Search by department",
					Action: "search",
					Params: map[string]string{
						"domain":     "{{input.domain}}",
						"department": "{{input.department}}",
						"limit":      "50",
					},
				},
			},
			Output: SkillOutput{
				Format: "table",
				Title:  "Lead Generation Results",
				Fields: []string{"email", "first_name", "last_name", "position", "department", "score", "linkedin"},
			},
		},
		{
			Name:        "linkedin-intel",
			Description: "Get full intel from a LinkedIn profile — email, enrichment, phone, verification",
			Author:      "tomba",
			Version:     "1.0.0",
			Tags:        []string{"linkedin", "enrich", "verify", "phone"},
			Inputs: []SkillInput{
				{Name: "url", Description: "LinkedIn profile URL", Type: "linkedin", Required: true},
			},
			Steps: []SkillStep{
				{
					ID:     "linkedin",
					Name:   "Find email from LinkedIn",
					Action: "linkedin",
					Params: map[string]string{
						"url": "{{input.url}}",
					},
				},
				{
					ID:     "enrich",
					Name:   "Enrich contact",
					Action: "enrich",
					Params: map[string]string{
						"email": "{{steps.linkedin.data.email}}",
					},
				},
				{
					ID:     "verify",
					Name:   "Verify email",
					Action: "verify",
					Params: map[string]string{
						"email": "{{steps.linkedin.data.email}}",
					},
				},
			},
			Output: SkillOutput{
				Format: "summary",
				Title:  "LinkedIn Intelligence",
				Fields: []string{"email", "name", "position", "company", "country", "linkedin", "twitter", "phone", "verification"},
			},
		},
		{
			Name:        "domain-check",
			Description: "Quick health check on a domain — status, email count, disposable/webmail detection",
			Author:      "tomba",
			Version:     "1.0.0",
			Tags:        []string{"domain", "status", "count"},
			Inputs: []SkillInput{
				{Name: "domain", Description: "Domain to check", Type: "domain", Required: true},
			},
			Steps: []SkillStep{
				{
					ID:     "status",
					Name:   "Check domain status",
					Action: "status",
					Params: map[string]string{
						"domain": "{{input.domain}}",
					},
				},
				{
					ID:     "count",
					Name:   "Count available emails",
					Action: "count",
					Params: map[string]string{
						"domain": "{{input.domain}}",
					},
				},
			},
			Output: SkillOutput{
				Format: "summary",
				Title:  "Domain Health Check",
				Fields: []string{"domain", "webmail", "disposable", "total_emails", "personal", "generic", "departments"},
			},
		},
		{
			Name:        "author-verify",
			Description: "Find an article author's email and verify it",
			Author:      "tomba",
			Version:     "1.0.0",
			Tags:        []string{"author", "verify", "article"},
			Inputs: []SkillInput{
				{Name: "url", Description: "Article URL", Type: "url", Required: true},
			},
			Steps: []SkillStep{
				{
					ID:     "author",
					Name:   "Find author email",
					Action: "author",
					Params: map[string]string{
						"url": "{{input.url}}",
					},
				},
				{
					ID:     "verify",
					Name:   "Verify author email",
					Action: "verify",
					Params: map[string]string{
						"email": "{{steps.author.data.email}}",
					},
				},
			},
			Output: SkillOutput{
				Format: "summary",
				Title:  "Author Email Verified",
				Fields: []string{"email", "name", "position", "company", "verification"},
			},
		},
	}
}
