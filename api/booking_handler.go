package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/nivb52/hotel-rent/db"
	"github.com/nivb52/hotel-rent/types"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

// --- GET Booking by Id ---
func (h *BookingHandler) GetBookingsById(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingsById(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(fiber.ErrInternalServerError.Message)
	}

	return c.JSON(booking)
}

// --- GET Bookings ---
func (h *BookingHandler) GetBookings(c *fiber.Ctx) error {
	var (
		whereClause types.BookingFilter
	)

	bookings, err := h.store.Booking.GetBookings(c.Context(), &whereClause)
	if err != nil {
		fmt.Println("GetBookings with filter failed, due: ", err)
		return c.Status(fiber.StatusInternalServerError).SendString(fiber.ErrInternalServerError.Message)
	}

	if len(bookings) > 0 {
		return c.JSON(bookings)
	}

	return c.Status(fiber.StatusNotFound).SendString(fiber.ErrNotFound.Message)
}

// --- GET Bookings ---
func (h *BookingHandler) GetBookingsByFilter(c *fiber.Ctx, isUserOnlyBookings bool) error {
	var (
		roomID      = c.Params("id")
		whereClause types.BookingFilter
	)

	if roomID == "" && !isUserOnlyBookings {
		fmt.Printf("\n Booking Data missing RoomID: %s", roomID)
		return c.Status(fiber.ErrBadRequest.Code).SendString(fiber.ErrBadRequest.Message)
	}

	err := c.BodyParser(&whereClause)
	if err != nil {
		fmt.Println("Booking Params Filter  - Failed to parse, due: ", err)
		return c.Status(fiber.ErrBadRequest.Code).SendString(fiber.ErrBadRequest.Message)
	}

	fmt.Println("whereClause: ", whereClause)
	if !isUserOnlyBookings {
		whereClause.RoomID = roomID
	}

	isAdmin := c.Context().UserValue("isAdmin").(bool)
	if !isAdmin || isUserOnlyBookings {
		whereClause.UserID = c.Context().UserValue("userID").(string)
	}

	bookings, err := h.store.Booking.GetBookings(c.Context(), &whereClause)
	if err != nil {
		fmt.Println("GetBookings with filter failed, due: ", err)
		return c.Status(fiber.StatusInternalServerError).SendString(fiber.ErrInternalServerError.Message)
	}

	fmt.Println("bookings: ", bookings)
	if len(bookings) > 0 {
		return c.JSON(bookings)
	}

	return c.Status(fiber.StatusNotFound).SendString(fiber.ErrNotFound.Message)
}

//	--- Make Bookings ---
//
// Make a booking without register
// not implemented
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
		fmt.Println("Booking Failed - no available room, due: ", err)
		return c.Status(fiber.ErrConflict.Code).SendString(fiber.ErrConflict.Message)
	}

	if isBooked {
		fmt.Println("Booked Already")
		return c.Status(fiber.ErrConflict.Code).SendString("Those dates where just booked for this room")
	}

	booking, err := h.store.Booking.InsertBooking(c.Context(), &params)
	if err != nil {
		fmt.Println("Booking Insertion Failed, due: ", err)
		return c.Status(fiber.ErrConflict.Code).SendString(fiber.ErrConflict.Message)
	}

	return c.JSON(booking)
}

// Cancel Booking
func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingsById(c.Context(), id)
	if err != nil {
		fmt.Println("GetBookingsById failed, due: ", err)
		return c.Status(fiber.StatusInternalServerError).SendString(fiber.ErrInternalServerError.Message)
	}

	userID, ok := c.Context().UserValue("userID").(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).SendString(fiber.ErrInternalServerError.Message)
	}

	isSameId, err := db.CompareIDs(userID, booking.UserID)
	if !isSameId {
		fmt.Println("CompareIDs failed, due: ", err)
		return c.Status(fiber.StatusMethodNotAllowed).SendString(fiber.ErrMethodNotAllowed.Message)
	}

	err = h.store.Booking.CancelBooking(c.Context(), id)
	if err != nil {
		fmt.Println("CancelBooking failed, due: ", err)
		return c.Status(fiber.StatusInternalServerError).SendString(fiber.ErrInternalServerError.Message)
	}

	return c.Status(fiber.StatusAccepted).SendString("success")
}
