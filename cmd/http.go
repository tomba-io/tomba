package cmd

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/spf13/cobra"

	"github.com/tomba-io/tomba/pkg/start"
)

// httpCmd represents the http command
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "Runs a HTTP server (reverse proxy).",
	Long:  Long,
	Run:   httpRun,
}

// httpRun the actual work http
func httpRun(cmd *cobra.Command, args []string) {
	init := start.New(conn)
	app := fiber.New(fiber.Config{
		DisableStartupMessage: false,
		AppName:               "tomba",
	})
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${ip}  ${status} - ${latency} ${method} ${path} ${queryParams}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", //
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET,POST,PUT,DELETE,HEAD,OPTIONS",
		Next:         nil,
	}))
	// if you want to prevent crashes
	app.Use(recover.New())
	setUpRoutes(app, init)
	_ = app.Listen(`:` + strconv.Itoa(init.Port))
}

func setUpRoutes(app *fiber.App, conn *start.Conn) {
	app.Get("/", conn.Home)
	app.Post("author", conn.Author)
	app.Post("count", conn.Count)
	app.Post("enrich", conn.Enrich)
	app.Post("finder", conn.Finder)
	app.Post("linkedin", conn.Linkedin)
	app.Get("logs", conn.Logs)
	app.Post("phone-finder", conn.PhoneFinder)
	app.Post("phone-validator", conn.PhoneValidator)
	app.Post("reveal", conn.Reveal)
	app.Post("search", conn.Search)
	app.Post("similar", conn.Similar)
	app.Post("sources", conn.Sources)
	app.Post("status", conn.Status)
	app.Post("technology", conn.Technology)
	app.Get("usage", conn.Usage)
	app.Post("verify", conn.Verify)
	app.Get("whoami", conn.Whoami)

	// Format, Location, AutoComplete
	app.Post("format", conn.Format)
	app.Post("location", conn.Location)
	app.Post("autocomplete", conn.AutoCompleteHandler)

	// Enrichment lookups
	app.Post("companies/find", conn.CompanyFindHandler)
	app.Post("people/find", conn.PersonFindHandler)
	app.Post("combined/find", conn.CombinedFindHandler)

	// Leads CRUD
	app.Get("leads", conn.ListLeadsHandler)
	app.Post("leads", conn.CreateLeadHandler)
	app.Get("leads/:id", conn.GetLeadHandler)
	app.Put("leads/:id", conn.UpdateLeadHandler)
	app.Delete("leads/:id", conn.DeleteLeadHandler)

	// Lead Lists CRUD
	app.Get("leads_lists", conn.ListLeadsListsHandler)
	app.Post("leads_lists", conn.CreateLeadsListHandler)
	app.Get("leads_lists/:id", conn.GetLeadsListHandler)
	app.Put("leads_lists/:id", conn.UpdateLeadsListHandler)
	app.Delete("leads_lists/:id", conn.DeleteLeadsListHandler)

	// Attributes CRUD
	app.Get("attributes", conn.ListAttributesHandler)
	app.Post("attributes", conn.CreateAttributeHandler)
	app.Get("attributes/:id", conn.GetAttributeHandler)
	app.Put("attributes/:id", conn.UpdateAttributeHandler)
	app.Delete("attributes/:id", conn.DeleteAttributeHandler)

	// Keys
	app.Get("keys", conn.ListKeysHandler)
	app.Post("keys", conn.CreateKeyHandler)
	app.Get("keys/:id", conn.GetKeyHandler)
	app.Delete("keys/:id", conn.DeleteKeyHandler)
	app.Put("keys/:id", conn.ResetKeyHandler)

	// Flags
	app.Get("flag", conn.ListFlagsHandler)
	app.Post("flag", conn.CreateFlagHandler)
}
