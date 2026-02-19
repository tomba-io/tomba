package start

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/tomba-io/go/tomba"
	"github.com/tomba-io/go/tomba/models"

	"github.com/tomba-io/tomba/pkg/slack"
	_domain "github.com/tomba-io/tomba/pkg/validation/domain"
	_email "github.com/tomba-io/tomba/pkg/validation/email"
	_url "github.com/tomba-io/tomba/pkg/validation/url"
)

// Request structure
type Request struct {
	URL          string   `json:"url" form:"url"`
	Domain       string   `json:"domain" form:"domain"`
	Email        string   `json:"email" form:"email"`
	SlackText    string   `json:"text" form:"text"`
	FirstName    string   `json:"first_name" form:"first_name"`
	LastName     string   `json:"last_name" form:"last_name"`
	FullName     string   `json:"full_name" form:"full_name"`
	Phone        string   `json:"phone" form:"phone"`
	CountryCode  string   `json:"country_code" form:"country_code"`
	Linkedin     string   `json:"linkedin" form:"linkedin"`
	Full         bool     `json:"full" form:"full"`
	EnrichMobile bool     `json:"enrich_mobile" form:"enrich_mobile"`
	Query        string   `json:"query" form:"query"`
	Page         int      `json:"page" form:"page"`
	Country      []string `json:"country" form:"country"`
	Industry     []string `json:"industry" form:"industry"`
	Size         []string `json:"size" form:"size"`
}

func request(c *fiber.Ctx) Request {
	req := Request{}
	c.BodyParser(&req)
	return req
}

// Home redirect to tomba home page
func (d *Conn) Home(c *fiber.Ctx) error {
	return c.Redirect("http://tomba.io?ref=tomba_cli", 301)
}

// Author query author finder
func (init *Conn) Author(c *fiber.Ctx) error {
	req := request(c)
	url := req.URL
	if url == "" && c.QueryBool("slack") {
		url = req.SlackText
	}
	if !_url.IsValidURL(url) {
		log.Error(ErrArgumentsURL.Error())
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": ErrArgumentsURL.Error()})
	}
	result, err := init.Tomba.AuthorFinder(url)
	if err != nil {
		log.Error(ErrErrInvalidLogin.Error())
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": ErrErrInvalidLogin.Error()})
	}
	if result.Data.Email == "" {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Why doesn't the author finder return any result? https://help.tomba.io/en/questions/reasons-why-author-finder-don-t-find-any-result"})
	}
	if c.QueryBool("slack") {
		return c.Status(200).JSON(slack.FinderResponse(result))
	}
	return c.Status(200).JSON(result)
}

// Count query email counter
func (init *Conn) Count(c *fiber.Ctx) error {
	req := request(c)
	domain := req.Domain
	if domain == "" && c.QueryBool("slack") {
		domain = req.SlackText
	}
	if !_domain.IsValidDomain(domain) {
		log.Error(ErrArgumentsDomain.Error())
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": ErrArgumentsDomain.Error()})
	}
	result, err := init.Tomba.Count(domain)
	if err != nil {
		log.Error(ErrErrInvalidLogin.Error())
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": ErrErrInvalidLogin.Error()})
	}
	if result.Data.Total == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "TombaPublicWebCrawler is searching the internet for the best leads that relate to this company, but we don't have any for it yet. That said, we're working on it"})
	}
	return c.Status(200).JSON(result)
}

// Enrich query enrichment
func (init *Conn) Enrich(c *fiber.Ctx) error {
	req := request(c)
	email := req.Email
	if email == "" && c.QueryBool("slack") {
		email = req.SlackText
	}
	if !_email.IsValidEmail(email) {
		log.Error(ErrArgumentEmail.Error())
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": ErrArgumentEmail.Error()})
	}
	params := tomba.Params{"email": email}
	if req.EnrichMobile {
		params["enrich_mobile"] = true
	}
	result, err := init.Tomba.Enrichment(params)
	if err != nil {
		log.Error(ErrErrInvalidLogin.Error())
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": ErrErrInvalidLogin.Error()})
	}
	if result.Data.Email == "" {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Why doesn't the enrichment return any result? https://help.tomba.io/en/questions/reasons-why-enrichment-don-t-find-any-emails"})
	}
	if c.QueryBool("slack") {
		return c.Status(200).JSON(slack.FinderResponse(result))
	}
	return c.Status(200).JSON(result)
}

// Linkedin query linkedin finder
func (init *Conn) Linkedin(c *fiber.Ctx) error {
	req := request(c)
	url := req.URL
	if url == "" && c.QueryBool("slack") {
		url = req.SlackText
	}
	if !_url.IsValidURL(url) {
		log.Error(ErrArgumentsURL.Error())
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": ErrArgumentsURL.Error()})
	}
	params := tomba.Params{"url": url}
	if req.EnrichMobile {
		params["enrich_mobile"] = true
	}
	result, err := init.Tomba.LinkedinFinder(params)
	if err != nil {
		log.Error(ErrErrInvalidLogin.Error())
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": ErrErrInvalidLogin.Error()})
	}
	if result.Data.Email == "" {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Why doesn't the Linkedin return any result? https://help.tomba.io/en/questions/reasons-why-linkedin-don-t-find-any-emails"})
	}
	if c.QueryBool("slack") {
		return c.Status(200).JSON(slack.FinderResponse(result))
	}
	return c.Status(200).JSON(result)
}

