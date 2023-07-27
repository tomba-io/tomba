package cmd

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
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
		Format: "[${time}] ${ip}  ${status} - ${latency} ${method} ${path}\n",
	}))
	setUpRoutes(app, init)
	app.Listen(`:` + strconv.Itoa(init.Port))
}

func setUpRoutes(app *fiber.App, conn *start.Conn) {
	app.Get("/", conn.Home)
	app.Get("author", conn.Author)
	app.Get("count", conn.Count)
	app.Get("enrich", conn.Enrich)
	app.Get("linkedin", conn.Linkedin)
	app.Get("search", conn.Search)
	app.Get("status", conn.Status)
	app.Get("verify", conn.Verify)
}
