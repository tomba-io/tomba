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
		AllowMethods: "GET,POST,HEAD,OPTIONS",
		Next:         nil,
	}))
	// if you want to prevent crashes
	app.Use(recover.New())
	setUpRoutes(app, init)
	app.Listen(`:` + strconv.Itoa(init.Port))
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
}