// Logs query logs
func (init *Conn) Logs(c *fiber.Ctx) error {
	result, err := init.Tomba.Logs()
	if err != nil {
		log.Error(ErrErrInvalidLogin.Error())
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": ErrErrInvalidLogin.Error()})
	}
	return c.Status(200).JSON(result)
}

// Search query domain search
func (init *Conn) Search(c *fiber.Ctx) error {
	req := request(c)
	domain := req.Domain
	if domain == "" && c.QueryBool("slack") {
		domain = req.SlackText
	}
	if !_domain.IsValidDomain(domain) {
		log.Error(ErrArgumentsDomain.Error())
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": ErrArgumentsDomain.Error()})
	}
	result, err := init.Tomba.DomainSearch(tomba.Params{"domain": domain})
	if err != nil {
		log.Error(ErrErrInvalidLogin.Error())
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": ErrErrInvalidLogin.Error()})
	}
	if result.Meta.Total == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "TombaPublicWebCrawler is searching the internet for the best leads that relate to this company, but we don't have any for it yet. That said, we're working on it"})
	}
	if c.QueryBool("slack") {
		return c.Status(200).JSON(slack.SearchResponse(result))
	}
	return c.Status(200).JSON(result)
}

// Status query Domain status
func (init *Conn) Status(c *fiber.Ctx) error {
	req := request(c)
	if !_domain.IsValidDomain(req.Domain) {
		log.Error(ErrArgumentsDomain.Error())
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": ErrArgumentsDomain.Error()})
	}
	result, err := init.Tomba.Status(req.Domain)
	if err != nil {
		log.Error(ErrErrInvalidLogin.Error())
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": ErrErrInvalidLogin.Error()})
	}
	return c.Status(200).JSON(result)
}

// Usage query usage
func (init *Conn) Usage(c *fiber.Ctx) error {
	result, err := init.Tomba.Usage()
	if err != nil {
		log.Error(ErrErrInvalidLogin.Error())
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": ErrErrInvalidLogin.Error()})
	}
	return c.Status(200).JSON(result)
}

// Verify query email verifier
func (init *Conn) Verify(c *fiber.Ctx) error {
	req := request(c)
	email := req.Email
	if email == "" && c.QueryBool("slack") {
		email = req.SlackText
	}
	if !_email.IsValidEmail(email) {
		log.Error(ErrArgumentEmail.Error())
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": ErrArgumentEmail.Error()})
	}
	result, err := init.Tomba.EmailVerifier(tomba.Params{"email": email})
	if err != nil {
		log.Error(ErrErrInvalidLogin.Error())
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": ErrErrInvalidLogin.Error()})
	}
	if result.Data.Email.Email != "" {
		if result.Data.Email.Disposable {
			return c.Status(404).JSON(fiber.Map{"status": "error", "message": "The domain name is used by a disposable email addresses provider."})
		}
		if result.Data.Email.Webmail {
			return c.Status(404).JSON(fiber.Map{"status": "error", "message": "The domain name  is webmail provider."})
		}
		if c.QueryBool("slack") {
			return c.Status(200).JSON(slack.VerifyResponse(result))
		}
		return c.Status(200).JSON(result)
	}
	return c.Status(222).JSON(fiber.Map{"status": "error", "message": "The Email Verification failed because of an unexpected response from the remote SMTP server. This failure is outside of our control. We advise you to retry later."})
}

// Finder query email finder
func (init *Conn) Finder(c *fiber.Ctx) error {
	req := request(c)
	if !_domain.IsValidDomain(req.Domain) {
		log.Error(ErrArgumentsDomain.Error())
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": ErrArgumentsDomain.Error()})
	}
	if req.FirstName == "" && req.LastName == "" && req.FullName == "" {
		log.Error(ErrArgumentsFinder.Error())
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": ErrArgumentsFinder.Error()})
	}
	params := tomba.Params{"domain": req.Domain}
	if req.FullName != "" {
		params["full_name"] = req.FullName
	} else {
		params["first_name"] = req.FirstName
		params["last_name"] = req.LastName
	}
	if req.EnrichMobile {
		params["enrich_mobile"] = true
	}
	result, err := init.Tomba.EmailFinder(params)
	if err != nil {
		log.Error(ErrErrInvalidLogin.Error())
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": ErrErrInvalidLogin.Error()})
	}
	if result.Data.Email == "" {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Why doesn't the finder return any result? https://help.tomba.io/en/questions/reasons-why-finder-don-t-find-any-emails"})
	}
	return c.Status(200).JSON(result)
}

