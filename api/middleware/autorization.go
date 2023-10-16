package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func IsAdminAuth(c *fiber.Ctx) error {
	isAdmin := c.Context().UserValue("isAdmin")
	if isAdmin == nil {
		return c.Status(fiber.StatusForbidden).SendString(fiber.ErrForbidden.Message)
	}
	isAdmin = isAdmin.(bool)
	if isAdmin == true {
		return nil
	}

	return c.Status(fiber.StatusForbidden).SendString(fiber.ErrForbidden.Message)
}
