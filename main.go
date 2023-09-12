// Package main provides ...
package main

import (
	"flag"

	"github.com/gofiber/fiber/v2"
)
func main() {
	listenAddr := flag.String("listenAddr", ":5000", "API PORT")
    flag.Parse()


	app := fiber.New()
	apiv1 := app.Group("api/v1")
    app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Nothing here ! (choose api version)")
	})
	apiv1.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World V1 👋!")
	})
	
	app.Listen(*listenAddr)
}