package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/nivb52/hotel-rent/db"
	"github.com/nivb52/hotel-rent/types"
)

type BookingHandler struct {
	// bookingStore db.BookingStore
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		// bookingStore: bookingStore, // ?
		store: store,
	}
}

//	--- Make Bookings ---
//
// Make a booking without register
func (h *BookingHandler) BookARoomByGuest(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString(fiber.ErrNotImplemented.Message)

}

// Make a booking for register user
func (h *BookingHandler) BookARoomByUser(c *fiber.Ctx) error {

	var (
		roomID = c.Params("id")
		userID = c.Context().UserValue("userID")
		params types.BookingParamsForCreate
	)

	err := c.BodyParser(&params)
	if err != nil {
		fmt.Println("Booking Params to create - Failed to parse, due: ", err)
		return c.Status(fiber.ErrBadRequest.Code).SendString(fiber.ErrBadRequest.Message)
	}

	if userID == "" || roomID == "" {
		fmt.Printf("\n Booking Data missing RoomID: %s or UserID: %s", roomID, userID)
		return c.Status(fiber.ErrBadRequest.Code).SendString(fiber.ErrBadRequest.Message)
	}

	params.RoomID = roomID
	params.UserID = userID.(string)
	errors := params.Validate()
	if len(errors) > 0 {
		return c.Status(fiber.ErrBadRequest.Code).JSON(errors)
	}

	var whereClause types.BookingFilter
	whereClause.RoomID = roomID
	whereClause.FromDate = params.FromDate
	whereClause.TillDate = params.TillDate
	isBooked, err := h.store.Booking.IsRoomAvailable(c.Context(), &whereClause)
	if err != nil {
		fmt.Println("Booking Failed, due: ", err)
		return c.Status(fiber.ErrConflict.Code).SendString(fiber.ErrConflict.Message)
	}

	if isBooked {
		fmt.Println("Booked Already")
		return c.Status(fiber.ErrConflict.Code).SendString("Those dates where just booked for this room")
	}

	booking, err := h.store.Booking.InsertBooking(c.Context(), &params)
	if err != nil {
		fmt.Println("Booking Failed, due: ", err)
		return c.Status(fiber.ErrConflict.Code).SendString(fiber.ErrConflict.Message)
	}

	return c.JSON(booking)
}
