package rest

import (
	"github.com/Abraxas-365/ia-search/internal/application"
	"github.com/gofiber/fiber/v2"
)

func ControllerFactory(fiberApp *fiber.App, app application.Application) {
	r := fiberApp.Group("/api")
	r.Get("/completition", func(c *fiber.Ctx) error {
		var requestBody struct {
			Query string `json:"query"`
		}

		if err := c.BodyParser(&requestBody); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err})
		}
		answer, err := app.GetGptResposeWithContext(c.Context(), requestBody.Query)
		if err != nil {

			return c.Status(500).JSON(fiber.Map{"error": err})
		}

		return c.Status(200).JSON(fiber.Map{"answer": answer})

	})

}
