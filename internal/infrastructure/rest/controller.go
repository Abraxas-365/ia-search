package rest

import (
	"context"

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
		answer, err := app.GetGptResposeWithContext(context.Background(), requestBody.Query, 1000, "text-davinci-003", false)
		if err != nil {

			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(200).JSON(fiber.Map{"answer": answer})

	})

	//on develop dont use
	r.Get("/completition/chat", func(c *fiber.Ctx) error {
		var requestBody struct {
			Query string `json:"query"`
		}

		if err := c.BodyParser(&requestBody); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err})
		}
		answer, err := app.GetGptResposeWithContext(context.Background(), requestBody.Query, 1000, "gpt-3.5-turbo", true)
		if err != nil {

			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(200).JSON(fiber.Map{"answer": answer})

	})
}
