package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	e "github.com/nivb52/hotel-rent/api/errors"
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
		return e.ErrInternalServerError(c)
	}

	if booking == nil {
		return e.ErrResourceNotFound(c)
	}

	return c.JSON(booking)
}

// --- GET Bookings ---
func (h *BookingHandler) AdminGetBookings(c *fiber.Ctx) error {
	var (
		whereClause types.BookingFilter
	)

	bookings, err := h.store.Booking.GetBookings(c.Context(), &whereClause)
	if err != nil {
		fmt.Println("GetBookings with filter failed, due: ", err)
		return e.ErrInternalServerError(c)
	}

	if len(bookings) > 0 {
		return c.JSON(bookings)
	}

	return e.ErrResourceNotFound(c)
}

// --- GET Bookings ---
func (h *BookingHandler) GetBookingsByFilter(c *fiber.Ctx, opt *types.GetBookingOptions) error {
	var (
		roomID             = c.Params("id")
		whereClause        types.BookingFilter
		isUserBookingsOnly = true
	)

	if !opt.UserBookingOnly {
		isUserBookingsOnly = false
	}

	if roomID == "" && !isUserBookingsOnly {
		fmt.Printf("\n Booking Data missing RoomID: %s", roomID)
		return c.Status(fiber.ErrBadRequest.Code).SendString(fiber.ErrBadRequest.Message)
	}

	err := c.BodyParser(&whereClause)
	if err != nil {
		fmt.Println("Booking Params Filter  - Failed to parse, due: ", err)
		return c.Status(fiber.ErrBadRequest.Code).SendString(fiber.ErrBadRequest.Message)
	}

	fmt.Println("whereClause: ", whereClause)
	if !isUserBookingsOnly {
		whereClause.RoomID = roomID
	}

	isAdmin := c.Context().UserValue("isAdmin").(bool)
	if !isAdmin || isUserBookingsOnly {
		whereClause.UserID = c.Context().UserValue("userID").(string)
	}

	bookings, err := h.store.Booking.GetBookings(c.Context(), &whereClause)
	if err != nil {
		fmt.Println("GetBookings with filter failed, due: ", err)
		return e.ErrInternalServerError(c)
	}

	fmt.Println("bookings: ", bookings)
	if len(bookings) > 0 {
		return c.JSON(bookings)
	}

	return e.ErrResourceNotFound(c)
}

//	--- Make Bookings ---
//
// Make a booking without register
// not implemented
func (h *BookingHandler) BookARoomByGuest(c *fiber.Ctx) error {
	return e.ErrNotImplemented(c)

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
		return e.ErrBadRequest(c)
	}

	if userID == "" || roomID == "" {
		fmt.Printf("\n Booking Data missing RoomID: %s or UserID: %s", roomID, userID)
		return e.ErrBadRequest(c)
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
	isAvailable, err := h.store.Booking.IsRoomAvailable(c.Context(), &whereClause)
	if err != nil {
		fmt.Println("Booking Failed - no available room, due: ", err)
		return e.ErrConflict(c)
	}
	if !isAvailable {
		fmt.Println("Booked Already")
		return e.ErrConflict(c, "Those dates where just booked for this room")
	}

	booking, err := h.store.Booking.InsertBooking(c.Context(), &params)
	if err != nil {
		fmt.Println("Booking Insertion Failed, due: ", err)
		return e.ErrConflict(c)
	}

	return c.Status(fiber.StatusAccepted).JSON(booking)
}

// Cancel Booking
func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingsById(c.Context(), id)
	if err != nil {
		fmt.Println("GetBookingsById failed, due: ", err)
		return e.ErrInternalServerError(c)
	}

	userID, ok := c.Context().UserValue("userID").(string)
	if !ok {
		return e.ErrInternalServerError(c)
	}

	isSameId, err := db.CompareIDs(userID, booking.UserID)
	if !isSameId {
		fmt.Println("CompareIDs failed, due: ", err)
		return e.ErrMethodNotAllowed(c)
	}

	err = h.store.Booking.CancelBooking(c.Context(), id)
	if err != nil {
		fmt.Println("CancelBooking failed, due: ", err)
		return e.ErrInternalServerError(c)
	}

	return c.Status(fiber.StatusAccepted).JSON(map[string]string{"msg": "ok", "canceld": id})
}
