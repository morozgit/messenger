package main

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v3"
)

func main() {
	app := fiber.New()

	app.Post("/auth/login", func(c fiber.Ctx) error {
		jsonData := c.Body()
		bodyReader := bytes.NewReader(jsonData)

		resp, err := http.Post("http://127.0.0.1:8000/add_user", "application/json", bodyReader)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("auth service unavailable")
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("failed to read response")
		}

		c.Status(resp.StatusCode)
		return c.Send(body)
	})

	log.Fatal(app.Listen(":3000"))
}
