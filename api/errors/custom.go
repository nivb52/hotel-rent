package customerrors

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type ServerError struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Reason string `json:"reason"`
}

// impl. error iterface
func (e ServerError) Error() string {
	return e.Msg
}

func NewServerError(code int, msg string) ServerError {
	reason := ""
	return ServerError{
		Code:   code,
		Msg:    msg,
		Reason: reason,
	}
}

func ErrUnAuthorized() ServerError {
	return ServerError{
		Code:   http.StatusUnauthorized,
		Msg:    "Unauthorized",
		Reason: "",
	}
}

func ErrTokenExpired() ServerError {
	return ServerError{
		Code:   http.StatusUnauthorized,
		Msg:    "token expired",
		Reason: "",
	}
}

func ErrInvalidID(reason string) ServerError {
	return ServerError{
		Code:   http.StatusBadRequest,
		Msg:    "the provided ID is wrong",
		Reason: reason,
	}
}

func ErrForbidden(c *fiber.Ctx) error {
	return c.Status(fiber.StatusForbidden).SendString(fiber.ErrForbidden.Message)
}

func ErrBadRequest(c *fiber.Ctx, reason ...string) error {
	msg := "Invalid Input"
	if len(reason) > 0 {
		msg = reason[0]
	}
	return c.Status(fiber.ErrBadRequest.Code).SendString(msg)
}

func ErrResourceNotFound(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).SendString(fiber.ErrNotFound.Message)
}

func ErrMethodNotAllowed(c *fiber.Ctx, reason ...string) error {
	msg := fiber.ErrMethodNotAllowed.Message
	if len(reason) > 0 {
		msg = reason[0]
	}
	return c.Status(fiber.StatusMethodNotAllowed).SendString(msg)
}

func ErrConflict(c *fiber.Ctx, reason ...string) error {
	msg := fiber.ErrConflict.Message
	if len(reason) > 0 {
		msg = reason[0]
	}
	return c.Status(fiber.ErrConflict.Code).SendString(msg)
}

func ErrInternalServerError(c *fiber.Ctx, reason ...string) error {
	msg := fiber.ErrInternalServerError.Message
	if len(reason) > 0 {
		msg = reason[0]
	}
	return c.Status(fiber.StatusInternalServerError).SendString(msg)
}

func ErrNotImplemented(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString(fiber.ErrNotImplemented.Message)
}
