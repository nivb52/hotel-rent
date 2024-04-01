package middleware

import (
	"github.com/gofiber/fiber/v2"
	e "github.com/nivb52/hotel-rent/api/errors"
)

func IsAdminAuth(c *fiber.Ctx) error {
	isAdmin := c.Context().UserValue("isAdmin")
	if isAdmin == nil {
		return e.ErrForbidden(c)
	}
	isAdmin = isAdmin.(bool)
	if isAdmin == true {
		return c.Next()
	}

	return e.ErrForbidden(c)
}
