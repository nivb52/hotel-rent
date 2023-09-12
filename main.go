// Package main provides ...
package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/nivb52/hotel-rent/api"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dburi = "mongodb://localhost:27017"
func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		panic(err)
	}

	fmt.Println(client)
	
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
	
	apiv1User := app.Group("api/v1/user")

	apiv1User.Get("/", api.HandleGetUsers)
	apiv1User.Get("/:id", api.HandleGetUsers)
	app.Listen(*listenAddr)
}