package start

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	_domain "github.com/tomba-io/tomba/pkg/validation/domain"
	_email "github.com/tomba-io/tomba/pkg/validation/email"
	_url "github.com/tomba-io/tomba/pkg/validation/url"
)

// Home redirect to tomba home page
func (d *Conn) Home(c *fiber.Ctx) error {
	return c.Redirect("http://tomba.io?ref=tomba_cli", 301)
}

// Author query author finder
func (init *Conn) Author(c *fiber.Ctx) error {
	url := strings.ToLower(c.Query("url"))
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

	return c.Status(200).JSON(result)

}

// Count query email counter
func (init *Conn) Count(c *fiber.Ctx) error {
	domain := strings.ToLower(c.Query("domain"))
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
	email := strings.ToLower(c.Query("email"))

	if !_email.IsValidEmail(email) {
		log.Error(ErrArgumentEmail.Error())
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": ErrArgumentEmail.Error()})
	}
	result, err := init.Tomba.Enrichment(email)
	if err != nil {
		log.Error(ErrErrInvalidLogin.Error())
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": ErrErrInvalidLogin.Error()})
	}
	if result.Data.Email == "" {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Why doesn't the enrichment return any result? https://help.tomba.io/en/questions/reasons-why-enrichment-don-t-find-any-emails"})
	}
	return c.Status(200).JSON(result)
}

// Linkedin query linkedin finder
func (init *Conn) Linkedin(c *fiber.Ctx) error {
	url := strings.ToLower(c.Query("url"))

	if !_url.IsValidURL(url) {
		log.Error(ErrArgumentsURL.Error())
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": ErrArgumentsURL.Error()})
	}
	result, err := init.Tomba.LinkedinFinder(url)
	if err != nil {
		log.Error(ErrErrInvalidLogin.Error())
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": ErrErrInvalidLogin.Error()})
	}
	if result.Data.Email == "" {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Why doesn't the Linkedin return any result? https://help.tomba.io/en/questions/reasons-why-linkedin-don-t-find-any-emails"})
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
	domain := strings.ToLower(c.Query("domain"))

	if !_domain.IsValidDomain(domain) {
		log.Error(ErrArgumentsDomain.Error())
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": ErrArgumentsDomain.Error()})
	}
	result, err := init.Tomba.DomainSearch(domain, "10", "0")
	if err != nil {
		log.Error(ErrErrInvalidLogin.Error())
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": ErrErrInvalidLogin.Error()})
	}
	if result.Meta.Total == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "TombaPublicWebCrawler is searching the internet for the best leads that relate to this company, but we don't have any for it yet. That said, we're working on it"})
	}
	return c.Status(200).JSON(result)

}

// Status query Domain status
func (init *Conn) Status(c *fiber.Ctx) error {
	domain := strings.ToLower(c.Query("domain"))

	if !_domain.IsValidDomain(domain) {
		log.Error(ErrArgumentsDomain.Error())
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": ErrArgumentsDomain.Error()})
	}
	result, err := init.Tomba.Status(domain)
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
	email := strings.ToLower(c.Query("email"))

	if !_email.IsValidEmail(email) {
		log.Error(ErrArgumentEmail.Error())
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": ErrArgumentEmail.Error()})
	}

	result, err := init.Tomba.EmailVerifier(email)
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
		return c.Status(200).JSON(result)
	}
	return c.Status(222).JSON(fiber.Map{"status": "error", "message": "The Email Verification failed because of an unexpected response from the remote SMTP server. This failure is outside of our control. We advise you to retry later."})
}
