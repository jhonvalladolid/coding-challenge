package docs

import (
	"embed"

	"github.com/gofiber/fiber/v2"
)

//go:embed openapi.yaml swagger.html
var files embed.FS

func Register(app *fiber.App) {
	app.Get("/openapi.yaml", func(c *fiber.Ctx) error {
		data, err := files.ReadFile("openapi.yaml")
		if err != nil {
			return err
		}
		c.Type("yaml")
		return c.Send(data)
	})

	app.Get("/docs", func(c *fiber.Ctx) error {
		data, err := files.ReadFile("swagger.html")
		if err != nil {
			return err
		}
		c.Type("html")
		return c.Send(data)
	})
}
