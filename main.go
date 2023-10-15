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
	"github.com/nivb52/hotel-rent/api/middleware"
	"github.com/nivb52/hotel-rent/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dburi = db.DBURI
const dbname = db.DBNAME

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

	var (
		app        = fiber.New(appConfig)
		apiGeneral = app.Group("api")
		apiv1      = app.Group("api/v1")

		userStore    = db.NewMongoUserStore(client, dbname)
		hotelStore   = db.NewMongoHotelStore(client, dbname)
		roomStore    = db.NewMongoRoomStore(client, hotelStore, dbname)
		bookingStore = db.NewMongoBookingStore(client, dbname)
		store        = db.Store{
			User:    userStore,
			Hotel:   hotelStore,
			Room:    roomStore,
			Booking: bookingStore,
		}

		userHandler    = api.NewUserHandler(userStore)
		hotelHandler   = api.NewHotelHandler(&store)
		authHandler    = api.NewAuthHandler(userStore)
		bookingHandler = api.NewBookingHandler(&store)
	)

	// ROUTES
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Nothing here ! (choose api/version)")
	})

	apiGeneral.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Nothing here ! (choose api version - latest is 1 )")
	})

	// AUTH
	apiGeneral.Post("/auth", authHandler.HandleAuthenticate)

	//API V1
	apiv1.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World V1 👋!")
	})

	// ROUTES - USERS
	apiv1User := apiv1.Group("/users", middleware.JWTAuthentication)
	apiv1User.Get("/", userHandler.HandleGetUsers)
	apiv1User.Get("/:id", userHandler.HandleGetUser)
	apiv1User.Post("/", userHandler.HandleCreateUser)
	apiv1User.Delete("/:id", userHandler.HandleDeleteUser)
	apiv1User.Put("/:id", userHandler.HandleUpdateUser)

	// ROUTES - HOTEL
	apiv1Hotel := apiv1.Group("/hotels")
	apiv1Hotel.Get("/", hotelHandler.HandleGetHotels)
	apiv1Hotel.Get("/:id", hotelHandler.HandleGetHotel)
	apiv1Hotel.Get("/:id/rooms", hotelHandler.HandleGetHotelRooms)

	// ROUTES - ROOMS
	apiv1.Post("/room/:id/book", middleware.JWTAuthentication, bookingHandler.BookARoomByUser)
	apiv1.Post("/room/:id/gbook", bookingHandler.BookARoomByGuest)
	apiv1.Get("/room/:id/bookings", bookingHandler.GetBookings)

	// INIT
	app.Listen(*listenAddr)
}
