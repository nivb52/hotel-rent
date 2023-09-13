// Package main provides ...
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/nivb52/hotel-rent/api"
	"github.com/nivb52/hotel-rent/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dburi = "mongodb://localhost:27017"
const dbname = "hotel-rent"
const userColl = "users"

var appConfig = fiber.Config{
	// Override default error handler
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError // Status code defaults to 500
		// Retrieve the custom status code if it's a *fiber.Error
		var e *fiber.Error
		if errors.As(err, &e) {
			code = e.Code
		}
		return ctx.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Internal Server Error: %s | code: %d", err.Error(), code))
	},
}

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}

	listenAddr := flag.String("listenAddr", ":5000", "API PORT")
	flag.Parse()

	app := fiber.New(appConfig)

	// ROUTES
	apiv1 := app.Group("api/v1")
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Nothing here ! (choose api version)")
	})
	apiv1.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World V1 👋!")
	})

	// ROUTES - USERS
	userHandler := api.NewUserHandler(db.NewMongoUserStore(client))
	apiv1User := apiv1.Group("/users")
	apiv1User.Get("/", userHandler.HandleGetUsers)
	apiv1User.Get("/:id", userHandler.HandleGetUser)
	apiv1User.Post("/", userHandler.HandleCreateUser)
	apiv1User.Delete("/:id", userHandler.HandleDeleteUser)

	// INIT
	app.Listen(*listenAddr)
}