// PhoneFinder query phone finder
func (init *Conn) PhoneFinder(c *fiber.Ctx) error {
	req := request(c)
	params := tomba.Params{}

	if req.Email != "" {
		if !_email.IsValidEmail(req.Email) {
			log.Error(ErrArgumentEmail.Error())
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": ErrArgumentEmail.Error()})
		}
		params["email"] = req.Email
	} else if req.Domain != "" {
		if !_domain.IsValidDomain(req.Domain) {
			log.Error(ErrArgumentsDomain.Error())
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": ErrArgumentsDomain.Error()})
		}
		params["domain"] = req.Domain
	} else if req.Linkedin != "" {
		if !_url.IsValidLinkedInProfile(req.Linkedin) {
			log.Error(ErrArgumentsLinkedinURL.Error())
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": ErrArgumentsLinkedinURL.Error()})
		}
		params["linkedin"] = req.Linkedin
	} else {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "At least one of email, domain, or linkedin is required."})
	}

	if req.Full {
		params["full"] = true
	}

	result, err := init.Tomba.PhoneFinder(params)
	if err != nil {
		log.Error(ErrErrInvalidLogin.Error())
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": ErrErrInvalidLogin.Error()})
	}
	return c.Status(200).JSON(result)
}

// PhoneValidator query phone validator
func (init *Conn) PhoneValidator(c *fiber.Ctx) error {
	req := request(c)
	if req.Phone == "" {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Phone number is required."})
	}
	params := tomba.Params{"phone": req.Phone}
	if req.CountryCode != "" {
		params["country_code"] = req.CountryCode
	}
	result, err := init.Tomba.PhoneValidator(params)
	if err != nil {
		log.Error(ErrErrInvalidLogin.Error())
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": ErrErrInvalidLogin.Error()})
	}
	return c.Status(200).JSON(result)
}

// Reveal query company search
func (init *Conn) Reveal(c *fiber.Ctx) error {
	req := request(c)
	request := &models.RevealSearchRequest{}

	if req.Page > 0 {
		request.Page = req.Page
	} else {
		request.Page = 1
	}

	if req.Query != "" {
		request.Query = req.Query
	}

	if len(req.Country) > 0 || len(req.Industry) > 0 || len(req.Size) > 0 {
		request.Filters = &models.RevealSearchFilters{
			Company: &models.RevealCompanyFilters{},
		}
		if len(req.Country) > 0 {
			request.Filters.Company.LocationCountry = &models.RevealCircularFilter{Include: req.Country}
		}
		if len(req.Industry) > 0 {
			request.Filters.Company.Industry = &models.RevealCircularFilter{Include: req.Industry}
		}
		if len(req.Size) > 0 {
			request.Filters.Company.Size = &models.RevealCircularFilter{Include: req.Size}
		}
	}

	result, err := init.Tomba.SearchCompanies(request)
	if err != nil {
		log.Error(ErrErrInvalidLogin.Error())
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": ErrErrInvalidLogin.Error()})
	}
	return c.Status(200).JSON(result)
}

// Similar query similar domains
func (init *Conn) Similar(c *fiber.Ctx) error {
	req := request(c)
	if !_domain.IsValidDomain(req.Domain) {
		log.Error(ErrArgumentsDomain.Error())
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": ErrArgumentsDomain.Error()})
	}
	result, err := init.Tomba.SimilarDomains(req.Domain)
	if err != nil {
		log.Error(ErrErrInvalidLogin.Error())
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": ErrErrInvalidLogin.Error()})
	}
	return c.Status(200).JSON(result)
}

// Sources query email sources
func (init *Conn) Sources(c *fiber.Ctx) error {
	req := request(c)
	email := req.Email
	if email == "" && c.QueryBool("slack") {
		email = req.SlackText
	}
	if !_email.IsValidEmail(email) {
		log.Error(ErrArgumentEmail.Error())
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": ErrArgumentEmail.Error()})
	}
	result, err := init.Tomba.Sources(email)
	if err != nil {
		log.Error(ErrErrInvalidLogin.Error())
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": ErrErrInvalidLogin.Error()})
	}
	return c.Status(200).JSON(result)
}

// Technology query technology check
func (init *Conn) Technology(c *fiber.Ctx) error {
	req := request(c)
	if !_domain.IsValidDomain(req.Domain) {
		log.Error(ErrArgumentsDomain.Error())
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": ErrArgumentsDomain.Error()})
	}
	result, err := init.Tomba.TechnologyCheck(req.Domain)
	if err != nil {
		log.Error(ErrErrInvalidLogin.Error())
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": ErrErrInvalidLogin.Error()})
	}
	return c.Status(200).JSON(result)
}

// Whoami query account info
func (init *Conn) Whoami(c *fiber.Ctx) error {
	result, err := init.Tomba.Account()
	if err != nil {
		log.Error(ErrErrInvalidLogin.Error())
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": ErrErrInvalidLogin.Error()})
	}
	return c.Status(200).JSON(result)
}
