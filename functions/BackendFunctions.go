package functions

import (
	"github.com/gofiber/fiber/v2"
)

func Hello(c *fiber.Ctx) error {
	return c.SendString("Hello ✋")
}

func Sup(c *fiber.Ctx) error {
	return c.SendString("sup")
}